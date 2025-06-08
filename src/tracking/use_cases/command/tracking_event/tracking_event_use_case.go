package tracking_event

import (
	"ludiks/src/kernel/database"
	"ludiks/src/tracking/domain/models"
	domain_repositories "ludiks/src/tracking/domain/repositories"

	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrackingEventUseCase struct {
	progressionRepository domain_repositories.ProgressionRepository
	endUserRepository     domain_repositories.EndUserRepository
	circuitRepository     domain_repositories.CircuitRepository
}

func NewTrackingEventUseCase(
	progressionRepository domain_repositories.ProgressionRepository,
	endUserRepository domain_repositories.EndUserRepository,
	circuitRepository domain_repositories.CircuitRepository,
) *TrackingEventUseCase {
	return &TrackingEventUseCase{
		progressionRepository: progressionRepository,
		endUserRepository:     endUserRepository,
		circuitRepository:     circuitRepository,
	}
}

type TrackingEventResponse struct {
	Success          bool    `json:"success"`
	Updated          bool    `json:"updated"`
	Message          *string `json:"message"`
	StepCompleted    bool    `json:"stepCompleted"`
	CircuitCompleted bool    `json:"circuitCompleted"`
	AlreadyCompleted bool    `json:"alreadyCompleted"`
	Points           *int    `json:"points,omitempty"`
}

func (u *TrackingEventUseCase) Execute(command *TrackingEventCommand) *TrackingEventResponse {
	response := &TrackingEventResponse{
		Success:          true,
		Updated:          false,
		StepCompleted:    false,
		CircuitCompleted: false,
		AlreadyCompleted: false,
	}

	endUser, err := u.endUserRepository.Find(command.UserID)
	if err != nil {
		response.Success = false
		response.Success = false
		return response
	}

	circuitProgression, err := u.progressionRepository.FindUserProgressionByEventName(
		command.ProjectID,
		endUser.ID,
		command.EventName,
	)

	if err != nil && err != gorm.ErrRecordNotFound {
		response.Success = false
		return response
	}

	var circuit *models.Circuit

	if circuitProgression == nil || err == gorm.ErrRecordNotFound {
		circuit, err = u.circuitRepository.FindByEventName(command.ProjectID, command.EventName)
		if err != nil {
			response.Success = false
			return response
		}

		circuitProgression = models.StartUserProgression(
			uuid.New().String(),
			endUser.ID,
			circuit.ID,
		)

		circuitProgression, err = u.progressionRepository.CreateProgression(circuitProgression)
		if err != nil {
			response.Success = false
			return response
		}
	} else {
		circuit = &circuitProgression.Circuit
	}

	switch circuit.Type {
	case database.TypeObjective:
		return u.handleObjectiveCircuit(circuit, circuitProgression, command)
	case database.TypePoints, database.TypeActions, database.TypeHybrid:
		return u.handleProgressiveCircuit(circuit, circuitProgression, command)
	default:
		response := newTrackingEventResponse()
		response.Success = false
		response.Message = &[]string{"Invalid circuit type"}[0]
		return response
	}
}

func (u *TrackingEventUseCase) handleObjectiveCircuit(
	circuit *models.Circuit,
	progression *models.UserCircuitProgression,
	command *TrackingEventCommand,
) *TrackingEventResponse {
	response := newTrackingEventResponse()
	timestamp := time.Now()
	if command.Timestamp != nil {
		timestamp = *command.Timestamp
	}

	if progression.CompletedAt != nil {
		response.AlreadyCompleted = true
		response.Message = &[]string{"Circuit already completed"}[0]
		return response
	}

	var step *models.Step
	steps, err := u.circuitRepository.GetCircuitSteps(circuit.ID)
	if err != nil {
		response.Success = false
		response.Message = &[]string{"Failed to get circuit steps"}[0]
		return response
	}

	for _, s := range *steps {
		if s.EventName == command.EventName {
			step = &s
			break
		}
	}

	if step == nil {
		response.Success = false
		response.Message = &[]string{"Step not found for this event"}[0]
		return response
	}

	stepProgression := progression.GetStepProgressionByEventName(command.EventName)

	if stepProgression == nil {
		stepProgression = models.StartUserStepProgression(
			uuid.New().String(),
			progression.ID,
			step.ID,
		)
		stepProgression.StartedAt = timestamp
		response.Updated = true
		progression.AddStepProgression(stepProgression)
	}

	if stepProgression.Status == database.UserStepProgressionStatusCompleted {
		response.AlreadyCompleted = true
		response.Message = &[]string{"Step already completed"}[0]
		return response
	}

	pointToAdd := 1
	if command.Value != nil {
		pointToAdd = *command.Value
	}

	stepProgression.ProgressCount += pointToAdd
	progression.Points += pointToAdd
	if stepProgression.ProgressCount >= step.CompletionThreshold {
		stepProgression.Status = database.UserStepProgressionStatusCompleted
		stepProgression.CompletedAt = &timestamp
		response.StepCompleted = true
		response.Updated = true
		progression.UpdateStepProgression(stepProgression)
	}

	response.Points = &progression.Points

	progression, err = u.updateProgressionAndCheckCompletion(progression)
	if err != nil {
		response.Success = false
		response.Message = &[]string{"Failed to update progression"}[0]
		return response
	}

	response.CircuitCompleted = progression.CompletedAt != nil
	return response
}

func (u *TrackingEventUseCase) handleProgressiveCircuit(
	circuit *models.Circuit,
	progression *models.UserCircuitProgression,
	command *TrackingEventCommand,
) *TrackingEventResponse {
	response := newTrackingEventResponse()
	timestamp := time.Now()
	if command.Timestamp != nil {
		timestamp = *command.Timestamp
	}

	incrementValue := 1
	if command.Value != nil {
		incrementValue = *command.Value
	}
	progression.Points += incrementValue

	sortedSteps, err := u.circuitRepository.GetCircuitSteps(circuit.ID)
	if err != nil {
		response.Success = false
		response.Message = &[]string{"Failed to get circuit steps"}[0]
		return response
	}

	sort.Slice(*sortedSteps, func(i, j int) bool {
		return (*sortedSteps)[i].CompletionThreshold < (*sortedSteps)[j].CompletionThreshold
	})

	var previousThreshold int = 0

	for _, step := range *sortedSteps {
		stepProg := progression.GetStepProgressionByStepID(step.ID)

		if stepProg == nil {
			if progression.Points >= previousThreshold {
				stepProg = models.StartUserStepProgression(
					uuid.New().String(),
					progression.ID,
					step.ID,
				)
				stepProg.StartedAt = timestamp
				stepProg.ProgressCount = progression.Points
				progression.AddStepProgression(stepProg)
				response.Updated = true
			} else {
				break
			}
		}

		if stepProg.Status == database.UserStepProgressionStatusInProgress {
			if progression.Points >= step.CompletionThreshold {
				stepProg.Status = database.UserStepProgressionStatusCompleted
				stepProg.CompletedAt = &timestamp
				stepProg.ProgressCount = step.CompletionThreshold
				progression.UpdateStepProgression(stepProg)
				response.StepCompleted = true
				response.Updated = true
			} else {
				break
			}
		}

		previousThreshold = step.CompletionThreshold
	}

	response.Points = &progression.Points

	progression, err = u.updateProgressionAndCheckCompletion(progression)
	if err != nil {
		response.Success = false
		response.Message = &[]string{"Failed to update progression: " + err.Error()}[0]
		return response
	}

	response.CircuitCompleted = progression.CompletedAt != nil
	return response
}

func (u *TrackingEventUseCase) updateProgressionAndCheckCompletion(
	progression *models.UserCircuitProgression,
) (*models.UserCircuitProgression, error) {
	updatedProgression, err := u.progressionRepository.UpdateProgression(progression)
	if err != nil {
		return nil, err
	}

	steps, _ := u.circuitRepository.GetCircuitSteps(progression.CircuitID)
	if err != nil {
		return nil, err
	}

	if len(*steps) == len(*updatedProgression.StepProgressions) {
		allCompleted := true
		for _, step := range *updatedProgression.StepProgressions {
			if step.Status != database.UserStepProgressionStatusCompleted {
				allCompleted = false
				break
			}
		}

		if allCompleted {
			completedAt := time.Now()
			updatedProgression.CompletedAt = &completedAt
			updatedProgression, err = u.progressionRepository.UpdateProgression(updatedProgression)
			if err != nil {
				return nil, err
			}
		}
	}

	return updatedProgression, nil
}

func newTrackingEventResponse() *TrackingEventResponse {
	return &TrackingEventResponse{
		Success:          true,
		Updated:          false,
		StepCompleted:    false,
		CircuitCompleted: false,
		AlreadyCompleted: false,
	}
}

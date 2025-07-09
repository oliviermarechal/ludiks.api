package tracking_event

import (
	"errors"
	"ludiks/src/kernel/app/database"
	"ludiks/src/tracking/domain/models"
	domain_providers "ludiks/src/tracking/domain/providers"
	domain_repositories "ludiks/src/tracking/domain/repositories"

	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrackingEventUseCase struct {
	progressionRepository  domain_repositories.ProgressionRepository
	endUserRepository      domain_repositories.EndUserRepository
	circuitRepository      domain_repositories.CircuitRepository
	organizationRepository domain_repositories.OrganizationRepository
	billingUsageProvider   domain_providers.BillingUsageProvider
}

func NewTrackingEventUseCase(
	progressionRepository domain_repositories.ProgressionRepository,
	endUserRepository domain_repositories.EndUserRepository,
	circuitRepository domain_repositories.CircuitRepository,
	organizationRepository domain_repositories.OrganizationRepository,
	billingUsageProvider domain_providers.BillingUsageProvider,
) *TrackingEventUseCase {
	return &TrackingEventUseCase{
		progressionRepository:  progressionRepository,
		endUserRepository:      endUserRepository,
		circuitRepository:      circuitRepository,
		organizationRepository: organizationRepository,
		billingUsageProvider:   billingUsageProvider,
	}
}

type RewardsTrackingEventResponse struct {
	Name string `json:"name"`
}

type TrackingEventResponse struct {
	Success          bool                           `json:"success"`
	Updated          bool                           `json:"updated"`
	Message          *string                        `json:"message"`
	StepCompleted    bool                           `json:"stepCompleted"`
	CircuitCompleted bool                           `json:"circuitCompleted"`
	AlreadyCompleted bool                           `json:"alreadyCompleted"`
	Points           *int                           `json:"points,omitempty"`
	Rewards          []RewardsTrackingEventResponse `json:"rewards"`
}

func (u *TrackingEventUseCase) Execute(command *TrackingEventCommand) *TrackingEventResponse {
	response := &TrackingEventResponse{
		Success:          true,
		Updated:          false,
		StepCompleted:    false,
		CircuitCompleted: false,
		AlreadyCompleted: false,
		Rewards:          []RewardsTrackingEventResponse{},
	}

	endUser, err := u.endUserRepository.FindByExternalID(command.UserID)
	if err != nil {
		message := err.Error()
		response.Success = false
		response.Message = &message

		return response
	}

	organization, err := u.organizationRepository.FindByProjectID(command.ProjectID)
	if err != nil {
		message := err.Error()
		response.Success = false
		response.Message = &message

		return response
	}

	if organization.HasQuotasReached() {
		message := errors.New("monthly quota events reached").Error()
		response.Success = false
		response.Message = &message

		return response
	}
	organization.IncrementQuotaUsed()
	u.organizationRepository.IncrementQuotaUsed(organization)

	if organization.Plan != "free" {
		u.billingUsageProvider.IncrementUsage(organization.StripeCustomerID)
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
		stepRewards, _ := u.circuitRepository.GetStepRewards(step.ID)

		if stepRewards != nil && len(*stepRewards) > 0 {
			for _, reward := range *stepRewards {
				response.Rewards = append(response.Rewards, RewardsTrackingEventResponse{Name: reward.Name})
			}
		}

		stepProgression.Status = database.UserStepProgressionStatusCompleted
		stepProgression.CompletedAt = &timestamp
		response.StepCompleted = true
		response.Updated = true
		progression.UpdateStepProgression(stepProgression)

		nextStep := getNextStep(step, steps)
		if nextStep != nil && progression.GetStepProgressionByStepID(nextStep.ID) == nil {
			nextStepProg := models.StartUserStepProgression(
				uuid.New().String(),
				progression.ID,
				nextStep.ID,
			)
			nextStepProg.StartedAt = timestamp
			nextStepProg.ProgressCount = progression.Points // Commencer avec le total actuel
			progression.AddStepProgression(nextStepProg)
		}
	}

	response.Points = &progression.Points

	progression, err = u.updateProgressionAndCheckCompletion(progression, response)
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

	if progression.CompletedAt != nil {
		response.AlreadyCompleted = true
		response.Message = &[]string{"Circuit already completed"}[0]
		return response
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

	remainingPoints := incrementValue
	var lastCompletedStep *models.Step

	// Distribuer les points étape par étape
	for _, step := range *sortedSteps {
		if remainingPoints <= 0 {
			break
		}

		stepProg := progression.GetStepProgressionByStepID(step.ID)

		// Si pas de progression pour cette étape, la créer si on a assez de points
		if stepProg == nil {
			if progression.Points >= step.CompletionThreshold {
				stepProg = models.StartUserStepProgression(
					uuid.New().String(),
					progression.ID,
					step.ID,
				)
				stepProg.StartedAt = timestamp
				stepProg.ProgressCount = step.CompletionThreshold // Total cumulatif
				stepProg.Status = database.UserStepProgressionStatusCompleted
				stepProg.CompletedAt = &timestamp
				progression.AddStepProgression(stepProg)
				response.StepCompleted = true
				response.Updated = true
				lastCompletedStep = &step

				// Attribuer les points nécessaires pour compléter cette étape
				pointsNeeded := step.CompletionThreshold - (progression.Points - incrementValue)
				remainingPoints -= pointsNeeded

				// Récupérer les récompenses de l'étape
				stepRewards, _ := u.circuitRepository.GetStepRewards(step.ID)
				if stepRewards != nil && len(*stepRewards) > 0 {
					for _, reward := range *stepRewards {
						response.Rewards = append(response.Rewards, RewardsTrackingEventResponse{Name: reward.Name})
					}
				}
			} else {
				// Pas assez de points pour cette étape, arrêter
				break
			}
		} else {
			// Progression existante
			if stepProg.Status == database.UserStepProgressionStatusInProgress {
				// Vérifier si on a maintenant assez de points pour terminer cette étape
				if progression.Points >= step.CompletionThreshold {
					stepProg.ProgressCount = step.CompletionThreshold // Total cumulatif
					stepProg.Status = database.UserStepProgressionStatusCompleted
					stepProg.CompletedAt = &timestamp
					progression.UpdateStepProgression(stepProg)
					response.StepCompleted = true
					response.Updated = true
					lastCompletedStep = &step

					// Récupérer les récompenses de l'étape
					stepRewards, _ := u.circuitRepository.GetStepRewards(step.ID)
					if stepRewards != nil && len(*stepRewards) > 0 {
						for _, reward := range *stepRewards {
							response.Rewards = append(response.Rewards, RewardsTrackingEventResponse{Name: reward.Name})
						}
					}
				} else {
					// Mettre à jour le progress_count avec le total actuel
					stepProg.ProgressCount = progression.Points
					progression.UpdateStepProgression(stepProg)
					response.Updated = true
				}
			} else if stepProg.Status == database.UserStepProgressionStatusCompleted {
				// Étape déjà terminée, passer à la suivante
				continue
			}
		}
	}

	// Si on vient de terminer une étape, créer la progression pour l'étape suivante
	if lastCompletedStep != nil {
		nextStep := getNextStep(lastCompletedStep, sortedSteps)
		if nextStep != nil && progression.GetStepProgressionByStepID(nextStep.ID) == nil {
			nextStepProg := models.StartUserStepProgression(
				uuid.New().String(),
				progression.ID,
				nextStep.ID,
			)
			nextStepProg.StartedAt = timestamp
			nextStepProg.ProgressCount = progression.Points // Commencer avec le total actuel
			progression.AddStepProgression(nextStepProg)
		}
	}

	response.Points = &progression.Points

	progression, err = u.updateProgressionAndCheckCompletion(progression, response)
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
	response *TrackingEventResponse,
) (*models.UserCircuitProgression, error) {
	updatedProgression, err := u.progressionRepository.UpdateProgression(progression)
	if err != nil {
		return nil, err
	}

	steps, _ := u.circuitRepository.GetCircuitSteps(progression.CircuitID)

	if len(*steps) == len(*updatedProgression.StepProgressions) {
		allCompleted := true
		for _, step := range *updatedProgression.StepProgressions {
			if step.Status != database.UserStepProgressionStatusCompleted {
				allCompleted = false
				break
			}
		}

		if allCompleted {
			stepRewards, _ := u.circuitRepository.GetCircuitRewards(progression.CircuitID)

			if stepRewards != nil && len(*stepRewards) > 0 {
				for _, reward := range *stepRewards {
					response.Rewards = append(response.Rewards, RewardsTrackingEventResponse{Name: reward.Name})
				}
			}

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

func getNextStep(currentStep *models.Step, steps *[]models.Step) *models.Step {
	if currentStep == nil || steps == nil || len(*steps) == 0 {
		return nil
	}

	sorted := *steps
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].StepNumber != nil && sorted[j].StepNumber != nil {
			return *sorted[i].StepNumber < *sorted[j].StepNumber
		}
		return sorted[i].CompletionThreshold < sorted[j].CompletionThreshold
	})
	for idx, s := range sorted {
		if s.ID == currentStep.ID && idx+1 < len(sorted) {
			return &sorted[idx+1]
		}
	}
	return nil
}

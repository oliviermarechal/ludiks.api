package create_step

import (
	"ludiks/src/gamification/domain/models"
	domain_repositories "ludiks/src/gamification/domain/repositories"

	"github.com/google/uuid"
)

type CreateStepUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewCreateStepUseCase(circuitRepository domain_repositories.CircuitRepository) *CreateStepUseCase {
	return &CreateStepUseCase{circuitRepository: circuitRepository}
}

func (u *CreateStepUseCase) Execute(command *CreateStepCommand) (*models.Step, error) {
	circuit, err := u.circuitRepository.Find(command.CircuitID)
	if err != nil {
		return nil, err
	}

	step := models.CreateStep(uuid.New().String(), command.Name, command.Description, command.CompletionThreshold, circuit.ID, len(*circuit.Steps)+1, command.EventName)
	step, err = u.circuitRepository.CreateStep(step)
	if err != nil {
		return nil, err
	}

	return step, nil
}

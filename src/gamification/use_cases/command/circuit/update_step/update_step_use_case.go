package update_step

import (
	"ludiks/src/gamification/domain/models"
	domain_repositories "ludiks/src/gamification/domain/repositories"
)

type UpdateStepUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewUpdateStepUseCase(circuitRepository domain_repositories.CircuitRepository) *UpdateStepUseCase {
	return &UpdateStepUseCase{circuitRepository: circuitRepository}
}

func (u *UpdateStepUseCase) Execute(command *UpdateStepCommand) (*models.Step, error) {
	step, err := u.circuitRepository.UpdateStep(command.StepId, command.Name, command.Description, command.CompletionThreshold, command.EventName)
	if err != nil {
		return nil, err
	}

	return step, nil
}

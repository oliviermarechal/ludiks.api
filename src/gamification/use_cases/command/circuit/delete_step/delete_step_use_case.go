package delete_step

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
)

type DeleteStepUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewDeleteStepUseCase(circuitRepository domain_repositories.CircuitRepository) *DeleteStepUseCase {
	return &DeleteStepUseCase{circuitRepository: circuitRepository}
}

func (u *DeleteStepUseCase) Execute(command *DeleteStepCommand) error {
	return u.circuitRepository.DeleteStep(command.StepId)
}

package activate_circuit

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
)

type ActivateCircuitUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewActivateCircuitUseCase(circuitRepository domain_repositories.CircuitRepository) *ActivateCircuitUseCase {
	return &ActivateCircuitUseCase{circuitRepository: circuitRepository}
}

func (u *ActivateCircuitUseCase) Execute(command *ActivateCircuitCommand) error {
	return u.circuitRepository.Activate(command.CircuitID)
}

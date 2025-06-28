package delete_circuit

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
)

type DeleteCircuitUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewDeleteCircuitUseCase(circuitRepository domain_repositories.CircuitRepository) *DeleteCircuitUseCase {
	return &DeleteCircuitUseCase{circuitRepository: circuitRepository}
}

func (u *DeleteCircuitUseCase) Execute(command *DeleteCircuitCommand) error {
	return u.circuitRepository.Delete(command.CircuitID)
}

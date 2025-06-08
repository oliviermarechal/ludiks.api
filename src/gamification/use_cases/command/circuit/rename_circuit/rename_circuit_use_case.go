package rename_circuit

import (
	"ludiks/src/gamification/domain/models"
	domain_repositories "ludiks/src/gamification/domain/repositories"
)

type RenameCircuitUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewRenameCircuitUseCase(circuitRepository domain_repositories.CircuitRepository) *RenameCircuitUseCase {
	return &RenameCircuitUseCase{circuitRepository: circuitRepository}
}

func (u *RenameCircuitUseCase) Execute(command *RenameCircuitCommand) (*models.Circuit, error) {
	return u.circuitRepository.Rename(command.CircuitID, command.Name)
}

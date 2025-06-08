package create_circuit

import (
	"ludiks/src/gamification/domain/models"
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/kernel/database"

	"github.com/google/uuid"
)

type CreateCircuitUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewCreateCircuitUseCase(circuitRepository domain_repositories.CircuitRepository) *CreateCircuitUseCase {
	return &CreateCircuitUseCase{circuitRepository: circuitRepository}
}

func (u *CreateCircuitUseCase) Execute(command *CreateCircuitCommand) (*models.Circuit, error) {
	circuit := models.CreateCircuit(
		uuid.New().String(),
		command.Name,
		database.CircuitType(command.Type),
		command.ProjectID,
	)

	return u.circuitRepository.Create(circuit)
}

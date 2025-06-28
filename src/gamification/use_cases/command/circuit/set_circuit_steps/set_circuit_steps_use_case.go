package set_circuit_steps

import (
	"ludiks/src/gamification/domain/models"
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/kernel/app/database"

	"errors"
	"sort"

	"github.com/google/uuid"
)

type SetCircuitStepsUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewSetCircuitStepsUseCase(circuitRepository domain_repositories.CircuitRepository) *SetCircuitStepsUseCase {
	return &SetCircuitStepsUseCase{circuitRepository: circuitRepository}
}

func (u *SetCircuitStepsUseCase) Execute(command *SetCircuitStepsCommand) (*models.Circuit, error) {
	circuit, err := u.circuitRepository.Find(command.CircuitID)
	if err != nil {
		return nil, err
	}

	if circuit.Active {
		return nil, errors.New("circuit is active")
	}

	u.circuitRepository.DeleteSteps(command.CircuitID)

	i := 1

	if circuit.Type != database.TypeObjective {
		sort.Slice(command.Steps, func(i, j int) bool {
			return command.Steps[i].CompletionThreshold < command.Steps[j].CompletionThreshold
		})
	}

	for _, stepData := range command.Steps {
		step := models.CreateStep(uuid.New().String(), stepData.Name, stepData.Description, stepData.CompletionThreshold, circuit.ID, stepData.StepNumber, stepData.EventName)
		step, err = u.circuitRepository.CreateStep(step)
		if err != nil {
			return nil, err
		}

		i++
	}

	return u.circuitRepository.Find(command.CircuitID)
}

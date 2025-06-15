package reordering_steps

import (
	"ludiks/src/gamification/domain/models"
	domain_repositories "ludiks/src/gamification/domain/repositories"
)

type ReorderingStepsUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewReorderingStepsUseCase(circuitRepository domain_repositories.CircuitRepository) *ReorderingStepsUseCase {
	return &ReorderingStepsUseCase{circuitRepository: circuitRepository}
}

func (u *ReorderingStepsUseCase) Execute(command *ReorderingStepsCommand) (*models.Circuit, error) {
	circuit, err := u.circuitRepository.Find(command.CircuitID)
	if err != nil {
		return nil, err
	}

	stepsToUpdate := make(map[string]int)
	for _, step := range command.Steps {
		stepsToUpdate[step.StepId] = step.StepNumber
	}

	for _, step := range *circuit.Steps {
		if newStepNumber, exists := stepsToUpdate[step.ID]; exists {
			u.circuitRepository.UpdateStepNumber(step.ID, newStepNumber)
		}
	}

	return circuit, nil
}

package domain_repositories

import (
	"ludiks/src/tracking/domain/models"
)

type CircuitRepository interface {
	FindByEventName(projectId string, eventName string) (*models.Circuit, error)
	GetCircuitSteps(circuitId string) (*[]models.Step, error)
	GetStepRewards(stepID string) (*[]models.Reward, error)
	GetCircuitRewards(circuitID string) (*[]models.Reward, error)
}

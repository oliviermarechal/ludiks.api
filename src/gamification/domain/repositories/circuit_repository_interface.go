package domain_repositories

import (
	"ludiks/src/gamification/domain/models"
)

type CircuitRepository interface {
	Create(circuit *models.Circuit) (*models.Circuit, error)
	Rename(id string, name string) (*models.Circuit, error)
	CreateMultipleSteps(steps []*models.Step) ([]*models.Step, error)
	Find(id string) (*models.Circuit, error)
	Activate(id string) error
	CreateStep(step *models.Step) (*models.Step, error)
	DeleteSteps(circuitID string) error
	DeleteStep(stepId string) error
	UpdateStep(stepId string, name string, description string, completionThreshold int, eventName string) (*models.Step, error)
}

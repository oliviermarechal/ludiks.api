package infra_repositories

import (
	"ludiks/src/tracking/domain/models"

	"gorm.io/gorm"
)

type CircuitRepository struct {
	db *gorm.DB
}

func NewCircuitRepository(db *gorm.DB) *CircuitRepository {
	return &CircuitRepository{db: db}
}

func (r *CircuitRepository) FindByEventName(projectId string, eventName string) (*models.Circuit, error) {
	var circuit models.Circuit
	err := r.db.
		Joins("JOIN steps ON steps.circuit_id = circuits.id").
		Where("project_id = ? AND steps.event_name = ?", projectId, eventName).
		Preload("Steps").
		First(&circuit).
		Error
	if err != nil {
		return nil, err
	}
	return &circuit, nil
}

func (r *CircuitRepository) GetCircuitSteps(circuitId string) (*[]models.Step, error) {
	steps := []models.Step{}
	err := r.db.Where("circuit_id = ?", circuitId).Find(&steps).Error
	if err != nil {
		return nil, err
	}
	return &steps, nil
}

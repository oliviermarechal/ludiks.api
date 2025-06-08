package infra_repositories

import (
	"ludiks/src/gamification/domain/models"

	"gorm.io/gorm"
)

type CircuitRepository struct {
	db *gorm.DB
}

func NewCircuitRepository(db *gorm.DB) *CircuitRepository {
	return &CircuitRepository{db: db}
}

func (r *CircuitRepository) Create(circuit *models.Circuit) (*models.Circuit, error) {
	err := r.db.Create(circuit).Error
	if err != nil {
		return nil, err
	}

	return circuit, nil
}

func (r *CircuitRepository) Rename(id string, name string) (*models.Circuit, error) {
	err := r.db.Model(&models.Circuit{}).Where("id = ?", id).Update("name", name).Error
	if err != nil {
		return nil, err
	}

	return r.Find(id)
}

func (r *CircuitRepository) Find(id string) (*models.Circuit, error) {
	var circuit models.Circuit
	err := r.db.Where("id = ?", id).Preload("Steps").First(&circuit).Error
	if err != nil {
		return nil, err
	}

	return &circuit, nil
}

func (r *CircuitRepository) Activate(id string) error {
	err := r.db.Model(&models.Circuit{}).Where("id = ?", id).Update("active", true).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *CircuitRepository) CreateStep(step *models.Step) (*models.Step, error) {
	err := r.db.Create(step).Error
	if err != nil {
		return nil, err
	}

	return step, nil
}

func (r *CircuitRepository) CreateMultipleSteps(steps []*models.Step) ([]*models.Step, error) {
	err := r.db.Create(steps).Error
	if err != nil {
		return nil, err
	}

	return steps, nil
}

func (r *CircuitRepository) DeleteSteps(circuitID string) error {
	err := r.db.Where("circuit_id = ?", circuitID).Delete(&models.Step{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *CircuitRepository) UpdateStep(stepId string, name string, description string, completionThreshold int, eventName string) (*models.Step, error) {
	var existingStep models.Step
	if err := r.db.Where("id = ?", stepId).First(&existingStep).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	err := r.db.Model(&models.Step{}).Where("id = ?", stepId).Updates(map[string]interface{}{
		"name":                 name,
		"description":          description,
		"completion_threshold": completionThreshold,
		"event_name":           eventName,
	}).Error
	if err != nil {
		return nil, err
	}

	var step models.Step
	err = r.db.Where("id = ?", stepId).First(&step).Error
	if err != nil {
		return nil, err
	}

	return &step, nil
}

func (r *CircuitRepository) DeleteStep(stepId string) error {
	err := r.db.Where("id = ?", stepId).Delete(&models.Step{}).Error
	if err != nil {
		return err
	}

	return nil
}

package infra_repositories

import (
	"fmt"

	"gorm.io/gorm"

	"ludiks/src/tracking/domain/models"
)

type ProgressionRepository struct {
	db *gorm.DB
}

func NewProgressionRepository(db *gorm.DB) *ProgressionRepository {
	return &ProgressionRepository{db: db}
}

func (r *ProgressionRepository) FindUserProgressionByEventName(
	projectId string,
	endUserId string,
	eventName string,
) (*models.UserCircuitProgression, error) {
	var progression models.UserCircuitProgression

	if err := r.db.Joins("JOIN circuits ON user_circuit_progressions.circuit_id = circuits.id").
		Joins("JOIN steps ON steps.circuit_id = circuits.id").
		Where("circuits.project_id = ? AND user_circuit_progressions.end_user_id = ? AND steps.event_name = ?",
			projectId,
			endUserId,
			eventName,
		).Preload("StepProgressions").
		Preload("StepProgressions.Step").
		Preload("Circuit").
		First(&progression).Error; err != nil {
		return nil, err
	}

	return &progression, nil
}

func (r *ProgressionRepository) CreateProgression(progression *models.UserCircuitProgression) (*models.UserCircuitProgression, error) {
	if err := r.db.
		Preload("StepProgressions").
		Preload("StepProgressions.Step").
		Preload("Circuit").
		Create(progression).Error; err != nil {
		return nil, err
	}

	return progression, nil
}

func (r *ProgressionRepository) UpdateProgression(progression *models.UserCircuitProgression) (*models.UserCircuitProgression, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		println("Transaction begin error:", tx.Error.Error())
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(progression).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if progression.StepProgressions != nil {
		for i := range *progression.StepProgressions {
			step := &(*progression.StepProgressions)[i]
			var existingStep models.UserStepProgression
			err := tx.Where("id = ?", step.ID).First(&existingStep).Error

			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(step).Error; err != nil {
					tx.Rollback()
					return nil, fmt.Errorf("failed to create step: %w", err)
				}
			} else if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to find step: %w", err)
			} else {
				if err := tx.Save(step).Error; err != nil {
					tx.Rollback()
					return nil, fmt.Errorf("failed to save step: %w", err)
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	var updatedProgression models.UserCircuitProgression
	if err := r.db.Where("id = ?", progression.ID).Preload("StepProgressions").First(&updatedProgression).Error; err != nil {
		return nil, fmt.Errorf("failed to reload progression: %w", err)
	}

	return &updatedProgression, nil
}

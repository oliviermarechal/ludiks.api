package models

import (
	"ludiks/src/kernel/database"
	"time"
)

type UserCircuitProgression struct {
	ID               string                 `json:"id" gorm:"primarykey"`
	Points           int                    `json:"points"`
	StartedAt        time.Time              `json:"started_at"`
	CompletedAt      *time.Time             `json:"completed_at"`
	EndUserID        string                 `json:"end_user_id"`
	EndUser          EndUser                `json:"end_user"`
	CircuitID        string                 `json:"circuit_id"`
	Circuit          Circuit                `json:"circuit,omitempty"`
	StepProgressions *[]UserStepProgression `json:"step_progressions,omitempty" gorm:"foreignKey:UserCircuitProgressionID;constraint:OnDelete:CASCADE"`
}

func StartUserProgression(id string, endUserID string, circuitID string) *UserCircuitProgression {
	return &UserCircuitProgression{
		ID:        id,
		Points:    0,
		StartedAt: time.Now(),
		EndUserID: endUserID,
		CircuitID: circuitID,
	}
}

func (u *UserCircuitProgression) StepIsCompeleted(eventName string) bool {
	if u.StepProgressions == nil || len(*u.StepProgressions) == 0 {
		return false
	}

	for i := len(*u.StepProgressions) - 1; i >= 0; i-- {
		if (*u.StepProgressions)[i].Step.EventName == eventName {
			return (*u.StepProgressions)[i].Status == database.UserStepProgressionStatusCompleted
		}
	}

	return false
}

func (u *UserCircuitProgression) GetStepProgression() *UserStepProgression {
	if u.StepProgressions == nil || len(*u.StepProgressions) == 0 {
		return nil
	}

	for i := len(*u.StepProgressions) - 1; i >= 0; i-- {
		if (*u.StepProgressions)[i].Status != database.UserStepProgressionStatusInProgress {
			return &(*u.StepProgressions)[i]
		}
	}

	return nil
}

func (u *UserCircuitProgression) GetStepProgressionByEventName(eventName string) *UserStepProgression {
	if u.StepProgressions == nil || len(*u.StepProgressions) == 0 {
		return nil
	}

	for i := len(*u.StepProgressions) - 1; i >= 0; i-- {
		if (*u.StepProgressions)[i].Step.EventName == eventName {
			return &(*u.StepProgressions)[i]
		}
	}

	return nil
}

func (u *UserCircuitProgression) GetStepProgressionByStepID(stepID string) *UserStepProgression {
	if u.StepProgressions == nil || len(*u.StepProgressions) == 0 {
		return nil
	}

	for i := len(*u.StepProgressions) - 1; i >= 0; i-- {
		if (*u.StepProgressions)[i].Step.ID == stepID {
			return &(*u.StepProgressions)[i]
		}
	}

	return nil
}

func (u *UserCircuitProgression) AddStepProgression(stepProgression *UserStepProgression) {
	if u.StepProgressions == nil {
		u.StepProgressions = &[]UserStepProgression{}
	}
	*u.StepProgressions = append(*u.StepProgressions, UserStepProgression{
		ID:                       stepProgression.ID,
		StepID:                   stepProgression.StepID,
		UserCircuitProgressionID: u.ID,
		Status:                   stepProgression.Status,
		StartedAt:                stepProgression.StartedAt,
		CompletedAt:              stepProgression.CompletedAt,
		ProgressCount:            stepProgression.ProgressCount,
	})
}

func (u *UserCircuitProgression) UpdateStepProgression(stepProgression *UserStepProgression) bool {
	if u.StepProgressions == nil {
		u.StepProgressions = &[]UserStepProgression{}
		return false
	}

	for i := range *u.StepProgressions {
		if (*u.StepProgressions)[i].ID == stepProgression.ID {
			(*u.StepProgressions)[i] = *stepProgression
			return true
		}
	}

	return false
}

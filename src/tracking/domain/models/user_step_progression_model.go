package models

import (
	"time"

	"ludiks/src/kernel/app/database"
)

type UserStepProgression struct {
	ID                       string                             `json:"id"`
	StartedAt                time.Time                          `json:"started_at"`
	CompletedAt              *time.Time                         `json:"completed_at"`
	UserCircuitProgressionID string                             `json:"user_circuit_progression_id"`
	StepID                   string                             `json:"step_id"`
	Step                     Step                               `json:"step"`
	Status                   database.UserStepProgressionStatus `json:"status"`
	ProgressCount            int                                `json:"progress_count"`
}

func StartUserStepProgression(id string, userCircuitProgressionID string, stepID string) *UserStepProgression {
	return &UserStepProgression{
		ID:                       id,
		UserCircuitProgressionID: userCircuitProgressionID,
		StepID:                   stepID,
		StartedAt:                time.Now(),
		Status:                   database.UserStepProgressionStatusInProgress,
		ProgressCount:            0,
	}
}

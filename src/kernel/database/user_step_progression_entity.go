package database

import (
	"time"

	"github.com/google/uuid"
)

type UserStepProgressionStatus string

const (
	UserStepProgressionStatusInProgress UserStepProgressionStatus = "in_progress"
	UserStepProgressionStatusCompleted  UserStepProgressionStatus = "completed"
)

type UserStepProgressionEntity struct {
	ID                       uuid.UUID                    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	StartedAt                time.Time                    `gorm:"type:timestamp"`
	CompletedAt              time.Time                    `gorm:"type:timestamp"`
	UserCircuitProgressionID uuid.UUID                    `gorm:"type:uuid;not null"`
	UserCircuitProgression   UserCircuitProgressionEntity `gorm:"foreignKey:UserCircuitProgressionID"`
	StepID                   uuid.UUID                    `gorm:"type:uuid;not null"`
	Step                     StepEntity                   `gorm:"foreignKey:StepID"`
	Status                   UserStepProgressionStatus    `gorm:"type:text"`
	ProgressCount            int                          `gorm:"type:int;default:0"`
}

func (UserStepProgressionEntity) TableName() string {
	return "user_step_progressions"
}

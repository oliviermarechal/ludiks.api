package database

import (
	"time"

	"github.com/google/uuid"
)

type UserCircuitProgressionEntity struct {
	ID              uuid.UUID                   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Points          int                         `gorm:"type:int;default:0"`
	StartedAt       time.Time                   `gorm:"type:timestamp"`
	CompletedAt     time.Time                   `gorm:"type:timestamp"`
	EndUserID       uuid.UUID                   `gorm:"type:uuid;not null"`
	EndUser         EndUserEntity               `gorm:"foreignKey:EndUserID"`
	CircuitID       uuid.UUID                   `gorm:"type:uuid;not null"`
	Circuit         CircuitEntity               `gorm:"foreignKey:CircuitID"`
	StepProgression []UserStepProgressionEntity `gorm:"foreignKey:UserCircuitProgressionID"`
}

func (UserCircuitProgressionEntity) TableName() string {
	return "user_circuit_progressions"
}

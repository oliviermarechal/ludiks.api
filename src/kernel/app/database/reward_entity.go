package database

import (
	"time"

	"github.com/google/uuid"
)

type RewardEntity struct {
	ID                        uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name                      string        `gorm:"not null"`
	Description               string        `gorm:"default:null"`
	UnlockOnCircuitCompletion bool          `gorm:"not null;default:false"`
	CircuitID                 uuid.UUID     `gorm:"type:uuid;not null"`
	StepID                    *uuid.UUID    `gorm:"type:uuid;default:null"`
	Circuit                   CircuitEntity `gorm:"foreignKey:CircuitID"`
	Step                      *StepEntity   `gorm:"foreignKey:StepID"`
	CreatedAt                 time.Time     `gorm:"autoCreateTime"`
}

func (RewardEntity) TableName() string {
	return "rewards"
}

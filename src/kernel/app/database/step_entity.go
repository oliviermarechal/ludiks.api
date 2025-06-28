package database

import (
	"time"

	"github.com/google/uuid"
)

type StepEntity struct {
	ID                  uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name                string        `gorm:"not null"`
	Description         string        `gorm:"default:null"`
	CompletionThreshold int           `gorm:"not null;default:0"`
	CircuitID           uuid.UUID     `gorm:"type:uuid;not null"`
	StepNumber          *int          `gorm:"default:null"`
	EventName           string        `gorm:"not null"`
	Circuit             CircuitEntity `gorm:"foreignKey:CircuitID"`
	CreatedAt           time.Time     `gorm:"autoCreateTime"`
}

func (StepEntity) TableName() string {
	return "steps"
}

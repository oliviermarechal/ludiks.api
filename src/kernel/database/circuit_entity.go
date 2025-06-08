package database

import (
	"time"

	"github.com/google/uuid"
)

type CircuitType string

const (
	TypePoints    CircuitType = "points"
	TypeActions   CircuitType = "actions"
	TypeObjective CircuitType = "objective"
	TypeHybrid    CircuitType = "hybrid"
)

type CircuitEntity struct {
	ID        uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string        `gorm:"type:varchar(255);not null"`
	Type      CircuitType   `gorm:"type:varchar(55);not null"`
	Active    bool          `gorm:"type:boolean;not null"`
	ProjectID uuid.UUID     `gorm:"type:uuid;not null"`
	Project   ProjectEntity `gorm:"foreignKey:ProjectID"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
}

func (CircuitEntity) TableName() string {
	return "circuits"
}

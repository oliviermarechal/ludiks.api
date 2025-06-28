package database

import (
	"time"

	"github.com/google/uuid"
)

type ApiKeyEntity struct {
	ID        uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string        `gorm:"type:varchar(255);not null"`
	Value     string        `gorm:"type:varchar(255);not null"`
	ProjectID uuid.UUID     `gorm:"type:uuid;not null"`
	Project   ProjectEntity `gorm:"foreignKey:ProjectID"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
}

func (ApiKeyEntity) TableName() string {
	return "api_keys"
}

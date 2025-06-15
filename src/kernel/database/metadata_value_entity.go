package database

import (
	"time"

	"github.com/google/uuid"
)

type MetadataValueEntity struct {
	ID                   uuid.UUID                 `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectMetadataKeyID uuid.UUID                 `gorm:"type:uuid;not null;index"`
	ProjectMetadataKey   ProjectMetadataKeysEntity `gorm:"foreignKey:ProjectMetadataKeyID;references:ID"`
	Value                string                    `gorm:"type:varchar(100);not null"`
	CreatedAt            time.Time                 `gorm:"autoCreateTime"`
}

func (MetadataValueEntity) TableName() string {
	return "metadata_values"
}

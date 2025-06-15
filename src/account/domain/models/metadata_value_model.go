package models

import (
	"time"

	"github.com/google/uuid"
)

type MetadataValue struct {
	ID                   uuid.UUID           `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectMetadataKeyID uuid.UUID           `json:"projectMetadataKeyId" gorm:"type:uuid;not null;index"`
	ProjectMetadataKey   *ProjectMetadataKey `json:"project,omitempty" gorm:"foreignKey:ProjectMetadataKeyID;references:ID"`
	Value                string              `json:"value" gorm:"type:varchar(100);not null"`
	CreatedAt            time.Time           `json:"createdAt" gorm:"autoCreateTime"`
}

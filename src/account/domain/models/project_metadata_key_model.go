package models

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMetadataKey struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectID uuid.UUID       `json:"projectId" gorm:"type:uuid;not null"`
	KeyName   string          `json:"keyName" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	Values    []MetadataValue `json:"values" gorm:"foreignKey:ProjectMetadataKeyID;references:ID"`
}

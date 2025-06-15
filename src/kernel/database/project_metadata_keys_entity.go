package database

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMetadataKeysEntity struct {
	ID        uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectID uuid.UUID             `gorm:"type:uuid;not null;uniqueIndex:idx_project_key"`
	Project   ProjectEntity         `gorm:"foreignKey:ProjectID"`
	KeyName   string                `gorm:"type:varchar(100);not null;uniqueIndex:idx_project_key"`
	CreatedAt time.Time             `gorm:"autoCreateTime"`
	Values    []MetadataValueEntity `gorm:"foreignKey:ProjectMetadataKeyID;references:ID"`
}

func (ProjectMetadataKeysEntity) TableName() string {
	return "project_metadata_keys"
}

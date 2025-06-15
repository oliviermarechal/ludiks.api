package models

import (
	"time"
)

type MetadataValue struct {
	ID                   string    `json:"id"`
	ProjectMetadataKeyID string    `json:"project_metadata_key_id"`
	Value                string    `json:"value"`
	CreatedAt            time.Time `json:"created_at"`
}

func CreateMetadataValue(id string, projectMetadataKeyId string, value string) *MetadataValue {
	return &MetadataValue{
		ID:                   id,
		ProjectMetadataKeyID: projectMetadataKeyId,
		Value:                value,
		CreatedAt:            time.Now(),
	}
}

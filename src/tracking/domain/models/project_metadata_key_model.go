package models

import (
	"fmt"
	"time"
)

type ProjectMetadataKey struct {
	ID        string           `json:"id"`
	ProjectID string           `json:"project_id"`
	KeyName   string           `json:"key_name"`
	CreatedAt time.Time        `json:"created_at"`
	Values    []*MetadataValue `json:"values"`
}

func CreateProjectMetadataKey(ID string, ProjectID string, KeyName string) *ProjectMetadataKey {
	return &ProjectMetadataKey{
		ID:        ID,
		ProjectID: ProjectID,
		KeyName:   KeyName,
		CreatedAt: time.Now(),
		Values:    []*MetadataValue{},
	}
}

func (p *ProjectMetadataKey) HasMetadata(keyName string) bool {
	return p.KeyName == keyName
}

func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (p *ProjectMetadataKey) HasValue(value interface{}) bool {
	for _, v := range p.Values {
		if v.Value == value {
			return true
		}
	}

	return false
}

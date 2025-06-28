package models

import "time"

type ApiKey struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	ProjectID string    `json:"projectId" gorm:"type:uuid;not null;foreignKey:ProjectID;references:ID"`
	Project   *Project  `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	CreatedAt time.Time `json:"createdAt"`
}

func CreateApiKey(id string, project_id string, name string, value string) *ApiKey {
	return &ApiKey{
		ID:        id,
		ProjectID: project_id,
		Name:      name,
		Value:     value,
		CreatedAt: time.Now(),
	}
}

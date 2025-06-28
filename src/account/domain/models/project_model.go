package models

import "time"

type Project struct {
	ID             string        `json:"id" gorm:"primaryKey"`
	Name           string        `json:"name"`
	CreatedAt      time.Time     `json:"createdAt"`
	OrganizationId string        `json:"organizationId"`
	Organization   *Organization `json:"organization,omitempty"`
	ApiKeys        *[]ApiKey     `json:"apiKeys,omitempty" gorm:"foreignKey:ProjectID"`
}

func CreateProject(id string, name string, organizationId string) *Project {
	return &Project{
		ID:             id,
		Name:           name,
		OrganizationId: organizationId,
		CreatedAt:      time.Now(),
	}
}

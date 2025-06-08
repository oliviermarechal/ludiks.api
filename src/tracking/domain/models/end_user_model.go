package models

import "time"

type EndUser struct {
	ID         string    `json:"id"`
	ExternalID string    `json:"external_id"`
	FullName   string    `json:"full_name"`
	Email      *string   `json:"email"`
	Picture    *string   `json:"picture"`
	CreatedAt  time.Time `json:"created_at"`
	ProjectID  string    `json:"project_id"`
}

func CreateEndUser(ID string, ExternalID string, FullName string, Email *string, Picture *string, ProjectID string) *EndUser {
	return &EndUser{
		ID:         ID,
		ExternalID: ExternalID,
		FullName:   FullName,
		Email:      Email,
		Picture:    Picture,
		ProjectID:  ProjectID,
		CreatedAt:  time.Now(),
	}
}

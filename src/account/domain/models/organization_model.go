package models

import "time"

type Organization struct {
	ID          string                   `json:"id" gorm:"primaryKey"`
	Name        string                   `json:"name"`
	Plan        string                   `json:"plan"`
	EventsQuota int                      `json:"eventsQuota"`
	EventsUsed  int                      `json:"eventsUsed"`
	Pricing     int                      `json:"pricing"`
	CreatedAt   time.Time                `json:"created_at"`
	Projects    *[]Project               `json:"projects,omitempty"`
	Memberships []OrganizationMembership `json:"memberships,omitempty"`
}

func CreateOrganization(id string, name string) *Organization {
	return &Organization{
		ID:          id,
		Name:        name,
		CreatedAt:   time.Now(),
		EventsQuota: 5000,
		Plan:        "free",
	}
}

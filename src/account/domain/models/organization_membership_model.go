package models

import "time"

type OrganizationMembership struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	UserID         string    `json:"userId"`
	OrganizationID string    `json:"organizationId"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"createdAt"`
}

func CreateOrganizationMember(id string, userID string, organizationID string, role string) *OrganizationMembership {
	return &OrganizationMembership{
		ID:             id,
		UserID:         userID,
		OrganizationID: organizationID,
		Role:           role,
		CreatedAt:      time.Now(),
	}
}

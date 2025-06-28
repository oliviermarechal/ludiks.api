package models

import "time"

type OrganizationInvitation struct {
	ID             string        `json:"id" gorm:"primaryKey"`
	FromUserID     string        `json:"fromUserId"`
	FromUser       *User         `json:"fromUser,omitempty"`
	ToEmail        string        `json:"toEmail"`
	OrganizationID string        `json:"organizationId"`
	Organization   *Organization `json:"organization,omitempty"`
	Role           string        `json:"role"`
	CreatedAt      time.Time     `json:"createdAt"`
}

func CreateOrganizationInvitation(id string, fromUserID string, toEmail string, organizationID string, role string) *OrganizationInvitation {
	return &OrganizationInvitation{
		ID:             id,
		FromUserID:     fromUserID,
		ToEmail:        toEmail,
		OrganizationID: organizationID,
		Role:           role,
		CreatedAt:      time.Now(),
	}
}

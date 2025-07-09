package models

import "time"

type OrganizationSubscription struct {
	ID                   string     `json:"id" gorm:"primaryKey"`
	OrganizationID       string     `json:"organization_id"`
	StripeSubscriptionID string     `json:"-"`
	Status               string     `json:"status"`
	CurrentPeriodEnd     *time.Time `json:"current_period_end"`
	CancelRequestedAt    *time.Time `json:"cancel_requested_at"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func CreateOrganizationSubscription(
	id string,
	organizationID string,
	stripeSubscriptionID string,
) *OrganizationSubscription {
	now := time.Now()
	return &OrganizationSubscription{
		ID:                   id,
		OrganizationID:       organizationID,
		StripeSubscriptionID: stripeSubscriptionID,
		Status:               "active",
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

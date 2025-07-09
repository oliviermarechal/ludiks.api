package models

import "time"

type Invoice struct {
	ID              string    `json:"id"`
	StripeInvoiceID string    `json:"stripe_invoice_id"`
	OrganizationID  string    `json:"organization_id"`
	Amount          int       `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	PeriodMonth     int       `json:"period_month"`
	PeriodYear      int       `json:"period_year"`
	LinkUrl         *string   `json:"link_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func CreateInvoice(
	id string,
	stripeInvoiceID string,
	organizationID string,
	amount int,
	currency string,
	status string,
	periodMonth int,
	periodYear int,
	linkUrl *string,
) *Invoice {
	return &Invoice{
		ID:              id,
		StripeInvoiceID: stripeInvoiceID,
		OrganizationID:  organizationID,
		Amount:          amount,
		Currency:        currency,
		Status:          status,
		PeriodMonth:     periodMonth,
		PeriodYear:      periodYear,
		LinkUrl:         linkUrl,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

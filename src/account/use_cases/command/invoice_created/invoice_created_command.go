package invoice_created

import "time"

type InvoiceCreatedCommand struct {
	StripeInvoiceID string
	OrganizationID  string
	Amount          int
	Currency        string
	Status          string
	PeriodMonth     int
	PeriodYear      int
	LinkUrl         *string
	CreatedAt       time.Time
}

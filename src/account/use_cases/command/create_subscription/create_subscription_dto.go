package create_subscription

import "github.com/go-playground/validator/v10"

type CreateSubscriptionDTO struct {
	CustomerEmail string `json:"customer_email" validate:"required,email"`
	// Informations de facturation B2B
	CompanyName    string `json:"company_name" validate:"required"`
	CompanyAddress struct {
		Line1      string `json:"line1" validate:"required"`
		Line2      string `json:"line2"`
		City       string `json:"city" validate:"required"`
		State      string `json:"state"`
		PostalCode string `json:"postal_code" validate:"required"`
		Country    string `json:"country" validate:"required"`
	} `json:"company_address" validate:"required"`
	// Informations de contact
	ContactName  string `json:"contact_name" validate:"required"`
	ContactPhone string `json:"contact_phone"`
	// Informations fiscales (optionnel, Stripe peut auto-détecter le type)
	TaxIDValue string `json:"tax_id_value"`
	// ID de l'organisation
	OrganizationID string `json:"organization_id" validate:"required"`
}

func (dto *CreateSubscriptionDTO) Validate() error {
	return validator.New().Struct(dto)
}

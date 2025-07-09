package domain_repositories

import "ludiks/src/account/domain/models"

type InvoiceRepository interface {
	Create(invoice *models.Invoice) (*models.Invoice, error)
	FindByOrganizationID(organizationID string) ([]*models.Invoice, error)
}

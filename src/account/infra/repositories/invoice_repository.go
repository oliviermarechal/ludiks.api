package infra_repositories

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (i *InvoiceRepository) Create(invoice *models.Invoice) (*models.Invoice, error) {
	if err := i.db.Create(invoice).Error; err != nil {
		return nil, err
	}

	return invoice, nil
}

func (i *InvoiceRepository) FindByOrganizationID(organizationID string) ([]*models.Invoice, error) {
	var invoices []*models.Invoice
	if err := i.db.Where("organization_id = ?", organizationID).Find(invoices).Error; err != nil {
		return invoices, err
	}

	return invoices, nil
}

package invoice_created

import (
	"ludiks/src/account/domain/models"

	domain_repositories "ludiks/src/account/domain/repositories"

	"github.com/google/uuid"
)

type InvoiceCreatedUseCase struct {
	invoiceRepository domain_repositories.InvoiceRepository
}

func NewInvoiceCreatedUseCase(
	invoiceRepository domain_repositories.InvoiceRepository,
) *InvoiceCreatedUseCase {
	return &InvoiceCreatedUseCase{
		invoiceRepository: invoiceRepository,
	}
}

func (n *InvoiceCreatedUseCase) Execute(command InvoiceCreatedCommand) error {
	invoice := models.CreateInvoice(
		uuid.New().String(),
		command.StripeInvoiceID,
		command.OrganizationID,
		command.Amount,
		command.Currency,
		command.Status,
		command.PeriodMonth,
		command.PeriodYear,
		command.LinkUrl,
	)

	_, err := n.invoiceRepository.Create(invoice)
	if err != nil {
		return err
	}

	return nil
}

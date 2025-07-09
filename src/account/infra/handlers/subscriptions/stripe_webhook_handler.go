package subscriptions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/invoice_created"
	"ludiks/src/kernel/app/handlers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeWebhookHandler struct {
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository
	organizationRepository             domain_repositories.OrganizationRepository
	invoiceRepository                  domain_repositories.InvoiceRepository
	signingSecret                      string
}

func NewStripeWebhookHandler(
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository,
	organizationRepository domain_repositories.OrganizationRepository,
	invoiceRepository domain_repositories.InvoiceRepository,
	signingSecret string,
) *StripeWebhookHandler {
	return &StripeWebhookHandler{
		organizationSubscriptionRepository: organizationSubscriptionRepository,
		organizationRepository:             organizationRepository,
		invoiceRepository:                  invoiceRepository,
		signingSecret:                      signingSecret,
	}
}

func (h *StripeWebhookHandler) Handle(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, sigHeader, h.signingSecret)

	if err != nil {
		handlers.HandleBadRequest(c, errors.New("erreur"))
		return
	}

	fmt.Printf("Received event: %s - %s\n", event.ID, event.Type)

	switch event.Type {
	// TODO Case invoice payment failed
	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			handlers.HandleBadRequest(c, err)
			return
		}

		// TODO Update invoice status to paid
		fmt.Printf("invoice: %+v\n", invoice)

	case "invoice.created":
		var invoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			handlers.HandleBadRequest(c, err)
			return
		}

		periodTime := time.Unix(invoice.PeriodStart, 0)
		periodMonth := int(periodTime.Month())
		periodYear := periodTime.Year()

		var linkUrl *string
		if invoice.HostedInvoiceURL != "" {
			linkUrl = &invoice.HostedInvoiceURL
		}

		customerID, err := extractCustomerID(&invoice, event.Data.Raw)
		if err != nil {
			handlers.HandleBadRequest(c, err)
			return
		}
		fmt.Printf("customerID: %s\n", customerID)

		organization, err := h.organizationRepository.FindByStripeCustomerID("cus_SbxL5zczs1uBgm") // customerID)
		if err != nil {
			handlers.HandleBadRequest(c, err)
			return
		}

		command := invoice_created.InvoiceCreatedCommand{
			StripeInvoiceID: invoice.ID,
			OrganizationID:  organization.ID,
			Amount:          int(invoice.AmountDue),
			Currency:        string(invoice.Currency),
			Status:          string(invoice.Status),
			PeriodMonth:     periodMonth,
			PeriodYear:      periodYear,
			LinkUrl:         linkUrl,
			CreatedAt:       time.Unix(invoice.Created, 0),
		}

		useCase := invoice_created.NewInvoiceCreatedUseCase(h.invoiceRepository)
		if err := useCase.Execute(command); err != nil {
			handlers.HandleBadRequest(c, err)
			return
		}
	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	c.Status(http.StatusOK)
}

func extractCustomerID(invoice *stripe.Invoice, eventData json.RawMessage) (string, error) {
	if invoice.Customer == nil {
		return "", errors.New("customer field is nil")
	}

	var rawInvoice map[string]interface{}
	if err := json.Unmarshal(eventData, &rawInvoice); err != nil {
		return "", fmt.Errorf("failed to parse raw invoice data: %v", err)
	}

	customerField, exists := rawInvoice["customer"]
	if !exists {
		return "", errors.New("customer field not found in invoice data")
	}

	if customerID, ok := customerField.(string); ok {
		return customerID, nil
	}

	if customerObj, ok := customerField.(map[string]interface{}); ok {
		if customerID, ok := customerObj["id"].(string); ok {
			return customerID, nil
		}
	}

	return invoice.Customer.ID, nil
}

package create_subscription

import (
	"ludiks/config"
	domain_providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
)

type CreateSubscriptionResult struct {
	ClientSecret   string `json:"clientSecret"`
	CustomerID     string `json:"customerId"`
	SubscriptionID string `json:"subscriptionId,omitempty"`
}

type CreateSubscriptionUseCase struct {
	stripeProvider         domain_providers.StripeProvider
	organizationRepository domain_repositories.OrganizationRepository
}

func NewCreateSubscriptionUseCase(
	stripeProvider domain_providers.StripeProvider,
	organizationRepository domain_repositories.OrganizationRepository,
) *CreateSubscriptionUseCase {
	return &CreateSubscriptionUseCase{
		stripeProvider:         stripeProvider,
		organizationRepository: organizationRepository,
	}
}

func (c *CreateSubscriptionUseCase) Execute(command CreateSubscriptionCommand) (*CreateSubscriptionResult, error) {
	customerInfo := domain_providers.CustomerInfo{
		Email:       command.CustomerEmail,
		CompanyName: command.CompanyName,
		CompanyAddress: domain_providers.CustomerInfoCompanyAddress{
			Line1:      command.CompanyAddress.Line1,
			Line2:      command.CompanyAddress.Line2,
			City:       command.CompanyAddress.City,
			State:      command.CompanyAddress.State,
			PostalCode: command.CompanyAddress.PostalCode,
			Country:    command.CompanyAddress.Country,
		},
		ContactName:  command.ContactName,
		ContactPhone: command.ContactPhone,
		TaxIDValue:   command.TaxIDValue,
	}

	priceID := config.AppConfig.StripeSubscriptionPriceID

	organization, err := c.organizationRepository.Find(command.OrganizationID)
	if err != nil {
		return nil, err
	}

	customerID := organization.StripeCustomerID
	if customerID == nil {
		result, err := c.stripeProvider.CreateCustomer(customerInfo, organization.ID)
		if err != nil {
			return nil, err
		}

		customerID = &result
		c.organizationRepository.SetStripeCustomerID(command.OrganizationID, result)
	}

	result, err := c.stripeProvider.SetupPaymentMethod(*customerID, &priceID, command.OrganizationID)
	if err != nil {
		return nil, err
	}

	return &CreateSubscriptionResult{
		ClientSecret:   result.ClientSecret,
		SubscriptionID: result.SubscriptionID,
	}, nil
}

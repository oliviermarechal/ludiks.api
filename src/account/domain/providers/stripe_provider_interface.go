package providers

type CustomerInfo struct {
	Email          string
	CompanyName    string
	CompanyAddress CustomerInfoCompanyAddress
	ContactName    string
	ContactPhone   string
	TaxIDValue     string
}

type CustomerInfoCompanyAddress struct {
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	Country    string
}

type SetupPaymentMethodResult struct {
	ClientSecret   string
	SubscriptionID string
}

type CancelSubscriptionResult struct {
	SubscriptionID string
	CancelledAt    string
	EndDate        string
}

type SubscriptionDetails struct {
	ID     string
	Status string
}

type StripeProvider interface {
	CreateCustomer(customerInfo CustomerInfo, organizationID string) (string, error)
	SetupPaymentMethod(customerID string, priceID *string, organizationID string) (*SetupPaymentMethodResult, error)
	GetOrganizationIDFromIntentPayment(setupIntentID string) (string, error)
	CancelSubscription(subscriptionID string, cancelAtPeriodEnd bool) (*CancelSubscriptionResult, error)
	GetSubscriptionDetails(subscriptionID string) (*SubscriptionDetails, error)
	ReportUsage(customerID string) error
}

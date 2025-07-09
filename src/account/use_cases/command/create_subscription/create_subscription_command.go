package create_subscription

type CreateSubscriptionCommand struct {
	CustomerEmail  string
	CompanyName    string
	CompanyAddress CreateSubscriptionCommandCompanyAddress
	ContactName    string
	ContactPhone   string
	TaxIDValue     string
	OrganizationID string
}

type CreateSubscriptionCommandCompanyAddress struct {
	Line1      string
	Line2      string
	City       string
	State      string
	PostalCode string
	Country    string
}

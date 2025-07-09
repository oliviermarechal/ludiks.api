package providers

type BillingUsageProvider interface {
	IncrementUsage(customerID string) error
}

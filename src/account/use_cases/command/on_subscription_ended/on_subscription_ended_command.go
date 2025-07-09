package on_subscription_ended

type OnSubscriptionEndedCommand struct {
	SubscriptionID    string
	CancelAtPeriodEnd bool
}

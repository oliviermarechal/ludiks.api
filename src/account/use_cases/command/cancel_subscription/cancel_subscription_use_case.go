package cancel_subscription

import (
	"fmt"
	domain_providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"time"
)

type CancelSubscriptionResult struct {
	EndDate string `json:"endDate"`
}

type CancelSubscriptionUseCase struct {
	stripeProvider                     domain_providers.StripeProvider
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository
	organizationRepository             domain_repositories.OrganizationRepository
}

func NewCancelSubscriptionUseCase(
	stripeProvider domain_providers.StripeProvider,
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository,
	organizationRepository domain_repositories.OrganizationRepository,
) *CancelSubscriptionUseCase {
	return &CancelSubscriptionUseCase{
		stripeProvider:                     stripeProvider,
		organizationSubscriptionRepository: organizationSubscriptionRepository,
		organizationRepository:             organizationRepository,
	}
}

func (c *CancelSubscriptionUseCase) Execute(command CancelSubscriptionCommand) (*CancelSubscriptionResult, error) {
	subscription, err := c.organizationSubscriptionRepository.FindByID(command.SubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %v", err)
	}

	if subscription.Status != "active" {
		return nil, fmt.Errorf("subscription is not active (current status: %s)", subscription.Status)
	}

	stripeResult, err := c.stripeProvider.CancelSubscription(subscription.StripeSubscriptionID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel subscription in Stripe: %v", err)
	}

	cancelledAt, err := time.Parse(time.RFC3339, stripeResult.CancelledAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cancelled at date: %v", err)
	}

	endDate, err := time.Parse(time.RFC3339, stripeResult.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end date: %v", err)
	}

	err = c.organizationSubscriptionRepository.CancelSubscription(
		command.SubscriptionID,
		&cancelledAt,
		&endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription in database: %v", err)
	}

	return &CancelSubscriptionResult{
		EndDate: stripeResult.EndDate,
	}, nil
}

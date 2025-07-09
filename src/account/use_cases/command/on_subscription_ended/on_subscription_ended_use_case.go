package on_subscription_ended

import (
	"ludiks/src/account/domain/models"
	domain_providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"time"
)

type OnSubscriptionEndedUseCase struct {
	stripeProvider                     domain_providers.StripeProvider
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository
	organizationRepository             domain_repositories.OrganizationRepository
}

func NewOnSubscriptionEndedUseCase(
	stripeProvider domain_providers.StripeProvider,
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository,
	organizationRepository domain_repositories.OrganizationRepository,
) *OnSubscriptionEndedUseCase {
	return &OnSubscriptionEndedUseCase{
		stripeProvider:                     stripeProvider,
		organizationSubscriptionRepository: organizationSubscriptionRepository,
		organizationRepository:             organizationRepository,
	}
}

func (o *OnSubscriptionEndedUseCase) Execute() error {
	subscriptions, err := o.organizationSubscriptionRepository.FindAllActiveWithPeriodEnd()
	if err != nil {
		return err
	}

	now := time.Now()

	for _, subscription := range subscriptions {
		if subscription.CurrentPeriodEnd != nil && now.After(*subscription.CurrentPeriodEnd) {
			err = o.cancelExpiredSubscription(subscription)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func (o *OnSubscriptionEndedUseCase) cancelExpiredSubscription(subscription models.OrganizationSubscription) error {
	err := o.organizationSubscriptionRepository.ExpireSubscription(subscription.ID)
	if err != nil {
		return err
	}

	err = o.organizationRepository.UpdatePlan(subscription.OrganizationID, "free")
	if err != nil {
		return err
	}

	return nil
}

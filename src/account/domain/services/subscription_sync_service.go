package services

import (
	"ludiks/src/account/domain/models"
	domain_providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"time"
)

type SubscriptionSyncService struct {
	stripeProvider                     domain_providers.StripeProvider
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository
	organizationRepository             domain_repositories.OrganizationRepository
}

func NewSubscriptionSyncService(
	stripeProvider domain_providers.StripeProvider,
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository,
	organizationRepository domain_repositories.OrganizationRepository,
) *SubscriptionSyncService {
	return &SubscriptionSyncService{
		stripeProvider:                     stripeProvider,
		organizationSubscriptionRepository: organizationSubscriptionRepository,
		organizationRepository:             organizationRepository,
	}
}

func (s *SubscriptionSyncService) SyncExpiredSubscriptions() error {
	subscriptions, err := s.organizationSubscriptionRepository.FindAllActiveWithPeriodEnd()
	if err != nil {
		return err
	}

	now := time.Now()

	for _, subscription := range subscriptions {
		if subscription.CurrentPeriodEnd != nil && now.After(*subscription.CurrentPeriodEnd) {
			err = s.cancelExpiredSubscription(subscription)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func (s *SubscriptionSyncService) cancelExpiredSubscription(subscription models.OrganizationSubscription) error {
	err := s.organizationSubscriptionRepository.ExpireSubscription(subscription.ID)
	if err != nil {
		return err
	}

	err = s.organizationRepository.UpdatePlan(subscription.OrganizationID, "free")
	if err != nil {
		return err
	}

	return nil
}

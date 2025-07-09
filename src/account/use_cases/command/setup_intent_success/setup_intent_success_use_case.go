package setup_intent_success

import (
	"fmt"
	"ludiks/src/account/domain/models"
	domain_providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"

	"github.com/google/uuid"
)

type SetupIntentSuccessUseCase struct {
	stripeProvider                     domain_providers.StripeProvider
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository
	organizationRepository             domain_repositories.OrganizationRepository
}

func NewSetupIntentSuccessUseCase(
	stripeProvider domain_providers.StripeProvider,
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository,
	organizationRepository domain_repositories.OrganizationRepository,
) *SetupIntentSuccessUseCase {
	return &SetupIntentSuccessUseCase{
		stripeProvider:                     stripeProvider,
		organizationSubscriptionRepository: organizationSubscriptionRepository,
		organizationRepository:             organizationRepository,
	}
}

func (s *SetupIntentSuccessUseCase) Execute(command SetupIntentSuccessCommand) error {
	organizationID, err := s.stripeProvider.GetOrganizationIDFromIntentPayment(command.SetupIntentID)
	if err != nil {
		return fmt.Errorf("failed to get organization ID from setup intent: %v", err)
	}

	organization, err := s.organizationRepository.Find(organizationID)
	if err != nil {
		return fmt.Errorf("failed to find organization: %v", err)
	}

	if organization.StripeCustomerID == nil {
		return fmt.Errorf("organization has no Stripe customer ID")
	}

	subscriptionID := uuid.New().String()
	subscription := models.CreateOrganizationSubscription(
		subscriptionID,
		organizationID,
		command.StripeSubscriptionID,
	)

	_, err = s.organizationSubscriptionRepository.Create(subscription)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %v", err)
	}

	err = s.organizationRepository.UpdatePlan(organizationID, "pro")
	if err != nil {
		return fmt.Errorf("failed to update organization plan: %v", err)
	}

	return nil
}

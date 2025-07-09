package subscriptions

import (
	"fmt"
	"ludiks/config"
	domain_providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/setup_intent_success"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetupIntentSuccessHandler struct {
	stripeProvider                     domain_providers.StripeProvider
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository
	organizationRepository             domain_repositories.OrganizationRepository
}

func NewSetupIntentSuccessHandler(
	stripeProvider domain_providers.StripeProvider,
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository,
	organizationRepository domain_repositories.OrganizationRepository,
) *SetupIntentSuccessHandler {
	return &SetupIntentSuccessHandler{
		stripeProvider:                     stripeProvider,
		organizationSubscriptionRepository: organizationSubscriptionRepository,
		organizationRepository:             organizationRepository,
	}
}

func (h *SetupIntentSuccessHandler) Handle(c *gin.Context) {
	setupIntentID := c.Query("setup_intent")
	stripeSubscriptionID := c.Query("subscription_id")

	command := setup_intent_success.SetupIntentSuccessCommand{
		SetupIntentID:        setupIntentID,
		StripeSubscriptionID: stripeSubscriptionID,
	}

	useCase := setup_intent_success.NewSetupIntentSuccessUseCase(
		h.stripeProvider,
		h.organizationSubscriptionRepository,
		h.organizationRepository,
	)
	err := useCase.Execute(command)
	frontURL := config.AppConfig.FrontURL
	if err != nil {
		url := fmt.Sprintf("%s/dashboard/organization/billing?intent_status=error&intent_payment_error=%s", frontURL, err.Error())

		c.Redirect(http.StatusSeeOther, url)

		return
	}

	url := fmt.Sprintf("%s/dashboard/organization/billing?intent_status=success", frontURL)
	c.Redirect(http.StatusSeeOther, url)
}

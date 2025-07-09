package subscriptions

import (
	domain_providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/cancel_subscription"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CancelSubscriptionHandler struct {
	stripeProvider                     domain_providers.StripeProvider
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository
	organizationRepository             domain_repositories.OrganizationRepository
}

func NewCancelSubscriptionHandler(
	stripeProvider domain_providers.StripeProvider,
	organizationSubscriptionRepository domain_repositories.OrganizationSubscriptionRepository,
	organizationRepository domain_repositories.OrganizationRepository,
) *CancelSubscriptionHandler {
	return &CancelSubscriptionHandler{
		stripeProvider:                     stripeProvider,
		organizationSubscriptionRepository: organizationSubscriptionRepository,
		organizationRepository:             organizationRepository,
	}
}

func (h *CancelSubscriptionHandler) Handle(c *gin.Context) {
	subscriptionID := c.Param("id")
	result, err := cancel_subscription.NewCancelSubscriptionUseCase(
		h.stripeProvider,
		h.organizationSubscriptionRepository,
		h.organizationRepository,
	).Execute(cancel_subscription.CancelSubscriptionCommand{SubscriptionID: subscriptionID})
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

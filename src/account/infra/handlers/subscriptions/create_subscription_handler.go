package subscriptions

import (
	domain_providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/create_subscription"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateSubscriptionHandler struct {
	stripeProvider         domain_providers.StripeProvider
	organizationRepository domain_repositories.OrganizationRepository
}

func NewCreateSubscriptionHandler(
	stripeProvider domain_providers.StripeProvider,
	organizationRepository domain_repositories.OrganizationRepository,
) *CreateSubscriptionHandler {
	return &CreateSubscriptionHandler{
		stripeProvider:         stripeProvider,
		organizationRepository: organizationRepository,
	}
}

func (h *CreateSubscriptionHandler) Handle(c *gin.Context) {
	var dto create_subscription.CreateSubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := dto.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	command := create_subscription.CreateSubscriptionCommand{
		CustomerEmail: dto.CustomerEmail,
		CompanyName:   dto.CompanyName,
		CompanyAddress: create_subscription.CreateSubscriptionCommandCompanyAddress{
			Line1:      dto.CompanyAddress.Line1,
			Line2:      dto.CompanyAddress.Line2,
			City:       dto.CompanyAddress.City,
			State:      dto.CompanyAddress.State,
			PostalCode: dto.CompanyAddress.PostalCode,
			Country:    dto.CompanyAddress.Country,
		},
		ContactName:    dto.ContactName,
		ContactPhone:   dto.ContactPhone,
		TaxIDValue:     dto.TaxIDValue,
		OrganizationID: dto.OrganizationID,
	}

	useCase := create_subscription.NewCreateSubscriptionUseCase(h.stripeProvider, h.organizationRepository)
	result, err := useCase.Execute(command)
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

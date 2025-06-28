package organizations

import (
	"errors"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/invit_membership"
	"ludiks/src/kernel/app/handlers"
	"ludiks/src/kernel/domain/providers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InvitMembershipHandler struct {
	organizationRepository domain_repositories.OrganizationRepository
	userRepository         domain_repositories.UserRepository
	mailerProvider         providers.MailerProvider
}

func NewInvitMembershipHandler(
	organizationRepository domain_repositories.OrganizationRepository,
	userRepository domain_repositories.UserRepository,
	mailerProvider providers.MailerProvider,
) *InvitMembershipHandler {
	return &InvitMembershipHandler{
		organizationRepository: organizationRepository,
		userRepository:         userRepository,
		mailerProvider:         mailerProvider,
	}
}

func (h *InvitMembershipHandler) Handle(c *gin.Context) {
	organizationID := c.Param("id")

	var dto invit_membership.InvitMembershipDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := dto.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	user_id, ok := c.Get("user_id")
	if !ok {
		handlers.HandleBadRequest(c, errors.New("authentication required"))
		return
	}

	organizationInvitation, err := invit_membership.NewInvitMembershipUseCase(
		h.organizationRepository,
		h.userRepository,
		h.mailerProvider,
	).Execute(
		invit_membership.InvitMembershipCommand{
			OrganizationID: organizationID,
			FromUserID:     user_id.(string),
			Email:          dto.Email,
			Role:           dto.Role,
		},
	)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, organizationInvitation)
}

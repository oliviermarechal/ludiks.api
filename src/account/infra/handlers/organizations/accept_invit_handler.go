package organizations

import (
	"errors"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/accept_invit"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AcceptInvitHandler struct {
	organizationRepository domain_repositories.OrganizationRepository
	userRepository         domain_repositories.UserRepository
}

func NewAcceptInvitHandler(
	organizationRepository domain_repositories.OrganizationRepository,
	userRepository domain_repositories.UserRepository,
) *AcceptInvitHandler {
	return &AcceptInvitHandler{
		organizationRepository: organizationRepository,
		userRepository:         userRepository,
	}
}

func (h *AcceptInvitHandler) Handle(c *gin.Context) {
	invitID := c.Param("invit-id")
	userID, ok := c.Get("user_id")
	if !ok {
		handlers.HandleBadRequest(c, errors.New("authentication required"))
		return
	}

	err := accept_invit.NewAcceptInvitUseCase(
		h.organizationRepository,
		h.userRepository,
	).Execute(
		accept_invit.AcceptInvitCommand{
			InvitID: invitID,
			UserID:  userID.(string),
		},
	)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

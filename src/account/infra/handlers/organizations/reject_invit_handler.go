package organizations

import (
	"errors"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/reject_invit"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RejectInvitHandler struct {
	organizationRepository domain_repositories.OrganizationRepository
	userRepository         domain_repositories.UserRepository
}

func NewRejectInvitHandler(
	organizationRepository domain_repositories.OrganizationRepository,
	userRepository domain_repositories.UserRepository,
) *RejectInvitHandler {
	return &RejectInvitHandler{
		organizationRepository: organizationRepository,
		userRepository:         userRepository,
	}
}

func (h *RejectInvitHandler) Handle(c *gin.Context) {
	invitID := c.Param("invit-id")
	userID, ok := c.Get("user_id")
	if !ok {
		handlers.HandleBadRequest(c, errors.New("authentication required"))
		return
	}

	err := reject_invit.NewRejectInvitUseCase(
		h.organizationRepository,
		h.userRepository,
	).Execute(
		reject_invit.RejectInvitCommand{
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

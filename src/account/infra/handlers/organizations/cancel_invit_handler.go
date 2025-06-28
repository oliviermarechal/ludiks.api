package organizations

import (
	"errors"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/cancel_invit"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CancelInvitHandler struct {
	organizationRepository domain_repositories.OrganizationRepository
}

func NewCancelInvitHandler(
	organizationRepository domain_repositories.OrganizationRepository,
) *CancelInvitHandler {
	return &CancelInvitHandler{
		organizationRepository: organizationRepository,
	}
}

func (h *CancelInvitHandler) Handle(c *gin.Context) {
	invitID := c.Param("invit-id")
	userID, ok := c.Get("user_id")
	if !ok {
		handlers.HandleBadRequest(c, errors.New("authentication required"))
		return
	}

	err := cancel_invit.NewCancelInvitUseCase(
		h.organizationRepository,
	).Execute(
		cancel_invit.CancelInvitCommand{
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

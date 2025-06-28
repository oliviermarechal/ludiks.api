package organizations

import (
	"errors"
	"ludiks/src/account/use_cases/query"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListUserReceivedInvitesHandler struct {
	db *gorm.DB
}

func NewListUserReceivedInvitesHandler(
	db *gorm.DB,
) *ListUserReceivedInvitesHandler {
	return &ListUserReceivedInvitesHandler{
		db: db,
	}
}

func (h *ListUserReceivedInvitesHandler) Handle(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handlers.HandleUnauthorized(c, errors.New("user not authenticated"))
		return
	}

	invites := query.ListUserOrganizationInvitesQuery(h.db, userID.(string))

	c.JSON(http.StatusOK, invites)
}

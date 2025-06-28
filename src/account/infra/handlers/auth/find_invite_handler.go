package auth

import (
	"ludiks/src/account/use_cases/query"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FindInviteHandler struct {
	db *gorm.DB
}

func NewFindInviteHandler(
	db *gorm.DB,
) *FindInviteHandler {
	return &FindInviteHandler{
		db: db,
	}
}

func (h *FindInviteHandler) Handle(c *gin.Context) {
	inviteID := c.Param("id")
	loginResult, err := query.FindInviteQuery(h.db, inviteID)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResult)
}

package auth

import (
	"errors"
	"ludiks/src/account/use_cases/query"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeHandler struct {
	db *gorm.DB
}

func NewMeHandler(db *gorm.DB) *MeHandler {
	return &MeHandler{
		db: db,
	}
}

func (h *MeHandler) Handle(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handlers.HandleUnauthorized(c, errors.New("user not authenticated"))
		return
	}

	user, err := query.Me(h.db, userID.(string))
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, user)

}

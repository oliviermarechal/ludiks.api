package organizations

import (
	"errors"
	"ludiks/src/account/use_cases/query"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListOrganizationsHandler struct {
	db *gorm.DB
}

func NewListOrganizationsHandler(
	db *gorm.DB,
) *ListOrganizationsHandler {
	return &ListOrganizationsHandler{
		db: db,
	}
}

func (h *ListOrganizationsHandler) Handle(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handlers.HandleBadRequest(c, errors.New("user_id not found"))
		return
	}

	projects := query.ListUserOrganizationsQuery(h.db, userID.(string))

	c.JSON(http.StatusOK, projects)
}

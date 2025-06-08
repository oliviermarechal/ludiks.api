package projects

import (
	"errors"
	"ludiks/src/account/use_cases/query"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListProjectsHandler struct {
	db *gorm.DB
}

func NewListProjectsHandler(
	db *gorm.DB,
) *ListProjectsHandler {
	return &ListProjectsHandler{
		db: db,
	}
}

func (h *ListProjectsHandler) Handle(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		handlers.HandleBadRequest(c, errors.New("user_id not found"))
		return
	}

	projects := query.ListUserProjectsQuery(h.db, userID.(string))

	c.JSON(http.StatusOK, projects)
}

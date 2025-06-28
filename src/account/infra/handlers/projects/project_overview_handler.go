package projects

import (
	"ludiks/src/account/use_cases/query"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectOverviewHandler struct {
	db *gorm.DB
}

func NewProjectOverviewHandler(
	db *gorm.DB,
) *ProjectOverviewHandler {
	return &ProjectOverviewHandler{
		db: db,
	}
}

func (h *ProjectOverviewHandler) Handle(c *gin.Context) {
	projectID := c.Param("id")

	overview, err := query.ProjectOverviewQuery(h.db, projectID)
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, overview)
}

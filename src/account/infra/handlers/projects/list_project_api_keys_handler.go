package projects

import (
	"ludiks/src/account/use_cases/query"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListProjectApiKeysHandler struct {
	db *gorm.DB
}

func NewListProjectApiKeysHandler(
	db *gorm.DB,
) *ListProjectApiKeysHandler {
	return &ListProjectApiKeysHandler{
		db: db,
	}
}

func (h *ListProjectApiKeysHandler) Handle(c *gin.Context) {
	projectID := c.Param("id")

	apiKeys := query.ListProjectApiKeys(h.db, projectID)

	c.JSON(http.StatusOK, apiKeys)
}

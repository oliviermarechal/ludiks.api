package projects

import (
	"ludiks/src/account/use_cases/query"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListProjectMetadataHandler struct {
	db *gorm.DB
}

func NewListProjectMetadataHandler(
	db *gorm.DB,
) *ListProjectMetadataHandler {
	return &ListProjectMetadataHandler{
		db: db,
	}
}

func (h *ListProjectMetadataHandler) Handle(c *gin.Context) {
	projectID := c.Param("id")

	metadata, err := query.ListProjectMetadataQuery(h.db, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metadata)
}

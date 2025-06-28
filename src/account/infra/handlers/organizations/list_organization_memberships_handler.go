package organizations

import (
	"ludiks/src/account/use_cases/query"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListOrganizationMembershipsHandler struct {
	db *gorm.DB
}

func NewListOrganizationMembershipsHandler(
	db *gorm.DB,
) *ListOrganizationMembershipsHandler {
	return &ListOrganizationMembershipsHandler{
		db: db,
	}
}

func (h *ListOrganizationMembershipsHandler) Handle(c *gin.Context) {
	organizationID := c.Param("id")

	memberships := query.ListOrganizationMembership(h.db, organizationID)

	c.JSON(http.StatusOK, memberships)
}

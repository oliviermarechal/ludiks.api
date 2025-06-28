package organizations

import (
	"ludiks/src/account/use_cases/query"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListInvitesHandler struct {
	db *gorm.DB
}

func NewListInvitesHandler(
	db *gorm.DB,
) *ListInvitesHandler {
	return &ListInvitesHandler{
		db: db,
	}
}

func (h *ListInvitesHandler) Handle(c *gin.Context) {
	organizationID := c.Param("id")

	memberships := query.ListInvitesQuery(h.db, organizationID)

	c.JSON(http.StatusOK, memberships)
}

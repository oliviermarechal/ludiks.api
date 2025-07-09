package subscriptions

import (
	"ludiks/src/account/use_cases/query"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetOrganizationSubscriptionHandler struct {
	db *gorm.DB
}

func NewGetOrganizationSubscriptionHandler(db *gorm.DB) *GetOrganizationSubscriptionHandler {
	return &GetOrganizationSubscriptionHandler{db: db}
}

func (h *GetOrganizationSubscriptionHandler) Handle(c *gin.Context) {
	organizationID := c.Param("id")
	result := query.GetOrganizationSubscription(h.db, organizationID)

	c.JSON(http.StatusOK, result)
}

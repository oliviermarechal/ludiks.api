package circuit_handler

import (
	"ludiks/src/gamification/use_cases/query"
	"ludiks/src/kernel/handlers"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ListRewardHandler struct {
	db *gorm.DB
}

func NewListRewardHandler(db *gorm.DB) *ListRewardHandler {
	return &ListRewardHandler{db: db}
}

func (h *ListRewardHandler) Handle(c *gin.Context) {
	circuitID := c.Param("id")

	rewards, err := query.ListReward(h.db, circuitID)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, rewards)
}

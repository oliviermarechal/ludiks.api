package circuit_handler

import (
	"errors"
	"ludiks/src/gamification/use_cases/query"
	"ludiks/src/kernel/app/handlers"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetCircuitInsightHandler struct {
	db *gorm.DB
}

func NewGetCircuitInsightHandler(db *gorm.DB) *GetCircuitInsightHandler {
	return &GetCircuitInsightHandler{db: db}
}

func (h *GetCircuitInsightHandler) Handle(c *gin.Context) {
	circuitID := c.Param("id")
	if circuitID == "" {
		handlers.HandleBadRequest(c, errors.New("circuit is required"))
		return
	}

	circuitInsight, err := query.GetCircuitInsight(h.db, circuitID)
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, circuitInsight)
}

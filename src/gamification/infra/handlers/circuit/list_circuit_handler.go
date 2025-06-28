package circuit_handler

import (
	"errors"
	"ludiks/src/gamification/use_cases/query"
	"ludiks/src/kernel/app/handlers"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListCircuitHandler struct {
	db *gorm.DB
}

func NewListCircuitHandler(db *gorm.DB) *ListCircuitHandler {
	return &ListCircuitHandler{db: db}
}

func (h *ListCircuitHandler) Handle(c *gin.Context) {
	projectID := c.Query("projectId")
	if projectID == "" {
		handlers.HandleBadRequest(c, errors.New("project_id is required"))
		return
	}

	circuits, err := query.ListCircuit(h.db, projectID)
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, circuits)
}

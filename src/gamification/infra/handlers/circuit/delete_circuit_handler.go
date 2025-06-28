package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/delete_circuit"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteCircuitHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewDeleteCircuitHandler(circuitRepository domain_repositories.CircuitRepository) *DeleteCircuitHandler {
	return &DeleteCircuitHandler{circuitRepository: circuitRepository}
}

func (h *DeleteCircuitHandler) Handle(c *gin.Context) {
	circuitID := c.Param("id")

	deleteCircuitHandler := delete_circuit.NewDeleteCircuitUseCase(h.circuitRepository)
	err := deleteCircuitHandler.Execute(&delete_circuit.DeleteCircuitCommand{
		CircuitID: circuitID,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

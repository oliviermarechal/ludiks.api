package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/activate_circuit"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ActivateCircuitHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewActivateCircuitHandler(circuitRepository domain_repositories.CircuitRepository) *ActivateCircuitHandler {
	return &ActivateCircuitHandler{circuitRepository: circuitRepository}
}

func (h *ActivateCircuitHandler) Handle(c *gin.Context) {
	circuitID := c.Param("id")

	activateCircuitUseCase := activate_circuit.NewActivateCircuitUseCase(h.circuitRepository)
	err := activateCircuitUseCase.Execute(&activate_circuit.ActivateCircuitCommand{
		CircuitID: circuitID,
	})
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

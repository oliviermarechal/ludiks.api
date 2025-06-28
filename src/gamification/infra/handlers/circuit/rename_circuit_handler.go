package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/rename_circuit"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RenameCircuitHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewRenameCircuitHandler(circuitRepository domain_repositories.CircuitRepository) *RenameCircuitHandler {
	return &RenameCircuitHandler{circuitRepository: circuitRepository}
}

func (h *RenameCircuitHandler) Handle(c *gin.Context) {
	var renameCircuitDTO rename_circuit.RenameCircuitDTO
	if err := c.ShouldBindJSON(&renameCircuitDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := renameCircuitDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	circuitID := c.Param("id")

	renameCircuitUseCase := rename_circuit.NewRenameCircuitUseCase(h.circuitRepository)
	circuit, err := renameCircuitUseCase.Execute(&rename_circuit.RenameCircuitCommand{
		CircuitID: circuitID,
		Name:      renameCircuitDTO.Name,
	})
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, circuit)
}

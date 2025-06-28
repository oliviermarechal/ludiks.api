package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/create_circuit"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateCircuitHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewCreateCircuitHandler(circuitRepository domain_repositories.CircuitRepository) *CreateCircuitHandler {
	return &CreateCircuitHandler{circuitRepository: circuitRepository}
}

func (h *CreateCircuitHandler) Handle(c *gin.Context) {
	var createCircuitDTO create_circuit.CreateCircuitDTO
	if err := c.ShouldBindJSON(&createCircuitDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := createCircuitDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	createCircuitUseCase := create_circuit.NewCreateCircuitUseCase(h.circuitRepository)
	circuit, err := createCircuitUseCase.Execute(&create_circuit.CreateCircuitCommand{
		ProjectID: createCircuitDTO.ProjectID,
		Name:      createCircuitDTO.Name,
		Type:      createCircuitDTO.Type,
	})
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, circuit)
}

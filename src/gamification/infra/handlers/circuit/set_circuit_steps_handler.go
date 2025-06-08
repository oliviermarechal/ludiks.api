package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/set_circuit_steps"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SetCircuitStepsHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewSetCircuitStepsHandler(circuitRepository domain_repositories.CircuitRepository) *SetCircuitStepsHandler {
	return &SetCircuitStepsHandler{circuitRepository: circuitRepository}
}

func (h *SetCircuitStepsHandler) Handle(c *gin.Context) {
	circuitID := c.Param("id")
	var setCircuitStepsDTO set_circuit_steps.SetCircuitStepsDTO
	if err := c.ShouldBindJSON(&setCircuitStepsDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := setCircuitStepsDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	setCircuitStepsUseCase := set_circuit_steps.NewSetCircuitStepsUseCase(h.circuitRepository)

	steps := make([]struct {
		Name                string
		Description         *string
		CompletionThreshold int
		EventName           string
	}, len(setCircuitStepsDTO.Steps))

	for i, step := range setCircuitStepsDTO.Steps {
		steps[i] = struct {
			Name                string
			Description         *string
			CompletionThreshold int
			EventName           string
		}{
			Name:                step.Name,
			Description:         step.Description,
			CompletionThreshold: step.CompletionThreshold,
			EventName:           step.EventName,
		}
	}

	circuit, err := setCircuitStepsUseCase.Execute(&set_circuit_steps.SetCircuitStepsCommand{
		CircuitID: circuitID,
		Steps:     steps,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, circuit)
}

package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/generate_step"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenerateStepHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewGenerateStepHandler(circuitRepository domain_repositories.CircuitRepository) *GenerateStepHandler {
	return &GenerateStepHandler{circuitRepository: circuitRepository}
}

func (h *GenerateStepHandler) Handle(c *gin.Context) {
	circuitID := c.Param("id")
	var generateStepDTO generate_step.GenerateStepDTO
	if err := c.ShouldBindJSON(&generateStepDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := generateStepDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	generateStepUseCase := generate_step.NewGenerateStepUseCase(h.circuitRepository)
	circuit, err := generateStepUseCase.Execute(&generate_step.GenerateStepCommand{
		CircuitID:     circuitID,
		CircuitType:   generateStepDTO.CircuitType,
		NumberOfSteps: generateStepDTO.NumberOfSteps,
		Curve:         generateStepDTO.Curve,
		MaxValue:      generateStepDTO.MaxValue,
		Exponent:      generateStepDTO.Exponent,
		EventName:     generateStepDTO.EventName,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, circuit)
}

package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/reordering_steps"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReorderingStepsHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewReorderingStepsHandler(circuitRepository domain_repositories.CircuitRepository) *ReorderingStepsHandler {
	return &ReorderingStepsHandler{circuitRepository: circuitRepository}
}

func (h *ReorderingStepsHandler) Handle(c *gin.Context) {
	var reorderingStepsDTO reordering_steps.ReorderingStepsDTO
	if err := c.ShouldBindJSON(&reorderingStepsDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := reorderingStepsDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	steps := make([]struct {
		StepId     string
		StepNumber int
	}, len(reorderingStepsDTO.Steps))
	for i, step := range reorderingStepsDTO.Steps {
		steps[i] = struct {
			StepId     string
			StepNumber int
		}{
			StepId:     step.StepID,
			StepNumber: step.StepNumber,
		}
	}

	reorderingStepsUseCase := reordering_steps.NewReorderingStepsUseCase(h.circuitRepository)
	circuit, err := reorderingStepsUseCase.Execute(&reordering_steps.ReorderingStepsCommand{
		CircuitID: c.Param("id"),
		Steps:     steps,
	})
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, circuit)
}

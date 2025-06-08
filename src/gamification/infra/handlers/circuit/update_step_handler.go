package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/update_step"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateStepHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewUpdateStepHandler(circuitRepository domain_repositories.CircuitRepository) *UpdateStepHandler {
	return &UpdateStepHandler{circuitRepository: circuitRepository}
}

func (h *UpdateStepHandler) Handle(c *gin.Context) {
	stepID := c.Param("stepId")
	var updateStepDTO update_step.UpdateStepDTO
	if err := c.ShouldBindJSON(&updateStepDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := updateStepDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	updateStepUseCase := update_step.NewUpdateStepUseCase(h.circuitRepository)
	step, err := updateStepUseCase.Execute(&update_step.UpdateStepCommand{
		StepId:              stepID,
		Name:                updateStepDTO.Name,
		Description:         updateStepDTO.Description,
		CompletionThreshold: updateStepDTO.CompletionThreshold,
		EventName:           updateStepDTO.EventName,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, step)
}

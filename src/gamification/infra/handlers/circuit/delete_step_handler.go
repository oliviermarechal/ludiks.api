package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/delete_step"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteStepHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewDeleteStepHandler(circuitRepository domain_repositories.CircuitRepository) *DeleteStepHandler {
	return &DeleteStepHandler{circuitRepository: circuitRepository}
}

func (h *DeleteStepHandler) Handle(c *gin.Context) {
	stepID := c.Param("stepId")

	deleteStepUseCase := delete_step.NewDeleteStepUseCase(h.circuitRepository)
	err := deleteStepUseCase.Execute(&delete_step.DeleteStepCommand{
		StepId: stepID,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

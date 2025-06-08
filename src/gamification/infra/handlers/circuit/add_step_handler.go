package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/circuit/create_step"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddStepHandler struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewAddStepHandler(circuitRepository domain_repositories.CircuitRepository) *AddStepHandler {
	return &AddStepHandler{circuitRepository: circuitRepository}
}

func (h *AddStepHandler) Handle(c *gin.Context) {
	circuitID := c.Param("id")
	var addStepDTO create_step.CreateStepDTO
	if err := c.ShouldBindJSON(&addStepDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := addStepDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	createStepUseCase := create_step.NewCreateStepUseCase(h.circuitRepository)
	step, err := createStepUseCase.Execute(&create_step.CreateStepCommand{
		CircuitID:           circuitID,
		Name:                addStepDTO.Name,
		Description:         addStepDTO.Description,
		CompletionThreshold: addStepDTO.CompletionThreshold,
		EventName:           addStepDTO.EventName,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, step)
}

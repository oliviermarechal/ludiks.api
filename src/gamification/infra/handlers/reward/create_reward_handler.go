package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/reward/create_reward"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateRewardHandler struct {
	rewardRepository domain_repositories.RewardRepository
}

func NewCreateRewardHandler(rewardRepository domain_repositories.RewardRepository) *CreateRewardHandler {
	return &CreateRewardHandler{rewardRepository: rewardRepository}
}

func (h *CreateRewardHandler) Handle(c *gin.Context) {
	var createRewardDTO create_reward.CreateRewardDTO
	if err := c.ShouldBindJSON(&createRewardDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := createRewardDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	circuitID := c.Param("id")

	createRewardUseCase := create_reward.NewCreateRewardUseCase(h.rewardRepository)
	reward, err := createRewardUseCase.Execute(&create_reward.CreateRewardCommand{
		CircuitID:                 circuitID,
		Name:                      createRewardDTO.Name,
		Description:               createRewardDTO.Description,
		StepID:                    createRewardDTO.StepID,
		UnlockOnCircuitCompletion: createRewardDTO.UnlockOnCircuitCompletion,
	})
	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, reward)
}

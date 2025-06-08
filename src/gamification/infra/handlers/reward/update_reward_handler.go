package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/reward/update_reward"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateRewardHandler struct {
	rewardRepository domain_repositories.RewardRepository
}

func NewUpdateRewardHandler(rewardRepository domain_repositories.RewardRepository) *UpdateRewardHandler {
	return &UpdateRewardHandler{rewardRepository: rewardRepository}
}

func (h *UpdateRewardHandler) Handle(c *gin.Context) {
	var updateRewardDTO update_reward.UpdateRewardDTO
	if err := c.ShouldBindJSON(&updateRewardDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := updateRewardDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	rewardID := c.Param("rewardId")

	updateRewardUseCase := update_reward.NewUpdateRewardUseCase(h.rewardRepository)
	reward, err := updateRewardUseCase.Execute(&update_reward.UpdateRewardCommand{
		RewardID:                  rewardID,
		Name:                      updateRewardDTO.Name,
		Description:               updateRewardDTO.Description,
		StepID:                    updateRewardDTO.StepID,
		UnlockOnCircuitCompletion: updateRewardDTO.UnlockOnCircuitCompletion,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, reward)
}

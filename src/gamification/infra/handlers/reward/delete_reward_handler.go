package circuit_handler

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
	"ludiks/src/gamification/use_cases/command/reward/delete_reward"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteRewardHandler struct {
	rewardRepository domain_repositories.RewardRepository
}

func NewDeleteRewardHandler(rewardRepository domain_repositories.RewardRepository) *DeleteRewardHandler {
	return &DeleteRewardHandler{rewardRepository: rewardRepository}
}

func (h *DeleteRewardHandler) Handle(c *gin.Context) {
	rewardID := c.Param("rewardId")

	deleteRewardUseCase := delete_reward.NewDeleteRewardUseCase(h.rewardRepository)
	err := deleteRewardUseCase.Execute(&delete_reward.DeleteRewardCommand{
		RewardID: rewardID,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

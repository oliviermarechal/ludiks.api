package create_reward

import (
	"ludiks/src/gamification/domain/models"
	domain_repositories "ludiks/src/gamification/domain/repositories"

	"github.com/google/uuid"
)

type CreateRewardUseCase struct {
	rewardRepository domain_repositories.RewardRepository
}

func NewCreateRewardUseCase(rewardRepository domain_repositories.RewardRepository) *CreateRewardUseCase {
	return &CreateRewardUseCase{rewardRepository: rewardRepository}
}

func (u *CreateRewardUseCase) Execute(command *CreateRewardCommand) (*models.Reward, error) {
	return u.rewardRepository.Create(models.CreateReward(
		uuid.New().String(),
		command.Name,
		command.Description,
		command.CircuitID,
		command.StepID,
		command.UnlockOnCircuitCompletion,
	))
}

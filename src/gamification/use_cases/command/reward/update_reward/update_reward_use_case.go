package update_reward

import (
	"ludiks/src/gamification/domain/models"
	domain_repositories "ludiks/src/gamification/domain/repositories"
)

type UpdateRewardUseCase struct {
	rewardRepository domain_repositories.RewardRepository
}

func NewUpdateRewardUseCase(rewardRepository domain_repositories.RewardRepository) *UpdateRewardUseCase {
	return &UpdateRewardUseCase{rewardRepository: rewardRepository}
}

func (u *UpdateRewardUseCase) Execute(command *UpdateRewardCommand) (*models.Reward, error) {
	reward, err := u.rewardRepository.Find(command.RewardID)
	if err != nil {
		return nil, err
	}

	reward.Name = command.Name
	reward.Description = command.Description
	reward.StepID = command.StepID
	reward.UnlockOnCircuitCompletion = command.UnlockOnCircuitCompletion

	return u.rewardRepository.Update(command.RewardID, reward)
}

package delete_reward

import (
	domain_repositories "ludiks/src/gamification/domain/repositories"
)

type DeleteRewardUseCase struct {
	rewardRepository domain_repositories.RewardRepository
}

func NewDeleteRewardUseCase(rewardRepository domain_repositories.RewardRepository) *DeleteRewardUseCase {
	return &DeleteRewardUseCase{rewardRepository: rewardRepository}
}

func (u *DeleteRewardUseCase) Execute(command *DeleteRewardCommand) error {
	return u.rewardRepository.Delete(command.RewardID)
}

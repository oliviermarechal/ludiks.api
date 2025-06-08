package domain_repositories

import (
	"ludiks/src/gamification/domain/models"
)

type RewardRepository interface {
	Create(reward *models.Reward) (*models.Reward, error)
	Find(id string) (*models.Reward, error)
	Delete(id string) error
	Update(id string, reward *models.Reward) (*models.Reward, error)
}

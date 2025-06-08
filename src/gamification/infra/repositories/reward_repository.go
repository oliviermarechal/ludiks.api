package infra_repositories

import (
	"ludiks/src/gamification/domain/models"

	"gorm.io/gorm"
)

type RewardRepository struct {
	db *gorm.DB
}

func NewRewardRepository(db *gorm.DB) *RewardRepository {
	return &RewardRepository{db: db}
}

func (r *RewardRepository) Create(reward *models.Reward) (*models.Reward, error) {
	err := r.db.Create(reward).Error
	if err != nil {
		return nil, err
	}

	return reward, nil
}

func (r *RewardRepository) Find(id string) (*models.Reward, error) {
	var reward models.Reward
	err := r.db.Where("id = ?", id).First(&reward).Error
	if err != nil {
		return nil, err
	}

	return &reward, nil
}

func (r *RewardRepository) Update(id string, reward *models.Reward) (*models.Reward, error) {
	err := r.db.Model(&models.Reward{}).Where("id = ?", id).Updates(reward).Error
	if err != nil {
		return nil, err
	}

	return r.Find(id)
}

func (r *RewardRepository) Delete(id string) error {
	err := r.db.Model(&models.Reward{}).Where("id = ?", id).Delete(&models.Reward{}).Error
	if err != nil {
		return err
	}

	return nil
}

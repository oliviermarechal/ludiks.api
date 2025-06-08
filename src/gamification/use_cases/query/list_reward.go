package query

import (
	"ludiks/src/gamification/domain/models"

	"gorm.io/gorm"
)

func ListReward(db *gorm.DB, circuitID string) ([]*models.Reward, error) {
	var rewards []*models.Reward
	if err := db.Where("circuit_id = ?", circuitID).Find(&rewards).Error; err != nil {
		return nil, err
	}

	return rewards, nil
}

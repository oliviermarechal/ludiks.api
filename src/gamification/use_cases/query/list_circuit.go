package query

import (
	"ludiks/src/gamification/domain/models"

	"gorm.io/gorm"
)

func ListCircuit(db *gorm.DB, projectID string) ([]*models.Circuit, error) {
	var circuits []*models.Circuit
	if err := db.Where("project_id = ?", projectID).Preload("Steps").Find(&circuits).Error; err != nil {
		return nil, err
	}

	return circuits, nil
}

package query

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

func ListProjectApiKeys(db *gorm.DB, projectID string) []models.ApiKey {
	var project models.Project
	if err := db.Preload("ApiKeys").Find(&project, "id = ?", projectID).Error; err != nil {
		return []models.ApiKey{}
	}

	return *project.ApiKeys
}

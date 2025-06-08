package query

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

func ListUserProjectsQuery(db *gorm.DB, userID string) []models.Project {
	var user models.User
	if err := db.Preload("Projects").First(&user, "id = ?", userID).Error; err != nil {
		return []models.Project{}
	}

	return user.Projects
}

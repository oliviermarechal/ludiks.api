package query

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

func Me(db *gorm.DB, userID string) (models.User, error) {
	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

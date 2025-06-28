package query

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

func FindInviteQuery(db *gorm.DB, inviteID string) (*models.OrganizationInvitation, error) {
	var invite *models.OrganizationInvitation

	err := db.Preload("Organization").Where("id = ?", inviteID).First(&invite).Error

	return invite, err
}

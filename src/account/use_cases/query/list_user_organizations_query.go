package query

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

func ListUserOrganizationsQuery(db *gorm.DB, userID string) []models.Organization {
	var organizations []models.Organization
	if err := db.
		Preload("Projects").
		Preload("Memberships").
		Model(&models.Organization{}).
		Joins("JOIN organization_memberships ON organization_memberships.organization_id = organizations.id").
		Where("organization_memberships.user_id = ?", userID).
		Find(&organizations).
		Error; err != nil {
		return []models.Organization{}
	}

	return organizations
}

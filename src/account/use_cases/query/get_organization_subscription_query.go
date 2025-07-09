package query

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

func GetOrganizationSubscription(db *gorm.DB, organizationID string) *models.OrganizationSubscription {
	var subscription models.OrganizationSubscription

	if err := db.Where("organization_id = ? AND status = ?", organizationID, "active").First(&subscription).Error; err != nil {
		return nil
	}

	return &subscription
}

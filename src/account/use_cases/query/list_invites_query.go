package query

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

func ListInvitesQuery(db *gorm.DB, organizationID string) []*models.OrganizationInvitation {
	var invites []*models.OrganizationInvitation

	if err := db.Where("organization_id = ?", organizationID).Find(&invites).Error; err != nil {
		return []*models.OrganizationInvitation{}
	}

	return invites
}

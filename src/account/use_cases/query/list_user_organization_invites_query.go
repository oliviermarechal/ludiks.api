package query

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

func ListUserOrganizationInvitesQuery(db *gorm.DB, userID string) []models.OrganizationInvitation {
	var user models.User
	db.Select("email").Where("id = ?", userID).First(&user)

	var invites []models.OrganizationInvitation
	if err := db.
		Preload("Organization").
		Preload("FromUser").
		Model(&models.OrganizationInvitation{}).
		Where("organization_invitations.to_email = ?", user.Email).
		Find(&invites).
		Error; err != nil {
		return []models.OrganizationInvitation{}
	}

	return invites
}

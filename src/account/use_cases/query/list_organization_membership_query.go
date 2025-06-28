package query

import (
	"time"

	"gorm.io/gorm"
)

type OrganizationMembership struct {
	ID       string    `json:"id" gorm:"column:id"`
	UserID   string    `json:"userId" gorm:"column:user_id"`
	Email    string    `json:"email" gorm:"column:email"`
	Role     string    `json:"role" gorm:"column:role"`
	JoinedAt time.Time `json:"joinedAt" gorm:"column:joined_at"`
}

func ListOrganizationMembership(db *gorm.DB, organizationID string) []*OrganizationMembership {
	var memberships []*OrganizationMembership

	db.Table("organization_memberships om").
		Select(`
			om.id,
			om.created_at as "joined_at",
			om.role,
			u.id as "user_id",
			u.email
		`).
		Joins("LEFT JOIN users u ON u.id = om.user_id").
		Where("om.organization_id = ?", organizationID).
		Scan(&memberships)

	return memberships
}

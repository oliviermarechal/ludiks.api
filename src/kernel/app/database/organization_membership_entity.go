package database

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationMembershipEntity struct {
	ID             uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID         uuid.UUID          `gorm:"type:uuid;not null"`
	OrganizationID uuid.UUID          `gorm:"type:uuid;not null"`
	Role           string             `gorm:"type:varchar(50);default:'member'"` // member, admin
	CreatedAt      time.Time          `gorm:"autoCreateTime"`
	User           UserEntity         `gorm:"foreignKey:UserID;references:ID"`
	Organization   OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID"`
}

func (OrganizationMembershipEntity) TableName() string {
	return "organization_memberships"
}

package database

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationInvitationEntity struct {
	ID             uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FromUserID     uuid.UUID          `gorm:"type:uuid;not null"`
	OrganizationID uuid.UUID          `gorm:"type:uuid;not null"`
	Role           string             `gorm:"type:varchar(50);default:'member'"` // member, admin
	CreatedAt      time.Time          `gorm:"autoCreateTime"`
	FromUser       UserEntity         `gorm:"foreignKey:FromUserID;references:ID"`
	Organization   OrganizationEntity `gorm:"foreignKey:OrganizationID;references:ID"`
	ToEmail        string             `gorm:"type:varchar;not null"`
}

func (OrganizationInvitationEntity) TableName() string {
	return "organization_invitations"
}

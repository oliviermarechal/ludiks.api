package database

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationEntity struct {
	ID          uuid.UUID                      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string                         `gorm:"type:varchar(255);not null"`
	Plan        string                         `gorm:"type:varchar(50);default:'free'"` // free, pro, scale, custom
	EventsQuota int                            `gorm:"default:5000"`
	EventsUsed  int                            `gorm:"default:0"`
	Pricing     int                            `gorm:"default:0"`
	CreatedAt   time.Time                      `gorm:"autoCreateTime"`
	Projects    []ProjectEntity                `gorm:"foreignKey:OrganizationID"`
	Memberships []OrganizationMembershipEntity `gorm:"foreignKey:OrganizationID"`
	// Billing     OrganizationBillingEntity      `gorm:"foreignKey:OrganizationID"`
}

func (OrganizationEntity) TableName() string {
	return "organizations"
}

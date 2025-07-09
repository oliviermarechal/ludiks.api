package database

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationEntity struct {
	ID               uuid.UUID                        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name             string                           `gorm:"type:varchar(255);not null"`
	Plan             string                           `gorm:"type:varchar(50);default:'free'"` // free, pro, scale, custom
	EventsUsed       int                              `gorm:"default:0"`
	Pricing          int                              `gorm:"default:0"`
	StripeCustomerID string                           `gorm:"type:varchar(255)"`
	CreatedAt        time.Time                        `gorm:"autoCreateTime"`
	Projects         []ProjectEntity                  `gorm:"foreignKey:OrganizationID"`
	Memberships      []OrganizationMembershipEntity   `gorm:"foreignKey:OrganizationID"`
	Subscriptions    []OrganizationSubscriptionEntity `gorm:"foreignKey:OrganizationID"`
}

func (OrganizationEntity) TableName() string {
	return "organizations"
}

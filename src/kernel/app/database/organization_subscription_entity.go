package database

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationSubscriptionEntity struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrganizationID       uuid.UUID `gorm:"type:uuid;not null"`
	StripeSubscriptionID string    `gorm:"type:varchar(255)"`
	Status               string    `gorm:"type:varchar(50);not null"`
	CurrentPeriodEnd     *time.Time
	CancelRequestedAt    *time.Time
	CreatedAt            time.Time          `gorm:"autoCreateTime"`
	UpdatedAt            time.Time          `gorm:"autoUpdateTime"`
	Organization         OrganizationEntity `gorm:"foreignKey:OrganizationID"`
}

func (OrganizationSubscriptionEntity) TableName() string {
	return "organization_subscriptions"
}

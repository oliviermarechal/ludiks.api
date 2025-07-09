package database

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceEntity struct {
	ID              uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	StripeInvoiceID string             `gorm:"type:varchar(255);uniqueIndex;not null"`
	OrganizationID  uuid.UUID          `gorm:"type:uuid;not null"`
	Organization    OrganizationEntity `gorm:"foreignKey:OrganizationID"`
	Amount          int                `gorm:"not null"` // Amount in cents
	Currency        string             `gorm:"type:varchar(3);not null;default:'usd'"`
	Status          string             `gorm:"type:varchar(50);not null"` // draft, open, paid, void, uncollectible
	PeriodMonth     int                `gorm:"not null"`                  // 1-12
	PeriodYear      int                `gorm:"not null"`
	LinkUrl         *string            `gorm:"type:text"`
	CreatedAt       time.Time          `gorm:"autoCreateTime"`
	UpdatedAt       time.Time          `gorm:"autoUpdateTime"`
}

func (InvoiceEntity) TableName() string {
	return "invoices"
}

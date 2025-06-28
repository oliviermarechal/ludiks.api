package database

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	ID            uuid.UUID                      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email         string                         `gorm:"type:varchar(100);unique;not null"`
	Password      string                         `gorm:"type:varchar(100);default:null"`
	GoogleId      string                         `gorm:"type:varchar(100);default:null"`
	CreatedAt     time.Time                      `gorm:"autoCreateTime"`
	Organizations []OrganizationMembershipEntity `gorm:"foreignKey:UserID;references:ID"`
}

func (UserEntity) TableName() string {
	return "users"
}

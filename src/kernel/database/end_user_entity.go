package database

import (
	"time"

	"github.com/google/uuid"
)

type EndUserEntity struct {
	ID          uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ExternalID  string        `gorm:"type:varchar(100);uniqueIndex:idx_project_external"`
	FullName    string        `gorm:"type:varchar(100);default:null"`
	Email       string        `gorm:"type:varchar(100);default:null"`
	Picture     string        `gorm:"type:varchar(100);default:null"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	LastLoginAt time.Time     `gorm:"autoCreateTime"`
	ProjectID   uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex:idx_project_external"`
	Project     ProjectEntity `gorm:"foreignKey:ProjectID"`
}

func (EndUserEntity) TableName() string {
	return "end_users"
}

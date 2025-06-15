package database

import (
	"time"

	"github.com/google/uuid"
)

type EndUserMetadataEntity struct {
	EndUserID uuid.UUID     `gorm:"type:uuid;not null;primaryKey"`
	EndUser   EndUserEntity `gorm:"foreignKey:EndUserID"`
	KeyName   string        `gorm:"type:varchar(100);not null;primaryKey"`
	Value     string        `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
}

func (EndUserMetadataEntity) TableName() string {
	return "end_user_metadatas"
}

package models

import "time"

type EndUserMetadata struct {
	EndUserID string    `json:"end_user_id" gorm:"primaryKey"`
	KeyName   string    `json:"key_name" gorm:"primaryKey"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

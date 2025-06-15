package models

import "time"

type EndUser struct {
	ID              string            `json:"id"`
	ExternalID      string            `json:"external_id"`
	FullName        string            `json:"full_name"`
	Email           *string           `json:"email"`
	Picture         *string           `json:"picture"`
	CreatedAt       time.Time         `json:"created_at"`
	ProjectID       string            `json:"project_id"`
	LastLoginAt     time.Time         `json:"last_login_at"`
	CurrentStreak   int               `json:"current_streak"`
	LongestStreak   int               `json:"longest_streak"`
	EndUserMetadata []EndUserMetadata `json:"end_user_metadata" gorm:"foreignKey:EndUserID;references:ID"`
}

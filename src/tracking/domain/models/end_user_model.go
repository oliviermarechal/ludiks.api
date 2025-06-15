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

func CreateEndUser(ID string, ExternalID string, FullName string, Email *string, Picture *string, ProjectID string) *EndUser {
	return &EndUser{
		ID:            ID,
		ExternalID:    ExternalID,
		FullName:      FullName,
		Email:         Email,
		Picture:       Picture,
		ProjectID:     ProjectID,
		CreatedAt:     time.Now(),
		LastLoginAt:   time.Now(),
		CurrentStreak: 1,
		LongestStreak: 1,
	}
}

func (e *EndUser) HasMetadataWithValue(keyName string, value string) bool {
	for _, metadata := range e.EndUserMetadata {
		if metadata.KeyName == keyName && convertToString(metadata.Value) == value {
			return true
		}
	}
	return false
}

func (e *EndUser) GetMetadata(keyName string) *EndUserMetadata {
	for _, metadata := range e.EndUserMetadata {
		if metadata.KeyName == keyName {
			return &metadata
		}
	}
	return nil
}

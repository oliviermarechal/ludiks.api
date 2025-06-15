package models

import "time"

type EndUserMetadata struct {
	EndUserID string    `json:"end_user_id" gorm:"primaryKey"`
	KeyName   string    `json:"key_name" gorm:"primaryKey"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

func (EndUserMetadata) TableName() string {
	return "end_user_metadatas"
}

func CreateEndUserMetadata(EndUserID string, KeyName string, Value string) *EndUserMetadata {
	return &EndUserMetadata{
		EndUserID: EndUserID,
		KeyName:   KeyName,
		Value:     Value,
		CreatedAt: time.Now(),
	}
}

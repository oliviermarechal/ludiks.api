package models

import "time"

type Project struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ApiKeys   *[]ApiKey `json:"api_keys,omitempty" gorm:"foreignKey:ProjectID"`
}

func CreateProject(id string, name string) *Project {
	return &Project{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

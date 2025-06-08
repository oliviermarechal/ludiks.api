package database

import "time"

type StepEntity struct {
	ID                  string        `gorm:"primary_key"`
	Name                string        `gorm:"not null"`
	Description         string        `gorm:"default:null"`
	CompletionThreshold int           `gorm:"not null;default:0"`
	CircuitID           string        `gorm:"not null"`
	StepNumber          int           `gorm:"not null"`
	EventName           string        `gorm:"not null"`
	Circuit             CircuitEntity `gorm:"foreignKey:CircuitID"`
	CreatedAt           time.Time     `gorm:"autoCreateTime"`
}

func (StepEntity) TableName() string {
	return "steps"
}

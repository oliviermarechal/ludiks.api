package database

import "time"

type RewardEntity struct {
	ID                        string        `gorm:"primary_key"`
	Name                      string        `gorm:"not null"`
	Description               string        `gorm:"default:null"`
	UnlockOnCircuitCompletion bool          `gorm:"not null;default:false"`
	CircuitID                 string        `gorm:"not null"`
	StepID                    string        `gorm:"default:null"`
	Circuit                   CircuitEntity `gorm:"foreignKey:CircuitID"`
	Step                      StepEntity    `gorm:"foreignKey:StepID"`
	CreatedAt                 time.Time     `gorm:"autoCreateTime"`
}

func (RewardEntity) TableName() string {
	return "rewards"
}

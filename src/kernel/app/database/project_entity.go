package database

import (
	"time"

	"github.com/google/uuid"
)

type ProjectEntity struct {
	ID             uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name           string             `gorm:"type:varchar(255);not null"`
	OrganizationID uuid.UUID          `gorm:"type:uuid;not null"`
	CreatedAt      time.Time          `gorm:"autoCreateTime"`
	Organization   OrganizationEntity `gorm:"foreignKey:OrganizationID"`
	ApiKeys        []ApiKeyEntity     `gorm:"many2many:project_api_keys;"`
	Circuits       []CircuitEntity    `gorm:"many2many:project_circuits;"`
}

func (ProjectEntity) TableName() string {
	return "projects"
}

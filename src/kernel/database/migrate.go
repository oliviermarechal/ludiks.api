package database

import (
	"gorm.io/gorm"
)

func RegisterEntities() []interface{} {
	return []interface{}{
		&UserEntity{},
		&ProjectEntity{},
		&ApiKeyEntity{},
		&CircuitEntity{},
		&StepEntity{},
		&EndUserEntity{},
		&EndUserMetadataEntity{},
		&UserCircuitProgressionEntity{},
		&UserStepProgressionEntity{},
		&RewardEntity{},
		&ProjectMetadataKeysEntity{},
	}
}

func RunMigrations(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	entities := RegisterEntities()
	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			return err
		}
	}

	return nil
}

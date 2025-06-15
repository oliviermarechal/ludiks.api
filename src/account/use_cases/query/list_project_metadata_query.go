package query

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

func ListProjectMetadataQuery(db *gorm.DB, projectId string) ([]*models.ProjectMetadataKey, error) {
	var projectMetadataKeys []*models.ProjectMetadataKey

	if err := db.Where("project_id = ?", projectId).Preload("Values").Find(&projectMetadataKeys).Error; err != nil {
		return nil, err
	}

	return projectMetadataKeys, nil
}

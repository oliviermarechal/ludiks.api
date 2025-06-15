package infra_repositories

import (
	"gorm.io/gorm"

	"ludiks/src/tracking/domain/models"
)

type MetadataRepository struct {
	db *gorm.DB
}

func NewMetadataRepository(db *gorm.DB) *MetadataRepository {
	return &MetadataRepository{db: db}
}

func (r *MetadataRepository) ListProjectMetadataKey(projectId string) (*[]models.ProjectMetadataKey, error) {
	var projectMetadataKeys []models.ProjectMetadataKey

	if err := r.db.Where("project_id = ?", projectId).Preload("Values").Find(&projectMetadataKeys).Error; err != nil {
		return nil, err
	}

	return &projectMetadataKeys, nil
}

func (r *MetadataRepository) BatchCreateEndUserMetadata(metadata []*models.EndUserMetadata) error {
	if len(metadata) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&metadata).Error
	})
}

func (r *MetadataRepository) BatchCreateProjectMetadataKeys(keys []*models.ProjectMetadataKey) error {
	if len(keys) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&keys).Error
	})
}

func (r *MetadataRepository) BatchCreateMetadataValues(values []*models.MetadataValue) error {
	if len(values) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&values).Error
	})
}

func (r *MetadataRepository) DeleteEndUserMetadata(endUserID string, keyName string) error {
	return r.db.Where("end_user_id = ? AND key_name = ?", endUserID, keyName).Delete(&models.EndUserMetadata{}).Error
}

func (r *MetadataRepository) UpdateEndUserMetadata(metadata *models.EndUserMetadata) error {
	return r.db.Model(&models.EndUserMetadata{}).
		Where("end_user_id = ? AND key_name = ?", metadata.EndUserID, metadata.KeyName).
		Update("value", metadata.Value).Error
}

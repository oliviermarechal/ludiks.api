package domain_repositories

import "ludiks/src/tracking/domain/models"

type MetadataRepository interface {
	ListProjectMetadataKey(projectId string) (*[]models.ProjectMetadataKey, error)
	BatchCreateEndUserMetadata(metadata []*models.EndUserMetadata) error
	BatchCreateProjectMetadataKeys(keys []*models.ProjectMetadataKey) error
	BatchCreateMetadataValues(values []*models.MetadataValue) error
	DeleteEndUserMetadata(endUserID string, keyName string) error
	UpdateEndUserMetadata(metadata *models.EndUserMetadata) error
}

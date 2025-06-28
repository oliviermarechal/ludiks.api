package infra_repositories

import (
	"ludiks/src/tracking/domain/models"

	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (o *OrganizationRepository) Find(id string) (*models.Organization, error) {
	var organization models.Organization
	if err := o.db.Where("id = ?", id).First(&organization).Error; err != nil {
		return nil, err
	}

	return &organization, nil
}

func (o *OrganizationRepository) FindByProjectID(projectID string) (*models.Organization, error) {
	var organization models.Organization
	if err := o.db.Select("organizations.*").
		Joins("INNER JOIN projects ON projects.organization_id = organizations.id").
		Where("projects.id = ?", projectID).
		First(&organization).Error; err != nil {
		return nil, err
	}

	return &organization, nil
}

func (o *OrganizationRepository) IncrementQuotaUsed(organization *models.Organization) error {
	return o.db.Model(&models.Organization{}).Where("id = ?", organization.ID).Update("events_used", organization.EventsUsed).Error
}

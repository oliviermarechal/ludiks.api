package domain_repositories

import "ludiks/src/tracking/domain/models"

type OrganizationRepository interface {
	Find(id string) (*models.Organization, error)
	FindByProjectID(projectID string) (*models.Organization, error)
	IncrementQuotaUsed(organization *models.Organization) error
}

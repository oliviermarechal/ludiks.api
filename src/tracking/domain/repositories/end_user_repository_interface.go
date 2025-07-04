package domain_repositories

import "ludiks/src/tracking/domain/models"

type EndUserRepository interface {
	Create(endUser *models.EndUser) (*models.EndUser, error)
	Update(endUser *models.EndUser) (*models.EndUser, error)
	Find(ID string) (*models.EndUser, error)
	FindByExternalID(externalID string) (*models.EndUser, error)
}

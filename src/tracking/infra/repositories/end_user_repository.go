package infra_repositories

import (
	"ludiks/src/tracking/domain/models"

	"gorm.io/gorm"
)

type EndUserRepository struct {
	db *gorm.DB
}

func NewEndUserRepository(db *gorm.DB) *EndUserRepository {
	return &EndUserRepository{db: db}
}

func (r *EndUserRepository) Create(endUser *models.EndUser) (*models.EndUser, error) {
	err := r.db.Create(endUser).Error

	return endUser, err
}

func (r *EndUserRepository) Update(endUser *models.EndUser) (*models.EndUser, error) {
	err := r.db.Save(endUser).Error

	return endUser, err
}

func (r *EndUserRepository) Find(ID string) (*models.EndUser, error) {
	var endUser models.EndUser
	err := r.db.Where("id = ?", ID).Preload("EndUserMetadata").First(&endUser).Error

	return &endUser, err
}

func (r *EndUserRepository) FindByExternalID(externalID string) (*models.EndUser, error) {
	var endUser models.EndUser
	err := r.db.Where("external_id = ?", externalID).Preload("EndUserMetadata").First(&endUser).Error

	return &endUser, err
}

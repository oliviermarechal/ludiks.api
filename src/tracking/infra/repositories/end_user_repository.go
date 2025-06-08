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

func (r *EndUserRepository) Find(ID string) (*models.EndUser, error) {
	var endUser models.EndUser
	err := r.db.Where("id = ?", ID).First(&endUser).Error

	return &endUser, err
}

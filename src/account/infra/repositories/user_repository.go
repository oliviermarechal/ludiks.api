package infra_repositories

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Find(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByGid(gID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("google_id = ?", gID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateGid(id string, gID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&user).Update("google_id", gID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) AssociateProject(user_id string, project_id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, err
	}

	project := models.Project{ID: project_id}
	if err := r.db.Model(&user).Association("Projects").Append(&project); err != nil {
		return nil, err
	}

	return &user, nil
}

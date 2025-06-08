package infra_repositories

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Find(id string) (*models.Project, error) {
	var project models.Project
	if err := r.db.Where("id = ?", id).First(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepository) Create(project *models.Project) (*models.Project, error) {
	if err := r.db.Create(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) CreateApiKey(apiKey *models.ApiKey) (*models.ApiKey, error) {
	if err := r.db.Create(apiKey).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&apiKey).Association("Project").Find(&apiKey.Project); err != nil {
		return nil, err
	}

	return apiKey, nil
}

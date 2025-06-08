package domain_repositories

import "ludiks/src/account/domain/models"

type ProjectRepository interface {
	Create(project *models.Project) (*models.Project, error)
	CreateApiKey(apiKey *models.ApiKey) (*models.ApiKey, error)
	Find(id string) (*models.Project, error)
}

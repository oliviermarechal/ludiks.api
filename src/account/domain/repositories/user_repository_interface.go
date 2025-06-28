package domain_repositories

import "ludiks/src/account/domain/models"

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	Find(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByGid(gID string) (*models.User, error)
	UpdateGid(id string, gID string) (*models.User, error)
}

package create_end_user

import (
	"ludiks/src/tracking/domain/models"
	domain_repositories "ludiks/src/tracking/domain/repositories"

	"github.com/google/uuid"
)

type CreateEndUserUseCase struct {
	endUserRepository domain_repositories.EndUserRepository
}

func NewCreateEndUserUseCase(
	endUserRepository domain_repositories.EndUserRepository,
) *CreateEndUserUseCase {
	return &CreateEndUserUseCase{
		endUserRepository: endUserRepository,
	}
}

func (u *CreateEndUserUseCase) Execute(command CreateEndUserCommand) (*models.EndUser, error) {
	endUser := models.CreateEndUser(uuid.New().String(), command.ID, command.FullName, command.Email, command.Picture, command.ProjectID)

	return u.endUserRepository.Create(endUser)
}

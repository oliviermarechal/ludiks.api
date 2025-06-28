package registration

import (
	"ludiks/src/account/domain/models"
	providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"

	"time"

	"github.com/google/uuid"
)

type RegistrationUseCaseResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

type RegistrationUseCase struct {
	userRepository         domain_repositories.UserRepository
	organizationRepository domain_repositories.OrganizationRepository
	encrypter              providers.Encrypter
	jwtProvider            providers.JwtProvider
}

func NewRegistrationUseCase(
	userRepository domain_repositories.UserRepository,
	organizationRepository domain_repositories.OrganizationRepository,
	encrypter providers.Encrypter,
	jwtProvider providers.JwtProvider,
) *RegistrationUseCase {
	return &RegistrationUseCase{
		userRepository:         userRepository,
		organizationRepository: organizationRepository,
		encrypter:              encrypter,
		jwtProvider:            jwtProvider,
	}
}

func (r *RegistrationUseCase) Execute(
	command RegistrationCommand,
) (*RegistrationUseCaseResponse, error) {
	encryptedPassword, err := r.encrypter.Encrypt(command.Password)
	if err != nil {
		return nil, err
	}

	user := models.CreateUserWithPassword(uuid.New().String(), command.Email, encryptedPassword)

	createdUser, err := r.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	if command.InviteID != nil && *command.InviteID != "" {
		invit, _ := r.organizationRepository.FindInvit(*command.InviteID)
		if invit != nil && time.Since(invit.CreatedAt) < 7*24*time.Hour {
			organizationMember := models.CreateOrganizationMember(
				uuid.New().String(),
				createdUser.ID,
				invit.OrganizationID,
				invit.Role,
			)
			r.organizationRepository.CreateOrganizationMember(organizationMember)
		}
	}

	token, err := r.jwtProvider.GenerateToken(createdUser)
	if err != nil {
		return nil, err
	}

	return &RegistrationUseCaseResponse{
		User:  createdUser,
		Token: token,
	}, nil
}

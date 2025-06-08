package registration

import (
	"ludiks/src/account/domain/models"
	providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"

	"github.com/google/uuid"
)

type RegistrationUseCaseResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func RegistrationUseCase(
	userRepository domain_repositories.UserRepository,
	encrypter providers.Encrypter,
	jwtProvider providers.JwtProvider,
	command RegistrationCommand,
) (*RegistrationUseCaseResponse, error) {
	encryptedPassword, err := encrypter.Encrypt(command.Password)
	if err != nil {
		return nil, err
	}

	user := models.CreateUserWithPassword(uuid.New().String(), command.Email, encryptedPassword)

	createdUser, err := userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	token, err := jwtProvider.GenerateToken(createdUser)
	if err != nil {
		return nil, err
	}

	return &RegistrationUseCaseResponse{
		User:  createdUser,
		Token: token,
	}, nil
}

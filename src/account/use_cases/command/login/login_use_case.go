package login

import (
	"errors"

	"ludiks/src/account/domain/models"
	providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
)

type LoginUseCaseResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func LoginUseCase(
	userRepository domain_repositories.UserRepository,
	jwtProvider providers.JwtProvider,
	encrypter providers.Encrypter,
	command LoginCommand,
) (*LoginUseCaseResponse, error) {
	user, err := userRepository.FindByEmail(command.Email)
	if err != nil {
		return nil, err
	}

	if !encrypter.Compare(command.Password, *user.Password) {
		return nil, errors.New("Invalid credentials")
	}

	token, err := jwtProvider.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginUseCaseResponse{
		Token: token,
		User:  user,
	}, nil
}

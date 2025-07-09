package google_auth

import (
	"context"
	"ludiks/config"
	"ludiks/src/account/domain/models"
	providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type GoogleAuthUseCase struct {
	UserRepository         domain_repositories.UserRepository
	organizationRepository domain_repositories.OrganizationRepository
	jwtProvider            providers.JwtProvider
}

type GoogleAuthResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

func NewGoogleAuthUseCase(
	UserRepository domain_repositories.UserRepository,
	organizationRepository domain_repositories.OrganizationRepository,
	jwtProvider providers.JwtProvider,
) *GoogleAuthUseCase {
	return &GoogleAuthUseCase{
		UserRepository:         UserRepository,
		organizationRepository: organizationRepository,
		jwtProvider:            jwtProvider,
	}
}

func (ga *GoogleAuthUseCase) Handle(command GoogleAuthCommand) (GoogleAuthResponse, error) {
	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, command.IdToken, config.AppConfig.GoogleClientID)
	if err != nil {
		return GoogleAuthResponse{}, err
	}

	gid := payload.Claims["sub"].(string)
	email := payload.Claims["email"].(string)

	user, _ := ga.UserRepository.FindByGid(gid)

	if user != nil {
		token, err := ga.jwtProvider.GenerateToken(user)
		if err != nil {
			return GoogleAuthResponse{}, err
		}

		return GoogleAuthResponse{
			User:  *user,
			Token: token,
		}, nil
	}

	user, _ = ga.UserRepository.FindByEmail(email)
	if user != nil {
		ga.UserRepository.UpdateGid(user.ID, gid)
		token, err := ga.jwtProvider.GenerateToken(user)
		if err != nil {
			return GoogleAuthResponse{}, err
		}

		return GoogleAuthResponse{
			User:  *user,
			Token: token,
		}, nil
	}

	user = models.CreateUserWithGoogleId(
		uuid.New().String(),
		email,
		gid,
	)

	createdUser, err := ga.UserRepository.Create(user)
	if err != nil {
		return GoogleAuthResponse{}, nil
	}

	if command.InviteID != nil {
		invit, err := ga.organizationRepository.FindInvit(*command.InviteID)
		if err == nil && time.Since(invit.CreatedAt) < 7*24*time.Hour {
			organizationMember := models.CreateOrganizationMember(
				uuid.New().String(),
				createdUser.ID,
				invit.OrganizationID,
				invit.Role,
			)
			ga.organizationRepository.CreateOrganizationMember(organizationMember)
		}
	}

	token, err := ga.jwtProvider.GenerateToken(createdUser)
	if err != nil {
		return GoogleAuthResponse{}, err
	}

	return GoogleAuthResponse{
		User:  *createdUser,
		Token: token,
	}, nil
}

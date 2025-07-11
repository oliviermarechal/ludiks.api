package auth

import (
	providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/google_auth"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GoogleAuthHandler struct {
	userRepository         domain_repositories.UserRepository
	organizationRepository domain_repositories.OrganizationRepository
	jwtProvider            providers.JwtProvider
}

func NewGoogleAuthHandler(
	userRepo domain_repositories.UserRepository,
	organizationRepository domain_repositories.OrganizationRepository,
	jwtProvider providers.JwtProvider,
) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		userRepository:         userRepo,
		organizationRepository: organizationRepository,
		jwtProvider:            jwtProvider,
	}
}

func (h *GoogleAuthHandler) Handle(c *gin.Context) {
	var dto google_auth.GoogleAuthDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := dto.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	googleAuthResult, err := google_auth.NewGoogleAuthUseCase(
		h.userRepository,
		h.organizationRepository,
		h.jwtProvider,
	).Handle(google_auth.GoogleAuthCommand{
		IdToken:  dto.IdToken,
		InviteID: dto.InviteID,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, googleAuthResult)
}

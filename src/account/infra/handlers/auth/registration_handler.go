package auth

import (
	providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/registration"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	userRepository domain_repositories.UserRepository
	encrypter      providers.Encrypter
	jwtProvider    providers.JwtProvider
}

func NewRegistrationHandler(
	userRepo domain_repositories.UserRepository,
	encrypter providers.Encrypter,
	jwtProvider providers.JwtProvider,
) *RegistrationHandler {
	return &RegistrationHandler{
		userRepository: userRepo,
		encrypter:      encrypter,
		jwtProvider:    jwtProvider,
	}
}

func (h *RegistrationHandler) Handle(c *gin.Context) {
	var dto registration.RegistrationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := dto.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	registrationResult, err := registration.RegistrationUseCase(
		h.userRepository,
		h.encrypter,
		h.jwtProvider,
		registration.RegistrationCommand{
			Email:    dto.Email,
			Password: dto.Password,
		},
	)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, registrationResult)
}

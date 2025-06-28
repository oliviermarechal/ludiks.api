package auth

import (
	providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/registration"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	userRepository         domain_repositories.UserRepository
	organizationRepository domain_repositories.OrganizationRepository
	encrypter              providers.Encrypter
	jwtProvider            providers.JwtProvider
}

func NewRegistrationHandler(
	userRepository domain_repositories.UserRepository,
	organizationRepository domain_repositories.OrganizationRepository,
	encrypter providers.Encrypter,
	jwtProvider providers.JwtProvider,
) *RegistrationHandler {
	return &RegistrationHandler{
		userRepository:         userRepository,
		organizationRepository: organizationRepository,
		encrypter:              encrypter,
		jwtProvider:            jwtProvider,
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

	registrationResult, err := registration.NewRegistrationUseCase(
		h.userRepository,
		h.organizationRepository,
		h.encrypter,
		h.jwtProvider,
	).Execute(
		registration.RegistrationCommand{
			Email:    dto.Email,
			Password: dto.Password,
			InviteID: dto.InviteID,
		},
	)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, registrationResult)
}

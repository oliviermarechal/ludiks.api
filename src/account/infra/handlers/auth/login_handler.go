package auth

import (
	providers "ludiks/src/account/domain/providers"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/login"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	userRepository domain_repositories.UserRepository
	encrypter      providers.Encrypter
	jwtProvider    providers.JwtProvider
}

func NewLoginHandler(
	userRepo domain_repositories.UserRepository,
	encrypter providers.Encrypter,
	jwtProvider providers.JwtProvider,
) *LoginHandler {
	return &LoginHandler{
		userRepository: userRepo,
		encrypter:      encrypter,
		jwtProvider:    jwtProvider,
	}
}

func (h *LoginHandler) Handle(c *gin.Context) {
	var dto login.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := dto.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	loginResult, err := login.LoginUseCase(
		h.userRepository,
		h.jwtProvider,
		h.encrypter,
		login.LoginCommand{
			Email:    dto.Email,
			Password: dto.Password,
		},
	)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusOK, loginResult)
}

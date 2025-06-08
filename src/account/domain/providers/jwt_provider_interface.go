package providers

import (
	"ludiks/src/account/domain/models"

	"github.com/golang-jwt/jwt/v5"
)

type JwtProvider interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(token string) (jwt.MapClaims, error)
}

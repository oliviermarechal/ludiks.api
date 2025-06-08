package providers

import (
	"errors"
	"os"
	"time"

	"ludiks/src/account/domain/models"

	"github.com/golang-jwt/jwt/v5"
)

type JwtProvider struct {
	secretKey []byte
}

func NewJwtProvider() *JwtProvider {
	return &JwtProvider{
		secretKey: []byte(os.Getenv("JWT_SECRET_KEY")),
	}
}

func (p *JwtProvider) GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(p.secretKey)
}

func (p *JwtProvider) ValidateToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return p.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

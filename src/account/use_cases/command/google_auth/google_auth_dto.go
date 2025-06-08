package google_auth

import "github.com/go-playground/validator/v10"

type GoogleAuthDTO struct {
	IdToken string `json:"idToken" validate:"required"`
}

func (g *GoogleAuthDTO) Validate() error {
	validator := validator.New()
	return validator.Struct(g)
}

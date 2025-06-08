package login

import "github.com/go-playground/validator/v10"

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=8"`
}

func (dto *LoginDTO) Validate() error {
	return validator.New().Struct(dto)
}

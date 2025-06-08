package create_end_user

import "github.com/go-playground/validator/v10"

type CreateEndUserDTO struct {
	ID       string  `json:"id" validate:"required"`
	FullName string  `json:"full_name" validate:"required"`
	Email    *string `json:"email"`
	Picture  *string `json:"picture"`
}

func (dto *CreateEndUserDTO) Validate() error {
	return validator.New().Struct(dto)
}

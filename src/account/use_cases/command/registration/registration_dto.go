package registration

import "github.com/go-playground/validator/v10"

type RegistrationDTO struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"min=8"`
	InviteID *string `json:"inviteId"`
}

func (dto *RegistrationDTO) Validate() error {
	return validator.New().Struct(dto)
}

package log_end_user

import "github.com/go-playground/validator/v10"

type LogEndUserDTO struct {
	ID       string                 `json:"id" validate:"required"`
	FullName string                 `json:"full_name" validate:"required"`
	Email    *string                `json:"email"`
	Picture  *string                `json:"picture"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (dto *LogEndUserDTO) Validate() error {
	return validator.New().Struct(dto)
}

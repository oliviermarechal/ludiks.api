package create_api_key

import "github.com/go-playground/validator/v10"

type CreateApiKeyDTO struct {
	Name string `json:"name" validate:"required"`
}

func (dto *CreateApiKeyDTO) Validate() error {
	return validator.New().Struct(dto)
}

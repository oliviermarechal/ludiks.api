package create_project

import "github.com/go-playground/validator/v10"

type CreateProjectDTO struct {
	Name           string `json:"name" validate:"required"`
	OrganizationID string `json:"organizationId" validate:"required"`
}

func (dto *CreateProjectDTO) Validate() error {
	return validator.New().Struct(dto)
}

package create_organization_and_project

import "github.com/go-playground/validator/v10"

type CreateOrganizationAndProjectDTO struct {
	Name string `json:"name" validate:"required"`
}

func (dto *CreateOrganizationAndProjectDTO) Validate() error {
	return validator.New().Struct(dto)
}

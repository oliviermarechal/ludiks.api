package update_project

import "github.com/go-playground/validator/v10"

type UpdateProjectDTO struct {
	Name string `json:"name" validate:"required"`
}

func (dto *UpdateProjectDTO) Validate() error {
	return validator.New().Struct(dto)
}

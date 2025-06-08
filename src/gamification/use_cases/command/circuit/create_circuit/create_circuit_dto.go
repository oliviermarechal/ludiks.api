package create_circuit

import "github.com/go-playground/validator/v10"

type CreateCircuitDTO struct {
	ProjectID string `json:"projectId" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Type      string `json:"type" validate:"required"`
}

func (c *CreateCircuitDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

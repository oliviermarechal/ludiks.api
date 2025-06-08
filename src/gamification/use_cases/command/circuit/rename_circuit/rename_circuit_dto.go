package rename_circuit

import "github.com/go-playground/validator/v10"

type RenameCircuitDTO struct {
	Name string `json:"name" validate:"required"`
}

func (c *RenameCircuitDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

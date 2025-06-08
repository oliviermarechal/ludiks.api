package set_circuit_steps

import "github.com/go-playground/validator/v10"

type SetCircuitStepsDTO struct {
	Steps []StepDto `json:"steps" validate:"required,dive"`
}

type StepDto struct {
	Name                string  `json:"name" validate:"required"`
	Description         *string `json:"description"`
	CompletionThreshold int     `json:"completionThreshold" validate:"required"`
	EventName           string  `json:"eventName" validate:"required"`
}

func (dto *SetCircuitStepsDTO) Validate() error {
	v := validator.New()
	return v.Struct(dto)
}

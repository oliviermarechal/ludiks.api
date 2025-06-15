package reordering_steps

import "github.com/go-playground/validator/v10"

type StepOrderingDTO struct {
	StepID     string `json:"stepId" validate:"required"`
	StepNumber int    `json:"stepNumber" validate:"required"`
}

type ReorderingStepsDTO struct {
	Steps []StepOrderingDTO `json:"steps" validate:"required"`
}

func (dto *ReorderingStepsDTO) Validate() error {
	v := validator.New()
	return v.Struct(dto)
}

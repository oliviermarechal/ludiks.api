package create_step

import "github.com/go-playground/validator/v10"

type CreateStepDTO struct {
	Name                string  `json:"name" validate:"required"`
	Description         *string `json:"description"`
	CompletionThreshold int     `json:"completionThreshold" validate:"required"`
	EventName           string  `json:"eventName" validate:"required"`
	StepNumber          *int    `json:"stepNumber"`
}

func (dto *CreateStepDTO) Validate() error {
	v := validator.New()
	return v.Struct(dto)
}

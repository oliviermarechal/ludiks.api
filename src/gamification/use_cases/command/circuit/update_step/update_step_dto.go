package update_step

import "github.com/go-playground/validator/v10"

type UpdateStepDTO struct {
	Name                string `json:"name" validate:"required"`
	Description         string `json:"description"`
	CompletionThreshold int    `json:"completionThreshold" validate:"required"`
	EventName           string `json:"eventName" validate:"required"`
}

func (dto *UpdateStepDTO) Validate() error {
	v := validator.New()
	return v.Struct(dto)
}

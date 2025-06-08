package update_reward

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type UpdateRewardDTO struct {
	Name                      string  `json:"name" validate:"required"`
	Description               *string `json:"description"`
	StepID                    *string `json:"stepId"`
	UnlockOnCircuitCompletion bool    `json:"unlockOnCircuitCompletion" validate:"reward_condition"`
}

func (dto *UpdateRewardDTO) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("reward_condition", func(fl validator.FieldLevel) bool {
		dto := fl.Parent().Interface().(UpdateRewardDTO)
		return dto.StepID != nil || dto.UnlockOnCircuitCompletion
	})

	if err := validate.Struct(dto); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

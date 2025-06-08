package generate_step

import (
	"ludiks/src/kernel/database"

	"github.com/go-playground/validator/v10"
)

type CurveType string

const (
	Linear      CurveType = "linear"
	Exponential CurveType = "exponential"
	Logarithmic CurveType = "logarithmic"
)

type GenerateStepDTO struct {
	CircuitType   database.CircuitType `json:"circuitType" validate:"required"`
	NumberOfSteps int                  `json:"numberOfSteps" validate:"required,lte=10"`
	Curve         CurveType            `json:"curve" validate:"required"`
	MaxValue      int                  `json:"maxValue" validate:"required"`
	Exponent      *float64             `json:"exponent"`
	EventName     string               `json:"eventName" validate:"required"`
}

func (dto *GenerateStepDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(dto)
}

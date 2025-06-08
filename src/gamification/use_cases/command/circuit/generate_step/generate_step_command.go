package generate_step

import (
	"ludiks/src/kernel/database"
)

type GenerateStepCommand struct {
	CircuitID     string
	CircuitType   database.CircuitType
	NumberOfSteps int
	Curve         CurveType
	MaxValue      int
	Exponent      *float64
	EventName     string
}

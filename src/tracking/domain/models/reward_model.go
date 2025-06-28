package models

import (
	"time"
)

type Reward struct {
	ID                        string    `json:"id"`
	Name                      string    `json:"name"`
	Description               *string   `json:"description"`
	UnlockOnCircuitCompletion bool      `json:"unlockOnCircuitCompletion"`
	CircuitID                 string    `json:"circuitId"`
	StepID                    *string   `json:"stepId"`
	Circuit                   Circuit   `json:"circuit,omitempty"`
	Step                      *Step     `json:"step,omitempty"`
	CreatedAt                 time.Time `json:"createdAt"`
}

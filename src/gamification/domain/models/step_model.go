package models

import (
	"time"
)

type Step struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Description         *string   `json:"description"`
	StepNumber          *int      `json:"stepNumber"`
	CompletionThreshold int       `json:"completionThreshold"`
	EventName           string    `json:"eventName"`
	CircuitID           string    `json:"circuitId"`
	Circuit             *Circuit  `json:"circuit,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
}

func CreateStep(id string, name string, description *string, completionThreshold int, circuitID string, stepNumber *int, eventName string) *Step {
	return &Step{
		ID:                  id,
		Name:                name,
		Description:         description,
		CompletionThreshold: completionThreshold,
		CircuitID:           circuitID,
		StepNumber:          stepNumber,
		EventName:           eventName,
		CreatedAt:           time.Now(),
	}
}

package models

import (
	"time"
)

type Step struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	StepNumber          int       `json:"stepNumber"`
	CompletionThreshold int       `json:"completionThreshold"`
	EventName           string    `json:"eventName"`
	CircuitID           string    `json:"circuitId"`
	Circuit             *Circuit  `json:"circuit,omitempty"`
	CreatedAt           time.Time `json:"createdAt"`
}

package models

import (
	"time"

	"ludiks/src/kernel/app/database"
)

type Circuit struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Type      database.CircuitType `json:"type"`
	Active    bool                 `json:"active"`
	ProjectID string               `json:"project_id"`
	Steps     *[]Step              `json:"steps"`
	CreatedAt time.Time            `json:"created_at"`
}

func CreateCircuit(id string, name string, circuitType database.CircuitType, projectID string) *Circuit {
	return &Circuit{
		ID:        id,
		Name:      name,
		Type:      circuitType,
		ProjectID: projectID,
		Active:    false,
	}
}

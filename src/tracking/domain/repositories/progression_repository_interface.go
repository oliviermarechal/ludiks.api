package domain_repositories

import "ludiks/src/tracking/domain/models"

type ProgressionRepository interface {
	FindUserProgressionByEventName(
		projectId string,
		endUserId string,
		eventName string,
	) (*models.UserCircuitProgression, error)
	CreateProgression(progression *models.UserCircuitProgression) (*models.UserCircuitProgression, error)
	UpdateProgression(progression *models.UserCircuitProgression) (*models.UserCircuitProgression, error)
}

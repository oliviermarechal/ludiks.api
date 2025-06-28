package tracking_handler

import (
	"errors"
	"ludiks/src/kernel/app/handlers"
	domain_repositories "ludiks/src/tracking/domain/repositories"
	"ludiks/src/tracking/use_cases/command/tracking_event"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProgressionTrackingHandler struct {
	progressionRepository  domain_repositories.ProgressionRepository
	endUserRepository      domain_repositories.EndUserRepository
	circuitRepository      domain_repositories.CircuitRepository
	organizationRepository domain_repositories.OrganizationRepository
}

func NewProgressionTrackingHandler(
	progressionRepository domain_repositories.ProgressionRepository,
	endUserRepository domain_repositories.EndUserRepository,
	circuitRepository domain_repositories.CircuitRepository,
	organizationRepository domain_repositories.OrganizationRepository,
) *ProgressionTrackingHandler {
	return &ProgressionTrackingHandler{
		progressionRepository:  progressionRepository,
		endUserRepository:      endUserRepository,
		circuitRepository:      circuitRepository,
		organizationRepository: organizationRepository,
	}
}

func (h *ProgressionTrackingHandler) Handle(c *gin.Context) {
	projectID, ok := c.Get("project_id")
	if !ok {
		handlers.HandleBadRequest(c, errors.New("project_id not found"))
		return
	}

	var trackingEvent tracking_event.TrackingEventDTO
	if err := c.ShouldBindJSON(&trackingEvent); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := trackingEvent.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	useCase := tracking_event.NewTrackingEventUseCase(
		h.progressionRepository,
		h.endUserRepository,
		h.circuitRepository,
		h.organizationRepository,
	)

	response := useCase.Execute(&tracking_event.TrackingEventCommand{
		ProjectID: projectID.(uuid.UUID).String(),
		UserID:    trackingEvent.UserID,
		EventName: trackingEvent.EventName,
		Value:     trackingEvent.Value,
		Timestamp: trackingEvent.Timestamp,
	})

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, response)
	}
}

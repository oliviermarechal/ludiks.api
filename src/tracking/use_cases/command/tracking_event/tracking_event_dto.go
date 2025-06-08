package tracking_event

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type TrackingEventDTO struct {
	UserID    string     `json:"user_id" validate:"required"`
	EventName string     `json:"event_name" validate:"required"`
	Value     *int       `json:"value"`
	Timestamp *time.Time `json:"timestamp"`
}

func (t *TrackingEventDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

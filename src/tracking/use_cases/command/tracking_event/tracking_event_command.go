package tracking_event

import "time"

type TrackingEventCommand struct {
	ProjectID string
	UserID    string
	EventName string
	Value     *int
	Timestamp *time.Time
}

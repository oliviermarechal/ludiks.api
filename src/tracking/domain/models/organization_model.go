package models

import "time"

type Organization struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Plan        string    `json:"plan"`
	EventsQuota int       `json:"eventsQuota"`
	EventsUsed  int       `json:"eventsUsed"`
	Pricing     int       `json:"pricing"`
	CreatedAt   time.Time `json:"created_at"`
}

func (o *Organization) IncrementQuotaUsed() {
	o.EventsUsed++
}

func (o *Organization) HasQuotasReached() bool {
	if o.EventsQuota == -1 {
		return false
	}

	return o.EventsUsed >= o.EventsQuota
}

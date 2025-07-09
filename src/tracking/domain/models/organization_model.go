package models

import (
	"ludiks/config"
	"time"
)

type Organization struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	Name             string    `json:"name"`
	Plan             string    `json:"plan"`
	EventsUsed       int       `json:"eventsUsed"`
	Pricing          int       `json:"pricing"`
	CreatedAt        time.Time `json:"created_at"`
	StripeCustomerID string    `json:"-"`
}

func (o *Organization) IncrementQuotaUsed() {
	o.EventsUsed++
}

func (o *Organization) HasQuotasReached() bool {
	if o.Plan == "free" && o.EventsUsed >= config.AppConfig.FreeEventsLimit {
		return true
	}

	return false
}

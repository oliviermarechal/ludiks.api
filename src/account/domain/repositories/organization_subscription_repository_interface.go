package domain_repositories

import (
	"ludiks/src/account/domain/models"
	"time"
)

type OrganizationSubscriptionRepository interface {
	Create(subscription *models.OrganizationSubscription) (*models.OrganizationSubscription, error)
	FindByID(id string) (*models.OrganizationSubscription, error)
	CancelSubscription(subscriptionID string, cancelledAt *time.Time, endDate *time.Time) error
	ExpireSubscription(subscriptionID string) error
	FindAllActiveWithPeriodEnd() ([]models.OrganizationSubscription, error)
}

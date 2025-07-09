package infra_repositories

import (
	"ludiks/src/account/domain/models"
	"time"

	"gorm.io/gorm"
)

type OrganizationSubscriptionRepository struct {
	db *gorm.DB
}

func NewOrganizationSubscriptionRepository(db *gorm.DB) *OrganizationSubscriptionRepository {
	return &OrganizationSubscriptionRepository{
		db: db,
	}
}

func (r *OrganizationSubscriptionRepository) Create(subscription *models.OrganizationSubscription) (*models.OrganizationSubscription, error) {
	if err := r.db.Create(subscription).Error; err != nil {
		return nil, err
	}
	return subscription, nil
}

func (r *OrganizationSubscriptionRepository) FindByID(id string) (*models.OrganizationSubscription, error) {
	var subscription models.OrganizationSubscription
	if err := r.db.Where("id = ?", id).First(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *OrganizationSubscriptionRepository) CancelSubscription(subscriptionID string, cancelledAt *time.Time, endDate *time.Time) error {
	now := time.Now()
	return r.db.Model(&models.OrganizationSubscription{}).
		Where("id = ?", subscriptionID).
		Updates(map[string]interface{}{
			"cancel_requested_at": now,
			"current_period_end":  endDate,
			"updated_at":          now,
		}).Error
}

func (r *OrganizationSubscriptionRepository) ExpireSubscription(subscriptionID string) error {
	now := time.Now()
	return r.db.Model(&models.OrganizationSubscription{}).
		Where("id = ?", subscriptionID).
		Updates(map[string]interface{}{
			"status":     "cancelled",
			"updated_at": now,
		}).Error
}

func (r *OrganizationSubscriptionRepository) FindAllActiveWithPeriodEnd() ([]models.OrganizationSubscription, error) {
	var subscriptions []models.OrganizationSubscription
	if err := r.db.Where("status = ? AND current_period_end IS NOT NULL", "active").Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

package infra_repositories

import (
	"ludiks/src/account/domain/models"

	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Find(id string) (*models.Organization, error) {
	var organization models.Organization
	if err := r.db.Where("id = ?", id).First(&organization).Error; err != nil {
		return nil, err
	}

	return &organization, nil
}

func (r *OrganizationRepository) Create(organization *models.Organization) (*models.Organization, error) {
	if err := r.db.Create(organization).Error; err != nil {
		return nil, err
	}

	return organization, nil
}

func (r *OrganizationRepository) UserIsMember(userID string, organizationID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.OrganizationMembership{}).
		Where("user_id = ? AND organization_id = ?", userID, organizationID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *OrganizationRepository) CreateOrganizationMember(organizationMember *models.OrganizationMembership) (*models.OrganizationMembership, error) {
	if err := r.db.Create(organizationMember).Error; err != nil {
		return nil, err
	}

	return organizationMember, nil
}

func (r *OrganizationRepository) CreateInvitation(invitation *models.OrganizationInvitation) (*models.OrganizationInvitation, error) {
	if err := r.db.Create(invitation).Error; err != nil {
		return nil, err
	}

	return invitation, nil
}

func (r *OrganizationRepository) FindInvit(invitationID string) (*models.OrganizationInvitation, error) {
	var invitation models.OrganizationInvitation
	if err := r.db.Where("id = ?", invitationID).First(&invitation).Error; err != nil {
		return nil, err
	}

	return &invitation, nil
}

func (r *OrganizationRepository) RemoveInvit(invitationID string) error {
	return r.db.Delete(&models.OrganizationInvitation{}, "id = ?", invitationID).Error
}

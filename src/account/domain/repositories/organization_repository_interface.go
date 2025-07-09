package domain_repositories

import "ludiks/src/account/domain/models"

type OrganizationRepository interface {
	Create(organization *models.Organization) (*models.Organization, error)
	Find(id string) (*models.Organization, error)
	FindByStripeCustomerID(stripeCustomerID string) (*models.Organization, error)
	CreateOrganizationMember(organizationMember *models.OrganizationMembership) (*models.OrganizationMembership, error)
	CreateInvitation(invitation *models.OrganizationInvitation) (*models.OrganizationInvitation, error)
	FindInvit(id string) (*models.OrganizationInvitation, error)
	RemoveInvit(id string) error
	UserIsMember(userID string, organizationID string) (bool, error)
	UpdatePlan(organizationID string, plan string) error
	SetStripeCustomerID(organizationID string, stripeCustomerId string) error
}

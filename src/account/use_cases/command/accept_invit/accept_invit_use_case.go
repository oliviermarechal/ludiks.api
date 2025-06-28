package accept_invit

import (
	"errors"
	"ludiks/src/account/domain/models"
	domain_repositories "ludiks/src/account/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type AcceptInvitUseCase struct {
	organizationRepository domain_repositories.OrganizationRepository
	userRepository         domain_repositories.UserRepository
}

func NewAcceptInvitUseCase(
	organizationRepository domain_repositories.OrganizationRepository,
	userRepository domain_repositories.UserRepository,
) *AcceptInvitUseCase {
	return &AcceptInvitUseCase{
		organizationRepository: organizationRepository,
		userRepository:         userRepository,
	}
}

func (a *AcceptInvitUseCase) Execute(command AcceptInvitCommand) error {
	invit, err := a.organizationRepository.FindInvit(command.InvitID)
	user, _ := a.userRepository.Find(command.UserID)
	if err != nil || user.Email != invit.ToEmail {
		return errors.New("invitation not found or expired")
	}

	if time.Since(invit.CreatedAt) > 7*24*time.Hour {
		a.organizationRepository.RemoveInvit(command.InvitID)
		return errors.New("invitation not found or expired")
	}

	organizationMembership := models.CreateOrganizationMember(
		uuid.New().String(),
		user.ID,
		invit.OrganizationID,
		invit.Role,
	)
	a.organizationRepository.CreateOrganizationMember(organizationMembership)
	a.organizationRepository.RemoveInvit(command.InvitID)

	return nil
}

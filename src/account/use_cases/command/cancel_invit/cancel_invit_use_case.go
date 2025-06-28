package cancel_invit

import (
	"errors"
	domain_repositories "ludiks/src/account/domain/repositories"
)

type CancelInvitUseCase struct {
	organizationRepository domain_repositories.OrganizationRepository
}

func NewCancelInvitUseCase(
	organizationRepository domain_repositories.OrganizationRepository,
) *CancelInvitUseCase {
	return &CancelInvitUseCase{
		organizationRepository: organizationRepository,
	}
}

func (a *CancelInvitUseCase) Execute(command CancelInvitCommand) error {
	invit, err := a.organizationRepository.FindInvit(command.InvitID)
	if err != nil {
		return errors.New("invitation not found or expired")
	}

	userIsOrganizationMember, err := a.organizationRepository.UserIsMember(command.UserID, invit.OrganizationID)
	if !userIsOrganizationMember || err != nil {
		return errors.New("invitation not found or expired")
	}

	a.organizationRepository.RemoveInvit(command.InvitID)

	return nil
}

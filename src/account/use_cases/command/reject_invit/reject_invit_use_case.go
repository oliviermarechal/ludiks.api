package reject_invit

import (
	"errors"
	domain_repositories "ludiks/src/account/domain/repositories"
)

type RejectInvitUseCase struct {
	organizationRepository domain_repositories.OrganizationRepository
	userRepository         domain_repositories.UserRepository
}

func NewRejectInvitUseCase(
	organizationRepository domain_repositories.OrganizationRepository,
	userRepository domain_repositories.UserRepository,
) *RejectInvitUseCase {
	return &RejectInvitUseCase{
		organizationRepository: organizationRepository,
		userRepository:         userRepository,
	}
}

func (a *RejectInvitUseCase) Execute(command RejectInvitCommand) error {
	invit, err := a.organizationRepository.FindInvit(command.InvitID)
	user, _ := a.userRepository.Find(command.UserID)
	if err != nil || user.Email != invit.ToEmail {
		return errors.New("invitation not found or expired")
	}

	a.organizationRepository.RemoveInvit(command.InvitID)

	return nil
}

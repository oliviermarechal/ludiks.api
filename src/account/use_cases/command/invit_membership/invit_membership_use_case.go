package invit_membership

import (
	"errors"
	"fmt"
	"ludiks/src/account/domain/models"
	domain_repositories "ludiks/src/account/domain/repositories"
	domain_providers "ludiks/src/kernel/domain/providers"
	"ludiks/src/kernel/infra/utils"
	"os"

	"github.com/google/uuid"
)

type InvitMembershipUseCase struct {
	organizationRepository domain_repositories.OrganizationRepository
	userRepository         domain_repositories.UserRepository
	mailer                 domain_providers.MailerProvider
}

func NewInvitMembershipUseCase(
	organizationRepository domain_repositories.OrganizationRepository,
	userRepository domain_repositories.UserRepository,
	mailer domain_providers.MailerProvider,
) *InvitMembershipUseCase {
	return &InvitMembershipUseCase{
		organizationRepository: organizationRepository,
		userRepository:         userRepository,
		mailer:                 mailer,
	}
}

func (i *InvitMembershipUseCase) Execute(command InvitMembershipCommand) (*models.OrganizationInvitation, error) {
	organization, _ := i.organizationRepository.Find(command.OrganizationID)
	user, _ := i.userRepository.Find(command.FromUserID)

	existingUser, _ := i.userRepository.FindByEmail(command.Email)

	isMember, _ := i.organizationRepository.UserIsMember(existingUser.ID, command.OrganizationID)
	if isMember {
		return nil, errors.New("user already exist in your organization")
	}

	invitation := models.CreateOrganizationInvitation(
		uuid.New().String(),
		command.FromUserID,
		command.Email,
		command.OrganizationID,
		command.Role,
	)

	createdInvit, err := i.organizationRepository.CreateInvitation(invitation)

	if err != nil {
		return nil, err
	}

	var invitationURL string
	if existingUser != nil {
		invitationURL = fmt.Sprintf("%s/invitations", os.Getenv("FRONT_URL"))
	} else {
		invitationURL = fmt.Sprintf("%s/auth/registration?token=%s", os.Getenv("FRONT_URL"), invitation.ID)
	}

	if false {
		data := map[string]interface{}{
			"from_email":        user.Email,
			"organization_name": organization.Name,
			"invitation_url":    invitationURL,
			"subject":           fmt.Sprintf("Invitation to join %s on Ludiks", organization.Name),
		}

		locale := "en"
		templateName := utils.GetInvitationTemplateName(locale)
		htmlContent, err := utils.ReadEmailTemplateWithBase(templateName, locale, data)
		if err != nil {
			return nil, fmt.Errorf("failed to read invitation template: %w", err)
		}

		err = i.mailer.SendTemplateMailWithContent(
			invitation.ToEmail,
			fmt.Sprintf("Invitation to join %s on Ludiks", organization.Name),
			htmlContent,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to send invitation email: %w", err)
		}
	}

	return createdInvit, nil
}

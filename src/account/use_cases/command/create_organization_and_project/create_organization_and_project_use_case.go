package create_organization_and_project

import (
	"ludiks/src/account/domain/models"
	domain_repositories "ludiks/src/account/domain/repositories"

	"github.com/google/uuid"
)

func CreateOrganizationAndProjectUseCase(
	projectRepository domain_repositories.ProjectRepository,
	organizationRepository domain_repositories.OrganizationRepository,
	userRepository domain_repositories.UserRepository,
	command CreateOrganizationAndProjectCommand,
) (*models.Organization, error) {
	user, err := userRepository.Find(command.UserID)
	if err != nil {
		return nil, err
	}

	organization := models.CreateOrganization(uuid.New().String(), command.Name)
	project := models.CreateProject(uuid.New().String(), command.Name, organization.ID)
	createdOrganization, err := organizationRepository.Create(organization)
	if err != nil {
		return nil, err
	}

	createdProject, err := projectRepository.Create(project)
	if err != nil {
		return nil, err
	}

	_, err = organizationRepository.CreateOrganizationMember(models.CreateOrganizationMember(uuid.New().String(), user.ID, createdOrganization.ID, "admin"))
	if err != nil {
		return nil, err
	}

	var organizationProject = &[]models.Project{
		*createdProject,
	}

	createdOrganization.Projects = organizationProject

	return createdOrganization, nil
}

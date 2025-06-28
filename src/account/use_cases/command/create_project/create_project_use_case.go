package create_project

import (
	"ludiks/src/account/domain/models"
	domain_repositories "ludiks/src/account/domain/repositories"

	"github.com/google/uuid"
)

func CreateProjectUseCase(
	projectRepository domain_repositories.ProjectRepository,
	command CreateProjectCommand,
) (*models.Project, error) {
	project := models.CreateProject(uuid.New().String(), command.Name, command.OrganizationID)
	createdProject, err := projectRepository.Create(project)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return createdProject, nil
}

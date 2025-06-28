package update_project

import (
	"ludiks/src/account/domain/models"
	domain_repositories "ludiks/src/account/domain/repositories"
)

func UpdateProjectUseCase(
	projectRepository domain_repositories.ProjectRepository,
	command UpdateProjectCommand,
) (*models.Project, error) {
	project, err := projectRepository.Find(command.ProjectID)
	if err != nil {
		return nil, err
	}

	project.Name = command.Name

	updated, err := projectRepository.Update(project)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

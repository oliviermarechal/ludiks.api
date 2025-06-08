package create_project

import (
	"ludiks/src/account/domain/models"
	domain_repositories "ludiks/src/account/domain/repositories"

	"github.com/google/uuid"
)

func CreateProjectUseCase(
	projectRepository domain_repositories.ProjectRepository,
	userRepository domain_repositories.UserRepository,
	command CreateProjectCommand,
) (*models.Project, error) {
	user, err := userRepository.Find(command.UserID)
	if err != nil {
		return nil, err
	}

	project := models.CreateProject(uuid.New().String(), command.Name)
	createdProject, err := projectRepository.Create(project)
	if err != nil {
		return nil, err
	}

	_, err = userRepository.AssociateProject(user.ID, createdProject.ID)
	if err != nil {
		return nil, err
	}

	return createdProject, nil
}

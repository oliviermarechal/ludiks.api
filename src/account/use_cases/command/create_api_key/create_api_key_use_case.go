package create_api_key

import (
	"ludiks/src/account/domain/models"
	domain_repositories "ludiks/src/account/domain/repositories"

	"github.com/google/uuid"
)

func CreateApiKeyUseCase(
	projectRepository domain_repositories.ProjectRepository,
	command CreateApiKeyCommand,
) (*models.ApiKey, error) {
	project, err := projectRepository.Find(command.ProjectID)
	if err != nil {
		return nil, err
	}

	apiKey := models.CreateApiKey(uuid.New().String(), project.ID, command.Name, uuid.New().String())
	createdApiKey, err := projectRepository.CreateApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	return createdApiKey, nil
}

package delete_api_key

import (
	domain_repositories "ludiks/src/account/domain/repositories"
)

func DeleteApiKeyUseCase(
	projectRepository domain_repositories.ProjectRepository,
	command DeleteApiKeyCommand,
) error {
	return projectRepository.DeleteApiKey(command.ApiKeyID)
}

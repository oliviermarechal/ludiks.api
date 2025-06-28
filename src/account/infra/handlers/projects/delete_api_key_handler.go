package projects

import (
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/delete_api_key"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteApiKeyHandler struct {
	projectRepository domain_repositories.ProjectRepository
}

func NewDeleteApiKeyHandler(
	projectRepo domain_repositories.ProjectRepository,
) *DeleteApiKeyHandler {
	return &DeleteApiKeyHandler{
		projectRepository: projectRepo,
	}
}

func (h *DeleteApiKeyHandler) Handle(c *gin.Context) {
	apiKeyID := c.Param("api-key-id")

	err := delete_api_key.DeleteApiKeyUseCase(
		h.projectRepository,
		delete_api_key.DeleteApiKeyCommand{
			ApiKeyID: apiKeyID,
		},
	)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

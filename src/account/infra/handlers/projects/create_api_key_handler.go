package projects

import (
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/create_api_key"
	"ludiks/src/kernel/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateApiKeyHandler struct {
	projectRepository domain_repositories.ProjectRepository
}

func NewCreateApiKeyHandler(
	projectRepo domain_repositories.ProjectRepository,
) *CreateApiKeyHandler {
	return &CreateApiKeyHandler{
		projectRepository: projectRepo,
	}
}

func (h *CreateApiKeyHandler) Handle(c *gin.Context) {
	projectID := c.Param("id")

	var dto create_api_key.CreateApiKeyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := dto.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	createdApiKey, err := create_api_key.CreateApiKeyUseCase(
		h.projectRepository,
		create_api_key.CreateApiKeyCommand{
			ProjectID: projectID,
			Name:      dto.Name,
		},
	)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdApiKey)
}

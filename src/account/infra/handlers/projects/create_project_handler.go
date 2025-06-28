package projects

import (
	"errors"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/create_project"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateProjectHandler struct {
	projectRepository domain_repositories.ProjectRepository
}

func NewCreateProjectHandler(
	projectRepository domain_repositories.ProjectRepository,
) *CreateProjectHandler {
	return &CreateProjectHandler{
		projectRepository: projectRepository,
	}
}

func (h *CreateProjectHandler) Handle(c *gin.Context) {
	var dto create_project.CreateProjectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := dto.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	user_id, ok := c.Get("user_id")
	if !ok {
		handlers.HandleBadRequest(c, errors.New("authentication required"))
		return
	}

	registrationResult, err := create_project.CreateProjectUseCase(
		h.projectRepository,
		create_project.CreateProjectCommand{
			UserID: user_id.(string),
			Name:   dto.Name,
		},
	)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, registrationResult)
}

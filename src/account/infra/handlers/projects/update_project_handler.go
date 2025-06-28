package projects

import (
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/update_project"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateProjectHandler struct {
	projectRepository domain_repositories.ProjectRepository
}

func NewUpdateProjectHandler(
	projectRepository domain_repositories.ProjectRepository,
) *UpdateProjectHandler {
	return &UpdateProjectHandler{
		projectRepository: projectRepository,
	}
}

func (h *UpdateProjectHandler) Handle(c *gin.Context) {
	var dto update_project.UpdateProjectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := dto.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	projectID := c.Param("id")

	registrationResult, err := update_project.UpdateProjectUseCase(
		h.projectRepository,
		update_project.UpdateProjectCommand{
			ProjectID: projectID,
			Name:      dto.Name,
		},
	)

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, registrationResult)
}

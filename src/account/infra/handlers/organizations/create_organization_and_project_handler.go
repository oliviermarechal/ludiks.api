package organizations

import (
	"errors"
	domain_repositories "ludiks/src/account/domain/repositories"
	"ludiks/src/account/use_cases/command/create_organization_and_project"
	"ludiks/src/kernel/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateOrganizationAndProjectHandler struct {
	userRepository         domain_repositories.UserRepository
	projectRepository      domain_repositories.ProjectRepository
	organizationRepository domain_repositories.OrganizationRepository
}

func NewCreateOrganizationAndProjectHandler(
	userRepository domain_repositories.UserRepository,
	projectRepository domain_repositories.ProjectRepository,
	organizationRepository domain_repositories.OrganizationRepository,
) *CreateOrganizationAndProjectHandler {
	return &CreateOrganizationAndProjectHandler{
		userRepository:         userRepository,
		projectRepository:      projectRepository,
		organizationRepository: organizationRepository,
	}
}

func (h *CreateOrganizationAndProjectHandler) Handle(c *gin.Context) {
	var dto create_organization_and_project.CreateOrganizationAndProjectDTO
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

	registrationResult, err := create_organization_and_project.CreateOrganizationAndProjectUseCase(
		h.projectRepository,
		h.organizationRepository,
		h.userRepository,
		create_organization_and_project.CreateOrganizationAndProjectCommand{
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

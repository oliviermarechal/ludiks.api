package end_user_handler

import (
	"ludiks/src/kernel/app/handlers"
	domain_repositories "ludiks/src/tracking/domain/repositories"
	"ludiks/src/tracking/use_cases/command/log_end_user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateEndUserHandler struct {
	endUserRepository  domain_repositories.EndUserRepository
	metadataRepository domain_repositories.MetadataRepository
}

func NewCreateEndUserHandler(endUserRepository domain_repositories.EndUserRepository, metadataRepository domain_repositories.MetadataRepository) *CreateEndUserHandler {
	return &CreateEndUserHandler{endUserRepository: endUserRepository, metadataRepository: metadataRepository}
}

func (h *CreateEndUserHandler) Handle(c *gin.Context) {
	projectID, _ := c.Get("project_id")
	var createEndUserDTO log_end_user.LogEndUserDTO
	if err := c.ShouldBindJSON(&createEndUserDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := createEndUserDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	useCase := log_end_user.NewLogEndUserUseCase(h.endUserRepository, h.metadataRepository)
	endUser, err := useCase.Execute(log_end_user.LogEndUserCommand{
		ID:        createEndUserDTO.ID,
		FullName:  createEndUserDTO.FullName,
		Email:     createEndUserDTO.Email,
		Picture:   createEndUserDTO.Picture,
		ProjectID: projectID.(uuid.UUID).String(),
		Metadata:  createEndUserDTO.Metadata,
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, endUser)
}

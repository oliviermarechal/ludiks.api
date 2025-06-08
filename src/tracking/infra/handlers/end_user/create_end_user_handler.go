package end_user_handler

import (
	"ludiks/src/kernel/handlers"
	domain_repositories "ludiks/src/tracking/domain/repositories"
	"ludiks/src/tracking/use_cases/command/create_end_user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateEndUserHandler struct {
	endUserRepository domain_repositories.EndUserRepository
}

func NewCreateEndUserHandler(endUserRepository domain_repositories.EndUserRepository) *CreateEndUserHandler {
	return &CreateEndUserHandler{endUserRepository: endUserRepository}
}

func (h *CreateEndUserHandler) Handle(c *gin.Context) {
	projectID, _ := c.Get("project_id")
	var createEndUserDTO create_end_user.CreateEndUserDTO
	if err := c.ShouldBindJSON(&createEndUserDTO); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	if err := createEndUserDTO.Validate(); err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	useCase := create_end_user.NewCreateEndUserUseCase(h.endUserRepository)
	endUser, err := useCase.Execute(create_end_user.CreateEndUserCommand{
		ID:        createEndUserDTO.ID,
		FullName:  createEndUserDTO.FullName,
		Email:     createEndUserDTO.Email,
		Picture:   createEndUserDTO.Picture,
		ProjectID: projectID.(uuid.UUID).String(),
	})

	if err != nil {
		handlers.HandleBadRequest(c, err)
		return
	}

	c.JSON(http.StatusCreated, endUser)
}

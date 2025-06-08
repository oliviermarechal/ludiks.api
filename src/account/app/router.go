package account

import (
	auth_handlers "ludiks/src/account/infra/handlers/auth"
	projects_handlers "ludiks/src/account/infra/handlers/projects"
	providers "ludiks/src/account/infra/providers"
	infra_repositories "ludiks/src/account/infra/repositories"
	"ludiks/src/kernel/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepository := infra_repositories.NewUserRepository(db)
	projectRepository := infra_repositories.NewProjectRepository(db)
	encrypter := providers.NewEncrypter()
	jwtProvider := providers.NewJwtProvider()

	accounts := router.Group("/accounts")
	{
		accounts.GET("/me", middleware.JwtMiddleware, auth_handlers.NewMeHandler(db).Handle)
		accounts.POST("/registration", auth_handlers.NewRegistrationHandler(userRepository, encrypter, jwtProvider).Handle)
		accounts.POST("/login", auth_handlers.NewLoginHandler(userRepository, encrypter, jwtProvider).Handle)
		accounts.POST("/google-auth", auth_handlers.NewGoogleAuthHandler(userRepository, jwtProvider).Handle)
	}

	projects := router.Group("/projects", middleware.JwtMiddleware)
	{
		projects.POST("", projects_handlers.NewCreateProjectHandler(userRepository, projectRepository).Handle)
		projects.POST("/:id/api-keys", projects_handlers.NewCreateApiKeyHandler(projectRepository).Handle)
		projects.GET("", projects_handlers.NewListProjectsHandler(db).Handle)
		projects.GET("/:id/api-keys", projects_handlers.NewListProjectApiKeysHandler(db).Handle)
		projects.GET("/:id/overview", projects_handlers.NewProjectOverviewHandler(db).Handle)
	}
}

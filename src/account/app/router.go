package account

import (
	auth_handlers "ludiks/src/account/infra/handlers/auth"
	organizations_handlers "ludiks/src/account/infra/handlers/organizations"
	projects_handlers "ludiks/src/account/infra/handlers/projects"
	providers "ludiks/src/account/infra/providers"
	infra_repositories "ludiks/src/account/infra/repositories"
	"ludiks/src/kernel/app/middleware"
	kernel_providers "ludiks/src/kernel/infra/providers"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepository := infra_repositories.NewUserRepository(db)
	projectRepository := infra_repositories.NewProjectRepository(db)
	organizationRepository := infra_repositories.NewOrganizationRepository(db)
	mailerProvider := kernel_providers.NewMailerProvider(os.Getenv("MAILER_SEND_API_KEY"))

	encrypter := providers.NewEncrypter()
	jwtProvider := providers.NewJwtProvider()

	accounts := router.Group("/accounts")
	{
		accounts.GET("/invite/:id", auth_handlers.NewFindInviteHandler(db).Handle)
		accounts.GET("/me", middleware.JwtMiddleware, auth_handlers.NewMeHandler(db).Handle)
		accounts.POST("/registration", auth_handlers.NewRegistrationHandler(userRepository, organizationRepository, encrypter, jwtProvider).Handle)
		accounts.POST("/login", auth_handlers.NewLoginHandler(userRepository, encrypter, jwtProvider).Handle)
		accounts.POST("/google-auth", auth_handlers.NewGoogleAuthHandler(userRepository, organizationRepository, jwtProvider).Handle)
	}

	organizations := router.Group("/organizations", middleware.JwtMiddleware)
	{
		organizations.POST("", organizations_handlers.NewCreateOrganizationAndProjectHandler(userRepository, projectRepository, organizationRepository).Handle)
		organizations.GET("", organizations_handlers.NewListOrganizationsHandler(db).Handle)
		organizations.GET("user/invites", organizations_handlers.NewListUserReceivedInvitesHandler(db).Handle)
		organizations.GET("/:id/memberships", organizations_handlers.NewListOrganizationMembershipsHandler(db).Handle)
		organizations.POST("/:id/memberships", organizations_handlers.NewInvitMembershipHandler(organizationRepository, userRepository, mailerProvider).Handle)
		organizations.GET("/:id/invites", organizations_handlers.NewListInvitesHandler(db).Handle)
		organizations.DELETE("/:id/invites/:invit-id", organizations_handlers.NewCancelInvitHandler(organizationRepository).Handle)
		organizations.POST("/:id/invites/:invit-id/accept", organizations_handlers.NewAcceptInvitHandler(organizationRepository, userRepository).Handle)
		organizations.POST("/:id/invites/:invit-id/reject", organizations_handlers.NewRejectInvitHandler(organizationRepository, userRepository).Handle)
	}

	projects := router.Group("/projects", middleware.JwtMiddleware)
	{
		projects.POST("", projects_handlers.NewCreateProjectHandler(projectRepository).Handle)
		projects.PUT("/:id", projects_handlers.NewUpdateProjectHandler(projectRepository).Handle)
		projects.POST("/:id/api-keys", projects_handlers.NewCreateApiKeyHandler(projectRepository).Handle)
		projects.DELETE("/:id/api-keys/:api-key-id", projects_handlers.NewDeleteApiKeyHandler(projectRepository).Handle)
		projects.GET("/:id/metadata", projects_handlers.NewListProjectMetadataHandler(db).Handle)
		projects.GET("/:id/api-keys", projects_handlers.NewListProjectApiKeysHandler(db).Handle)
		projects.GET("/:id/overview", projects_handlers.NewProjectOverviewHandler(db).Handle)
	}
}

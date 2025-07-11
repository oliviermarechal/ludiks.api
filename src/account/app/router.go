package account

import (
	"ludiks/config"
	auth_handlers "ludiks/src/account/infra/handlers/auth"
	organizations_handlers "ludiks/src/account/infra/handlers/organizations"
	projects_handlers "ludiks/src/account/infra/handlers/projects"
	subscriptions_handlers "ludiks/src/account/infra/handlers/subscriptions"
	providers "ludiks/src/account/infra/providers"
	infra_repositories "ludiks/src/account/infra/repositories"
	"ludiks/src/kernel/app/middleware"
	kernel_providers "ludiks/src/kernel/infra/providers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepository := infra_repositories.NewUserRepository(db)
	projectRepository := infra_repositories.NewProjectRepository(db)
	organizationRepository := infra_repositories.NewOrganizationRepository(db)
	organizationSubscriptionRepository := infra_repositories.NewOrganizationSubscriptionRepository(db)
	invoiceRepository := infra_repositories.NewInvoiceRepository(db)

	mailerProvider := kernel_providers.NewMailerProvider(config.AppConfig.MailerSendAPIKey)
	encrypter := providers.NewEncrypter()
	jwtProvider := providers.NewJwtProvider()
	stripeProvider := providers.NewStripeProvider(config.AppConfig.StripeSecretKey)

	router.POST("/webhook", subscriptions_handlers.NewStripeWebhookHandler(
		organizationSubscriptionRepository,
		organizationRepository,
		invoiceRepository,
		config.AppConfig.StripeWebhookKey,
	).Handle)

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
		projects.GET("/:id/end-users", projects_handlers.NewListProjectEndUserHandler(db).Handle)
		projects.GET("/:id/metadata", projects_handlers.NewListProjectMetadataHandler(db).Handle)
		projects.GET("/:id/api-keys", projects_handlers.NewListProjectApiKeysHandler(db).Handle)
		projects.GET("/:id/overview", projects_handlers.NewProjectOverviewHandler(db).Handle)
	}

	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("", middleware.JwtMiddleware, subscriptions_handlers.NewCreateSubscriptionHandler(stripeProvider, organizationRepository).Handle)
		subscriptions.GET("setup-intent-success", subscriptions_handlers.NewSetupIntentSuccessHandler(stripeProvider, organizationSubscriptionRepository, organizationRepository).Handle)
		subscriptions.GET("organizations/:id", middleware.JwtMiddleware, subscriptions_handlers.NewGetOrganizationSubscriptionHandler(db).Handle)
		subscriptions.POST("/:id/cancel", middleware.JwtMiddleware, subscriptions_handlers.NewCancelSubscriptionHandler(stripeProvider, organizationSubscriptionRepository, organizationRepository).Handle)
	}
}

package tracking

import (
	"ludiks/config"
	"ludiks/src/kernel/app/middleware"
	end_user_handler "ludiks/src/tracking/infra/handlers/end_user"
	tracking_handler "ludiks/src/tracking/infra/handlers/tracking"
	infra_providers "ludiks/src/tracking/infra/providers"
	infra_repositories "ludiks/src/tracking/infra/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(router *gin.RouterGroup, db *gorm.DB) {
	endUserRepository := infra_repositories.NewEndUserRepository(db)
	progressionRepository := infra_repositories.NewProgressionRepository(db)
	circuitRepository := infra_repositories.NewCircuitRepository(db)
	organizationRepository := infra_repositories.NewOrganizationRepository(db)
	metadataRepository := infra_repositories.NewMetadataRepository(db)
	billingUsageProvider := infra_providers.NewBillingUsageProvider(config.AppConfig.StripeSecretKey)

	endUserRouter := router.Group("/end-user")
	endUserRouter.Use(middleware.NewApiKeyMiddleware(db))
	{
		endUserRouter.POST("", end_user_handler.NewCreateEndUserHandler(endUserRepository, metadataRepository).Handle)
	}

	trackingRouter := router.Group("/tracking")
	trackingRouter.Use(middleware.NewApiKeyMiddleware(db))
	{
		trackingRouter.POST("", tracking_handler.NewProgressionTrackingHandler(progressionRepository, endUserRepository, circuitRepository, organizationRepository, billingUsageProvider).Handle)
	}
}

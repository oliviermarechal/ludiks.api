package gamification

import (
	circuit_handler "ludiks/src/gamification/infra/handlers/circuit"
	reward_handler "ludiks/src/gamification/infra/handlers/reward"
	infra_repositories "ludiks/src/gamification/infra/repositories"
	"ludiks/src/kernel/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.RouterGroup, db *gorm.DB) {
	circuitRepository := infra_repositories.NewCircuitRepository(db)
	rewardRepository := infra_repositories.NewRewardRepository(db)

	circuitRouter := router.Group("/circuits", middleware.JwtMiddleware)
	{
		circuitRouter.POST("", circuit_handler.NewCreateCircuitHandler(circuitRepository).Handle)
		circuitRouter.POST("/:id/steps", circuit_handler.NewAddStepHandler(circuitRepository).Handle)
		circuitRouter.POST("/:id/rename", circuit_handler.NewRenameCircuitHandler(circuitRepository).Handle)
		circuitRouter.POST("/:id/activate", circuit_handler.NewActivateCircuitHandler(circuitRepository).Handle)
		circuitRouter.POST("/:id/set-steps", circuit_handler.NewSetCircuitStepsHandler(circuitRepository).Handle)
		circuitRouter.POST("/:id/generate-steps", circuit_handler.NewGenerateStepHandler(circuitRepository).Handle)
		circuitRouter.PUT("/:id/steps/:stepId", circuit_handler.NewUpdateStepHandler(circuitRepository).Handle)
		circuitRouter.DELETE("/:id/steps/:stepId", circuit_handler.NewDeleteStepHandler(circuitRepository).Handle)
		circuitRouter.GET("", circuit_handler.NewListCircuitHandler(db).Handle)

		rewardRouter := circuitRouter.Group("/:id/rewards", middleware.JwtMiddleware)
		{
			rewardRouter.GET("", reward_handler.NewListRewardHandler(db).Handle)
			rewardRouter.POST("", reward_handler.NewCreateRewardHandler(rewardRepository).Handle)
			rewardRouter.PUT("/:rewardId", reward_handler.NewUpdateRewardHandler(rewardRepository).Handle)
			rewardRouter.DELETE("/:rewardId", reward_handler.NewDeleteRewardHandler(rewardRepository).Handle)
		}
	}
}

package main

import (
	"fmt"
	"log"
	"ludiks/config"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ludiks/src/account/infra/providers"
	infra_repositories "ludiks/src/account/infra/repositories"
	"ludiks/src/account/use_cases/command/on_subscription_ended"
)

func main() {
	config.LoadConfig()

	// Valider la configuration
	if err := config.ValidateConfig(); err != nil {
		log.Fatalf("❌ Erreur de configuration: %v", err)
	}
	log.Println("✅ Configuration validée avec succès")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	organizationSubscriptionRepo := infra_repositories.NewOrganizationSubscriptionRepository(db)
	organizationRepo := infra_repositories.NewOrganizationRepository(db)
	stripeProvider := providers.NewStripeProvider(config.AppConfig.StripeSecretKey)

	useCase := on_subscription_ended.NewOnSubscriptionEndedUseCase(
		stripeProvider,
		organizationSubscriptionRepo,
		organizationRepo,
	)

	log.Println("Starting subscription sync job...")

	err = useCase.Execute()
	if err != nil {
		log.Printf("Subscription sync failed: %v", err)
		os.Exit(1)
	}

	log.Println("Subscription sync completed successfully")
}

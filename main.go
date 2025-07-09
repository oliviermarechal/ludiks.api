package main

import (
	"fmt"
	"log"
	"ludiks/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	kernel "ludiks/src/kernel"
	database "ludiks/src/kernel/app/database"
)

func main() {
	config.LoadConfig()

	// Valider la configuration en production
	if config.AppConfig.GoEnv == "production" {
		if err := config.ValidateConfig(); err != nil {
			log.Fatalf("❌ Erreur de configuration: %v", err)
		}
		log.Println("✅ Configuration validée avec succès")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Erreur lors de la migration: %v", err)
	}

	router := kernel.SetupRouter(db)

	log.Println("Serveur démarré sur :" + config.AppConfig.AppPort)
	router.Run(":" + config.AppConfig.AppPort)
}

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	AppPort    string
}

var AppConfig Config

func LoadConfig() {
	envFile := ".env"
	if os.Getenv("GO_ENV") == "test" {
		envFile = ".env.test"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Aucun fichier %s trouvé, utilisation des variables système", envFile)
	}

	AppConfig = Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "mydatabase"),
		DBPort:     getEnv("DB_PORT", "5432"),
		AppPort:    getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

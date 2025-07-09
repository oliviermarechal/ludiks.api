package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string

	AppPort string
	GoEnv   string

	StripeSecretKey           string
	StripeSubscriptionPriceID string
	StripeWebhookKey          string

	MailerSendAPIKey string
	FreeEventsLimit  int
	JWTSecretKey     string
	GoogleClientID   string

	FrontURL string
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
		DBHost:     getRequiredEnv("DB_HOST"),
		DBUser:     getRequiredEnv("DB_USER"),
		DBPassword: getRequiredEnv("DB_PASSWORD"),
		DBName:     getRequiredEnv("DB_NAME"),
		DBPort:     getRequiredEnv("DB_PORT"),

		AppPort: getEnv("APP_PORT", "8080"),
		GoEnv:   getEnv("GO_ENV", "development"),

		StripeSecretKey:           getRequiredEnv("STRIPE_SECRET_KEY"),
		StripeSubscriptionPriceID: getRequiredEnv("STRIPE_SUBSCRIPTION_PRICE_ID"),
		StripeWebhookKey:          getRequiredEnv("STRIPE_WEBHOOK_KEY"),

		JWTSecretKey: getRequiredEnv("JWT_SECRET_KEY"),

		MailerSendAPIKey: getRequiredEnv("MAILER_SEND_API_KEY"),

		GoogleClientID: getRequiredEnv("GOOGLE_CLIENT_ID"),

		FrontURL: getRequiredEnv("FRONT_URL"),

		FreeEventsLimit: getIntEnv("FREE_EVENTS_LIMIT", 5000),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getRequiredEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}

	if os.Getenv("GO_ENV") == "test" {
		log.Printf("Variable d'environnement %s non définie en mode test, utilisation d'une valeur par défaut", key)
		return getDefaultValue(key)
	}

	log.Fatalf("Variable d'environnement requise %s n'est pas définie", key)
	return ""
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Variable d'environnement %s n'est pas un nombre valide, utilisation de la valeur par défaut %d", key, defaultValue)
	}
	return defaultValue
}

func getDefaultValue(key string) string {
	defaults := map[string]string{
		"DB_HOST":                      "localhost",
		"DB_USER":                      "postgres",
		"DB_PASSWORD":                  "password",
		"DB_NAME":                      "mydatabase",
		"DB_PORT":                      "5432",
		"STRIPE_SECRET_KEY":            "sk_test_...",
		"STRIPE_SUBSCRIPTION_PRICE_ID": "price_...",
		"STRIPE_WEBHOOK_KEY":           "stripe_webhook_key",
		"JWT_SECRET_KEY":               "test-secret-key",
		"MAILER_SEND_API_KEY":          "test-api-key",
		"GOOGLE_CLIENT_ID":             "test-client-id",
		"FRONT_URL":                    "http://localhost:3000",
	}

	if value, exists := defaults[key]; exists {
		return value
	}

	return ""
}

func ValidateConfig() error {
	requiredVars := []string{
		"DB_HOST",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_PORT",
		"STRIPE_SECRET_KEY",
		"STRIPE_SUBSCRIPTION_PRICE_ID",
		"STRIPE_WEBHOOK_KEY",
		"JWT_SECRET_KEY",
		"MAILER_SEND_API_KEY",
		"GOOGLE_CLIENT_ID",
		"FRONT_URL",
	}

	var missingVars []string

	for _, varName := range requiredVars {
		if value := os.Getenv(varName); value == "" {
			missingVars = append(missingVars, varName)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("variables d'environnement manquantes: %v", missingVars)
	}

	return nil
}

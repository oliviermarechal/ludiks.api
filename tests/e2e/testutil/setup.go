package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"ludiks/config"
	"ludiks/src/account/domain/models"
	kernel "ludiks/src/kernel"
	"ludiks/src/kernel/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestEnv struct {
	DB        *gorm.DB
	container testcontainers.Container
}

func NewTestEnv() *TestEnv {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container: %v", err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	os.Setenv("GO_ENV", "test")
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_NAME", "testdb")

	config.LoadConfig()

	dsn := "host=" + config.AppConfig.DBHost +
		" user=" + config.AppConfig.DBUser +
		" password=" + config.AppConfig.DBPassword +
		" dbname=" + config.AppConfig.DBName +
		" port=" + config.AppConfig.DBPort +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	return &TestEnv{
		DB:        db,
		container: container,
	}
}

func (env *TestEnv) Cleanup() {
	if env.container != nil {
		env.container.Terminate(context.Background())
	}
}

func (env *TestEnv) NewRouter() *gin.Engine {
	return kernel.SetupRouter(env.DB)
}

func (env *TestEnv) ResetDatabase() {
	// Désactive temporairement les contraintes de clés étrangères
	env.DB.Exec("SET session_replication_role = 'replica';")

	// Récupère toutes les tables de la base de données
	var tables []string
	env.DB.Raw(`
		SELECT tablename FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename != 'schema_migrations'
	`).Scan(&tables)

	// Tronque toutes les tables en une seule transaction
	env.DB.Transaction(func(tx *gorm.DB) error {
		for _, table := range tables {
			if err := tx.Exec("TRUNCATE TABLE " + table + " CASCADE").Error; err != nil {
				return err
			}
		}
		return nil
	})

	// Réactive les contraintes de clés étrangères
	env.DB.Exec("SET session_replication_role = 'origin';")
}

func (env *TestEnv) ExecSQL(sql string, values ...interface{}) {
	err := env.DB.Exec(sql, values...).Error
	if err != nil {
		log.Fatalf("Erreur lors de l'exécution de la requête SQL: %v", err)
	}
}

func (env *TestEnv) GenerateTestUserAndToken() (string, string) {
	password := "password"
	user := models.User{
		ID:       uuid.New().String(),
		Email:    "test@user.com",
		Password: &password,
	}

	if err := env.DB.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create test user: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	return tokenString, user.ID
}

func (env *TestEnv) AuthenticatedRequest(method, path string, body interface{}, token string) *http.Request {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	return req
}

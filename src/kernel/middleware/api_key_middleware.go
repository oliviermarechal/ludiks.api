package middleware

import (
	"net/http"
	"strings"

	"ludiks/src/kernel/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewApiKeyMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			return
		}

		apiKey = strings.TrimPrefix(apiKey, "Bearer ")

		var apiKeyEntity database.ApiKeyEntity
		result := db.Where("value = ?", apiKey).First(&apiKeyEntity)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		c.Set("project_id", apiKeyEntity.ProjectID)
		c.Next()
	}
}

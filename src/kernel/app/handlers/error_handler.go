package handlers

import (
	"github.com/gin-gonic/gin"
)

func HandleBadRequest(c *gin.Context, err error) {
	c.JSON(400, gin.H{"error": err.Error()})
}

func HandleUnauthorized(c *gin.Context, err error) {
	c.JSON(401, gin.H{"error": err.Error()})
}

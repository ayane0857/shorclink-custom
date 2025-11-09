package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RequireAPIToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("API") != "required" {
			c.Next()
			return
		}
		token := c.GetHeader("X-API-Token")
		if token != os.Getenv("API_TOKEN") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RequireAPIToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("API") == "invalid" {
			c.JSON(http.StatusForbidden, gin.H{"error": "API access is disabled"})
			c.Abort()
			return
		}
		if os.Getenv("API") != "required" {
			c.Next()
			return
		}
		token := c.GetHeader("X-API-Token")

		if token == "" && c.Request.Method == "GET" {
			token = c.Query("X-API-Token")
		}
		
		if token != os.Getenv("API_TOKEN") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

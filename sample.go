package main

import (
	"net/http"
	"os"
	"shorclick/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	os.Setenv("SHORT_CODE", "required")

	r := gin.Default()
	r.Use(handlers.SetShortCode())

	// POSTエンドポイントに変更
	r.POST("/test", func(c *gin.Context) {
		shortCode, _ := c.Get("short_code")
		c.JSON(http.StatusOK, gin.H{"short_code": shortCode})
	})

	r.Run(":8000")
}
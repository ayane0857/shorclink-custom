package main

import (
	"log"
	"shorclick/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r:= gin.Default()
	log.Println("Starting server on :8080")
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "welcome to shorclick",
		})
	})
	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "don't setup api token",
		})
	})
	r.GET("/:id", handlers.ShortLink())
	r.Run(":8080")
}
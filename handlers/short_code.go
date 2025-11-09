package handlers

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"shorclick/models"

	"github.com/gin-gonic/gin"
)

func SetShortCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("SHORT_CODE") == "required" {
			RequiredShortCode()(c)
		} else if os.Getenv("SHORT_CODE") == "auto" {
			length := 8
			if l := os.Getenv("SHORT_CODE_LENGTH"); l != "" {
				if val, err := strconv.Atoi(l); err == nil {
					length = val
				}
			}
			AutoShortCode(length)(c)
		} else {
			GenerateShortCode(8)(c)
		}
		c.Next()
		return
	}
}

func RequiredShortCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ShortLink
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			c.Abort()
			return
		}
		var short_code = req.ShortCode
		if short_code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
			c.Abort()
			return
		}
		c.Set("short_code", short_code)
		return
	}
}
func AutoShortCode(length int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ShortLink
		ShortCode := req.ShortCode
		if ShortCode == ""{
			var code = generateRandomString(length)
			c.Set("short_code", code)
			return
		}
		c.Set("short_code", ShortCode)
		return
	}
}
func GenerateShortCode(length int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var code = generateRandomString(length)
		c.Set("short_code", code)
		return
	}
}

func generateRandomString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		b[i] = chars[n.Int64()]
	}
	return string(b)
}

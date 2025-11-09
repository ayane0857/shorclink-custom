package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"shorclick/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetShortLinks(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        limitStr := c.DefaultQuery("limit", "10")
        offsetStr := c.DefaultQuery("offset", "0")

        limit, err := strconv.Atoi(limitStr)
        if err != nil || limit <= 0 || limit > 100 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
            return
        }
        offset, err := strconv.Atoi(offsetStr)
        if err != nil || offset < 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
            return
        }

        var ShortLinks []models.ShortLink
        if err := db.Limit(limit).Offset(offset).Find(&ShortLinks).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve ShortLinks"})
            return
        }
        if len(ShortLinks) == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "No ShortLinks found"})
            return
        }
        c.JSON(http.StatusOK, ShortLinks)
    }
}

func GetShortLink(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var ShortLink models.ShortLink
        id := c.Param("id")

        if err := db.First(&ShortLink, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "ShortLink not found"})
        return
        }

        c.JSON(http.StatusOK, ShortLink)
    }
}

func PostShortLink(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req models.ShortLink
		if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
            return
        }
		if req.URL  == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
			return
		}
		if !strings.HasPrefix(req.URL , "http://") && !strings.HasPrefix(req.URL , "https://") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL must start with http:// or https://"})
			return
		}

        if err := db.Create(&req).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ShortLink"})
            return
        }
        c.JSON(http.StatusOK, req)
    }
}

func PutShortLink(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var ShortLink models.ShortLink
        id := c.Param("id")

        if err := db.First(&ShortLink, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "ShortLink not found"})
            return
        }

        if err := c.ShouldBindJSON(&ShortLink); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        db.Save(&ShortLink)
        c.JSON(http.StatusOK, ShortLink)
    }
}

func DeleteShortLink(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var ShortLink models.ShortLink

        id := c.Param("id")

        if err := db.First(&ShortLink, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "ShortLink not found"})
            return
        }

        db.Delete(&ShortLink)
        c.JSON(http.StatusOK, ShortLink)
    }
}
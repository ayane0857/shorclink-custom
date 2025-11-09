package handlers

import (
	"net/http"
	"shorclick/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RedirectShortLink(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var model models.ShortLink
        var id = c.Param("id")
        result := db.First(&model, "short_code = ?", id)
        if result.Error != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Short link not found"})
            return
        }
        c.Redirect(http.StatusMovedPermanently, model.URL)
        return
    }
}
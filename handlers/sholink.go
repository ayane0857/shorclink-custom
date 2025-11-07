package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShortLinkHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        var id = c.Param("id")
        c.Redirect(http.StatusMovedPermanently, "https://"+id)
        return
    }
}
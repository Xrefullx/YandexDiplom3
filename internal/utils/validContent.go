package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidContent(c *gin.Context, ContentType string) bool {
	if c.GetHeader("Content-Type") != ContentType {
		c.String(http.StatusUnsupportedMediaType, "ожидалось %s", ContentType)
		return false
	}
	return true
}

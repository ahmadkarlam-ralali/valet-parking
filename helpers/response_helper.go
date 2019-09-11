package helpers

import (
	"github.com/gin-gonic/gin"
)

func HttpError(c *gin.Context, message string, code int) {
	c.Abort()
	c.JSON(code, gin.H{"status": "error", "message": message})
}

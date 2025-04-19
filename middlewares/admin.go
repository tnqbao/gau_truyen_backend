package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		permission := c.MustGet("permission").(string)
		if permission != "admin" {
			c.JSON(403, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}

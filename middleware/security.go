package middleware

import (
	"github.com/gin-gonic/gin"
)

func HeaderPolicy(c *gin.Context) {
	c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-Frame-Options", "DENY")
	c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*") //change * to allowed origin
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")


	c.Next()
}

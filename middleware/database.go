package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Middleware to inject db into gin context
func InjectDBToContext(postgresDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("postgresDB", postgresDB)
		c.Next()
	}
}
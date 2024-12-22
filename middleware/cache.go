package middleware

import (
	"engine.multifinance.com/cache"
	"github.com/gin-gonic/gin"
)

func InjectCacheToContext(lruCache cache.DoubleBufferCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cache", lruCache)
		c.Next()
	}
}
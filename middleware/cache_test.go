package middleware

import (
	"engine.multifinance.com/cache"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestInjectCacheToContext tests the InjectCacheToContext middleware
func TestInjectCacheToContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	lruCache := cache.NewDoubleBufferLru(cache.DoubleBufferLruConfig{
		CacheSize:        3,
		CacheExpiryMSec:  10000,
		CacheRefreshMSec: 3000,
	}).(*cache.DoubleBuffer)

	// Create a Gin router
	router := gin.Default()

	// Apply the middleware
	router.Use(InjectCacheToContext(lruCache))

	// Test route to check if cache is injected into context
	router.GET("/test-cache", func(c *gin.Context) {
		injectedCache, exists := c.Get("cache")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cache not injected"})
			return
		}

		// Assert that the cache in context is the same as the one passed to the middleware
		if injectedCache != lruCache {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cache in context does not match"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cache injected successfully"})
	})

	// Perform the request to test the middleware
	req, _ := http.NewRequest(http.MethodGet, "/test-cache", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "Cache injected successfully"}`, w.Body.String())
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test cases
func TestValidateParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success - Valid Params", func(t *testing.T) {
		router := gin.Default()
		router.Use(ValidateParams)

		router.GET("/test/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "validated"})
		})

		// Simulate valid UUID
		req, _ := http.NewRequest("GET", "/test/550e8400-e29b-41d4-a716-446655440000", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.Code)
		}
	})

	t.Run("Failure - Missing Param", func(t *testing.T) {
		router := gin.Default()
		router.Use(ValidateParams)

		router.GET("/test/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "validated"})
		})

		// Simulate missing UUID
		req, _ := http.NewRequest("GET", "/test/", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code == http.StatusOK {
			t.Errorf("Expected validation failure, but got status 200")
		}
	})
}

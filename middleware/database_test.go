package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestInjectDBToContext tests the InjectDBToContext middleware
func TestInjectDBToContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Create an in-memory SQLite database (for testing purposes)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}

	// Create a Gin router
	router := gin.Default()

	// Apply the middleware
	router.Use(InjectDBToContext(db))

	// Test route to check if db is injected into context
	router.GET("/test-db", func(c *gin.Context) {
		// Retrieve the db from context
		injectedDB, exists := c.Get("postgresDB")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB not injected"})
			return
		}

		// Assert that the db in context is the same as the one passed to the middleware
		if injectedDB != db {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB in context does not match"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "DB injected successfully"})
	})

	// Perform the request to test the middleware
	req, _ := http.NewRequest(http.MethodGet, "/test-db", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "DB injected successfully"}`, w.Body.String())
}

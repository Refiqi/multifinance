package handlers

import (
	"bytes"
	"encoding/json"
	"engine.multifinance.com/cache"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"multifinance.com/multifinance/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Mock data for testing
var (
	mockTransaction = model.Transaction{
		ID:         "1",
		UserID:     "user123",
		OTR:        1,
		AdminFee:   1000,
		Installments: 999999,
		Interest:   5.5,
		AssetName:  "Car",
	}

	mockLoanLimit = model.LoanLimit{
		ID: "limit123",
		Tenor1:  100000,
		Tenor2:  150000,
		Tenor3:  200000,
		Tenor6:  250000,
	}

	mockUser = model.User{
		NIK:            "user123",
		FullName:       "Bob",
		LegalName:      "Bobby",
		PlaceOfBirth:   "Jakarta",
		DateOfBirth:    time.Time{},
		Salary:         10000,
		KTPPhotoURL:    "url",
		SelfiePhotoURL: "url",
		LoanLimitID:    "limit123",
		LoanLimit:      mockLoanLimit,
		Transactions:   nil,
	}
)

func setupDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.Transaction{}, &model.LoanLimit{}, &model.User{})
	return db
}

func TestCreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupDB()
	router := setupRouter(db)

	// Seed LoanLimit for test
	db.Create(&mockLoanLimit)
	db.Create(&mockUser)

	lruCache := cache.NewDoubleBufferLru(cache.DoubleBufferLruConfig{
		CacheSize:        3,
		CacheExpiryMSec:  10000,
		CacheRefreshMSec: 3000,
	}).(*cache.DoubleBuffer)

	router.POST("/transactions", func(c *gin.Context) {
		c.Set("postgresDB", db)
		c.Set("cache", lruCache)
		CreateTransaction(c)
	})

	t.Run("Success - Create Transaction Within Limit", func(t *testing.T) {
		db.Delete(&mockTransaction) // Clear transactions before each test

		body, _ := json.Marshal(mockTransaction)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		var createdTransaction model.Transaction
		_ = json.Unmarshal(resp.Body.Bytes(), &createdTransaction)
		assert.Equal(t, mockTransaction.UserID, createdTransaction.UserID)
	})


	t.Run("Fail - Exceeds Loan Limit 6 Month", func(t *testing.T) {
		db.Delete(&mockTransaction) // Clear transactions before each test

		overLimitTransaction := mockTransaction
		overLimitTransaction.OTR = 300001 // Exceeds Tenor 6 (250000)

		body, _ := json.Marshal(overLimitTransaction)
		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		var response map[string]interface{}
		_ = json.Unmarshal(resp.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "Transaction exceeds loan limit for 6 month")
	})


}

package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"multifinance.com/multifinance/model"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func setupTestDB() (*gorm.DB, sqlmock.Sqlmock) {
	// Create a mock database connection and mock the sql.DB interface
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	return gormDB, mock
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("postgresDB", db)
	})
	SetupLoanLimitRoutes(r)
	return r
}

func TestCreateLoanLimit(t *testing.T) {
	// Setup mock database and router
	db, mock := setupTestDB()
	router := setupRouter(db)

	// Test data
	limit := model.LoanLimit{ID: "1", Tenor1: 1000, Tenor2: 12}
	limitJSON, _ := json.Marshal(limit)

	// Mock the Create method
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"loan_limits\"").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req, _ := http.NewRequest(http.MethodPost, "/loan-limits", bytes.NewBuffer(limitJSON))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "1000")
}

func TestGetLoanLimits(t *testing.T) {
	// Setup mock database and router
	db, mock := setupTestDB()
	router := setupRouter(db)

	// Mock the Find method
	rows := sqlmock.NewRows([]string{"id", "tenor1", "tenor2"}).
		AddRow("1", 1000, 12).
		AddRow("2", 2000, 24)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"loan_limits\"")).WillReturnRows(rows)

	req, _ := http.NewRequest(http.MethodGet, "/loan-limits", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "1000")
	assert.Contains(t, w.Body.String(), "2000")
}

func TestGetLoanLimitByID(t *testing.T) {
	// Setup mock database and router
	db, mock := setupTestDB()
	router := setupRouter(db)

	// Mock the First method
	rows := sqlmock.NewRows([]string{"id", "tenor1", "tenor2"}).
		AddRow("1", 1000, 12)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"loan_limits\" WHERE id = $1 ORDER BY \"loan_limits\".\"id\" LIMIT $2")).WillReturnRows(rows)

	req, _ := http.NewRequest(http.MethodGet, "/loan-limit/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "1000")
}

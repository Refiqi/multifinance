package handlers


import (
	"multifinance.com/multifinance/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateLoanLimit handles creating a new loan limit
func CreateLoanLimit(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var limit model.LoanLimit
	if err := c.ShouldBindJSON(&limit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := postgresDB.Create(&limit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, limit)
}

// GetLoanLimits retrieves all loan limits
func GetLoanLimits(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var limits []model.LoanLimit
	if err := postgresDB.Find(&limits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, limits)
}

// GetLoanLimitByID retrieves a specific loan limit by ID
func GetLoanLimitByID(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var limit model.LoanLimit
	id := c.Param("id")
	if err := postgresDB.First(&limit, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan limit not found"})
		return
	}
	c.JSON(http.StatusOK, limit)
}

// UpdateLoanLimit updates an existing loan limit
func UpdateLoanLimit(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var limit model.LoanLimit
	id := c.Param("id")
	if err := postgresDB.First(&limit, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan limit not found"})
		return
	}
	if err := c.ShouldBindJSON(&limit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit.ID = id // Ensure ID remains the same
	if err := postgresDB.Save(&limit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, limit)
}

// DeleteLoanLimit deletes a loan limit by ID
func DeleteLoanLimit(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	id := c.Param("id")
	if err := postgresDB.Delete(&model.LoanLimit{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Loan limit deleted"})
}

// SetupLoanLimitRoutes sets up the API routes for loan limits
func SetupLoanLimitRoutes(router *gin.Engine) {
	router.POST("/loan-limits", CreateLoanLimit)
	router.GET("/loan-limits", GetLoanLimits)
	router.GET("/loan-limit/:id", GetLoanLimitByID)
	router.POST("/loan-limit/:id", UpdateLoanLimit)
	router.DELETE("/loan-limit/:id", DeleteLoanLimit)
}
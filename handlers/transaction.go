package handlers

import (
	"engine.multifinance.com/cache"
	engineError "engine.multifinance.com/error"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"multifinance.com/multifinance/model"
	"net/http"
)

// CreateTransaction creates a new transaction within a database transaction
func CreateTransaction(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	lruCache := c.MustGet("cache").(cache.DoubleBufferCache)
	var transaction model.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := postgresDB.Transaction(func(tx *gorm.DB) error {
		var currentTotal float64
		var loanLimit model.LoanLimit

		// Step 1: Calculate Current Transaction Total for User (WITHIN TRANSACTION)
		// get from cache if there is no cache query from DB and refresh the cache
		lruCache.Get("transaction_total_user_id_"+ transaction.UserID , func() (data interface{}, err error) {
			if err := tx.
				Model(&model.Transaction{}).
				Select("COALESCE(SUM(otr), 0)").
				Where("user_id = ?", transaction.UserID).
				Scan(&currentTotal).Error; err != nil {
				return nil, err
			}
			return currentTotal, nil
		})

		// Step 2: LEFT JOIN to Retrieve Loan Limit from User (WITHIN TRANSACTION)
		// get from cache if there is no cache query from DB and refresh the cache
		lruCache.Get("limit_user_id_"+ transaction.UserID, func() (data interface{}, err error) {
			if err := tx.
				Table("users").
				Select("loan_limits.*").
				Joins("LEFT JOIN loan_limits ON users.loan_limit_id = loan_limits.id").
				Where("users.nik = ?", transaction.UserID).
				Limit(1).
				Order("users.nik").
				Take(&loanLimit).Error; err != nil {
				return nil, err
			}

			return loanLimit, nil
		})


		// Step 3: Calculate Total with New Transaction
		newTotal := currentTotal + transaction.OTR

		// Step 4: Check Against Loan Limit From Tenor
		switch {
		case newTotal > loanLimit.Tenor6:
			return &engineError.TransactionError{
				Message: "Transaction exceeds loan limit for 6 month",
				Limit:   loanLimit.Tenor6,
				Current: currentTotal,
				New:     newTotal,
			}
		case newTotal > loanLimit.Tenor3:
			return &engineError.TransactionError{
				Message: "Transaction exceeds loan limit for 3 month",
				Limit:   loanLimit.Tenor3,
				Current: currentTotal,
				New:     newTotal,
			}
		case newTotal > loanLimit.Tenor2:
			return &engineError.TransactionError{
				Message: "Transaction exceeds loan limit for 2 month",
				Limit:   loanLimit.Tenor2,
				Current: currentTotal,
				New:     newTotal,
			}
		case newTotal > loanLimit.Tenor1:
			return &engineError.TransactionError{
				Message: "Transaction exceeds loan limit for 1 month",
				Limit:   loanLimit.Tenor1,
				Current: currentTotal,
				New:     newTotal,
			}
		}


		// Step 5: Create Transaction (WITHIN TRANSACTION)
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// Step 6: If no errors, commit the transaction
		return nil
	})

	if err != nil {
		if tErr, ok := err.(*engineError.TransactionError); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":         tErr.Message,
				"limit":         tErr.Limit,
				"current_total": tErr.Current,
				"new_total":     tErr.New,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		}
		return
	}

	// Success Response
	c.JSON(http.StatusCreated, transaction)
}

// GetTransactionsByUserID retrieves all transactions
func GetTransactionsByUserID(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var transactions []model.Transaction

	if err := postgresDB.Where("user_id = ?", c.Param("user_id")).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transactions)
}

// GetTransactionByID retrieves a specific transaction by ID
func GetTransactionByID(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var transaction model.Transaction
	id := c.Param("id")

	if err := postgresDB.First(&transaction, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction updates an existing transaction
func UpdateTransaction(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var transaction model.Transaction
	id := c.Param("id")

	if err := postgresDB.First(&transaction, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction.ID = id

	if err := postgresDB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction deletes a transaction by ID
func DeleteTransaction(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	id := c.Param("id")

	if err := postgresDB.Delete(&model.Transaction{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}

// SetupTransactionRoutes sets up the API routes
func SetupTransactionRoutes(router *gin.Engine) {
	router.POST("/transactions", CreateTransaction)
	router.GET("/transactions", GetTransactionsByUserID)
	router.GET("/transaction/:id", GetTransactionByID)
	router.POST("/transactions/:id", UpdateTransaction)
	router.DELETE("/transactions/:id", DeleteTransaction)
}

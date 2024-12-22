package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"multifinance.com/multifinance/model"
	"net/http"
)

// CreateUser handles creating a new user
func CreateUser(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := postgresDB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// GetUsers retrieves all users
func GetUsers(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var users []model.User

	if err := postgresDB.Preload("LoanLimits").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUser updates an existing user
func UpdateUser(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	var user model.User

	id := c.Param("id")
	if err := postgresDB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	postgresDB.Save(&user)
	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	postgresDB := c.MustGet("postgresDB").(*gorm.DB)
	id := c.Param("id")
	if err := postgresDB.Delete(&model.User {}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func SetupUserRoutes(router *gin.Engine) {
	router.POST("/user", CreateUser)
	router.GET("/user", GetUsers)
	router.POST("/user/:id", UpdateUser)
	router.DELETE("/user/:id", DeleteUser)
}
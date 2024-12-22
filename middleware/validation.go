package middleware

import (
	"engine.multifinance.com/validation"
	"github.com/gin-gonic/gin"
)

func ValidateParams(c *gin.Context) {

	validator := validation.NewValidation()
	validator.Validate(c.Params)

	c.Next()
}

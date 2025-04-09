// Package middleware contains API middleware components.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-crypto-bot-clean/backend/internal/api/dto/response"
)

// RecoveryMiddleware recovers from panics and returns a standardized error response.
//
//	@summary	Panic recovery middleware
//	@description	Recovers from panics, logs details, and returns 500 error.
func RecoveryMiddleware(logger Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered:", r)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse{
					Code:    "internal_error",
					Message: "Internal server error",
					Details: "A server error occurred. Please try again later.",
				})
			}
		}()
		c.Next()
	}
}

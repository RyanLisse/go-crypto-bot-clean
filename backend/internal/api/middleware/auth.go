// Package middleware contains API middleware components.
package middleware

import (
	"net/http"

	"go-crypto-bot-clean/backend/internal/api/dto/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates API keys from the X-API-Key header.
//
//	@summary	API key authentication middleware
//	@description	Validates API keys from the X-API-Key header.
//	@security	ApiKeyAuth
func AuthMiddleware(validAPIKeys map[string]struct{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if _, ok := validAPIKeys[apiKey]; !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{
				Code:    "unauthorized",
				Message: "Invalid or missing API key",
			})
			return
		}
		c.Next()
	}
}

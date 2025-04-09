package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ryanlisse/go-crypto-bot/internal/api/handlers"
	"github.com/ryanlisse/go-crypto-bot/internal/api/middleware"
	"github.com/ryanlisse/go-crypto-bot/internal/config"
)

// AuthDependencies contains the dependencies for authentication.
type AuthDependencies struct {
	// Handlers
	AuthHandler *handlers.AuthHandler

	// Configuration
	Config *config.Config
}

// NewAuthDependencies creates a new AuthDependencies instance.
func NewAuthDependencies(cfg *config.Config) *AuthDependencies {
	deps := &AuthDependencies{
		Config: cfg,
	}

	// Initialize auth handler
	deps.AuthHandler = handlers.NewAuthHandler(
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTExpiry,
		cfg.Auth.CookieName,
	)

	return deps
}

// SetupAuthRoutes adds authentication routes to the router.
func SetupAuthRoutes(router *gin.Engine, deps *AuthDependencies) {
	// Authentication endpoints
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", deps.AuthHandler.Login)
		authGroup.POST("/logout", deps.AuthHandler.Logout)

		// Protected auth endpoints
		authProtected := authGroup.Group("")
		if deps.Config.Auth.Enabled {
			authProtected.Use(middleware.JWTAuthMiddleware(deps.Config.Auth.JWTSecret))
		}
		authProtected.GET("/me", deps.AuthHandler.GetCurrentUser)
	}
}

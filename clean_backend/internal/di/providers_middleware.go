package di

import (
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/service"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// provideAuthService creates and returns an authentication service
func provideAuthService(config *config.Config, logger *zerolog.Logger) port.AuthServiceInterface {
	authServiceFactory := service.NewAuthServiceFactory(config, logger)
	return authServiceFactory.CreateAuthService()
}

// provideAuthMiddlewareFactory creates and returns an authentication middleware factory
func provideAuthMiddlewareFactory(authService port.AuthServiceInterface, config *config.Config, logger *zerolog.Logger) *middleware.AuthFactory {
	return middleware.NewAuthFactory(authService, config, logger)
}

// provideAuthMiddleware creates and returns the default authentication middleware
func provideAuthMiddleware(factory *middleware.AuthFactory) middleware.AuthMiddleware {
	return factory.CreateDefaultMiddleware()
}

// provideUnifiedErrorMiddleware creates and returns a unified error middleware
func provideUnifiedErrorMiddleware(logger *zerolog.Logger) *middleware.UnifiedErrorMiddleware {
	return middleware.NewUnifiedErrorMiddleware(logger)
}

// provideMiddlewares creates and returns all middleware providers
func provideMiddlewares(config *config.Config, logger *zerolog.Logger) *MiddlewareProviders {
	// Create auth service
	authService := provideAuthService(config, logger)
	
	// Create auth middleware factory
	authFactory := provideAuthMiddlewareFactory(authService, config, logger)
	
	// Create default auth middleware
	authMiddleware := provideAuthMiddleware(authFactory)
	
	// Create error middleware
	errorMiddleware := provideUnifiedErrorMiddleware(logger)
	
	return &MiddlewareProviders{
		AuthService:     authService,
		AuthFactory:     authFactory,
		AuthMiddleware:  authMiddleware,
		ErrorMiddleware: errorMiddleware,
	}
}

// MiddlewareProviders holds all middleware providers
type MiddlewareProviders struct {
	AuthService     port.AuthServiceInterface
	AuthFactory     *middleware.AuthFactory
	AuthMiddleware  middleware.AuthMiddleware
	ErrorMiddleware *middleware.UnifiedErrorMiddleware
}

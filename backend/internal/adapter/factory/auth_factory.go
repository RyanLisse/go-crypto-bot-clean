package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/repository/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/rs/zerolog"
	gormdb "gorm.io/gorm"
)

// AuthFactory creates authentication-related components
type AuthFactory struct {
	db     *gormdb.DB
	logger *zerolog.Logger
}

// NewAuthFactory creates a new AuthFactory
func NewAuthFactory(db *gormdb.DB, logger *zerolog.Logger) *AuthFactory {
	return &AuthFactory{
		db:     db,
		logger: logger,
	}
}

// CreateUserRepository creates a user repository
func (f *AuthFactory) CreateUserRepository() *gorm.UserRepository {
	return gorm.NewUserRepository(f.db)
}

// CreateUserService creates a user service
func (f *AuthFactory) CreateUserService() *service.UserService {
	userRepo := f.CreateUserRepository()
	return service.NewUserService(userRepo)
}

// CreateAuthService creates an authentication service
func (f *AuthFactory) CreateAuthService(secretKey string) (*service.AuthService, error) {
	userService := f.CreateUserService()
	return service.NewAuthService(userService, secretKey)
}

// CreateEnhancedClerkMiddleware creates an enhanced Clerk middleware
func (f *AuthFactory) CreateEnhancedClerkMiddleware(secretKey string) (*middleware.EnhancedClerkMiddleware, error) {
	authService, err := f.CreateAuthService(secretKey)
	if err != nil {
		return nil, err
	}
	return middleware.NewEnhancedClerkMiddleware(authService, f.logger), nil
}

// CreateClerkMiddleware creates a basic Clerk middleware
func (f *AuthFactory) CreateClerkMiddleware(secretKey string) *middleware.ClerkMiddleware {
	return middleware.NewClerkMiddleware(secretKey, f.logger)
}

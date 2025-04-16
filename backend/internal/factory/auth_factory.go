package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
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
func (f *AuthFactory) CreateUserRepository() *repo.UserRepository {
	return repo.NewUserRepository(f.db, f.logger)
}

// CreateUserService creates a user service
func (f *AuthFactory) CreateUserService() service.UserServiceInterface {
	userRepo := f.CreateUserRepository()
	return service.NewUserService(userRepo)
}

// CreateAuthService creates an authentication service
func (f *AuthFactory) CreateAuthService(secretKey string) (service.AuthServiceInterface, error) {
	userService := f.CreateUserService()
	return service.NewAuthService(userService.(*service.UserService), secretKey)
}


// CreateAuthMiddleware creates an authentication middleware
func (f *AuthFactory) CreateAuthMiddleware(secret string) middleware.AuthMiddleware {
	authService, _ := f.CreateAuthService(secret)
	return middleware.NewAuthMiddleware(authService, f.logger)
}

// CreateTestAuthMiddleware creates a test authentication middleware
func (f *AuthFactory) CreateTestAuthMiddleware() middleware.AuthMiddleware {
	return middleware.NewTestAuthMiddleware(f.logger)
}

// CreateDisabledAuthMiddleware creates a disabled authentication middleware
func (f *AuthFactory) CreateDisabledAuthMiddleware() middleware.AuthMiddleware {
	return middleware.NewDisabledAuthMiddleware(f.logger)
}

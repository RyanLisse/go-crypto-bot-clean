package factory

import (
"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/repository/memory"
"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/service"
"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
portservice "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/service"
"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/usecase"
"github.com/rs/zerolog"
)

// UseCaseFactory is responsible for creating use case instances with their dependencies.
// It acts as an abstraction over the direct instantiation within the DI container,
// potentially allowing for different implementations (e.g., mocks for testing)
type UseCaseFactory struct {
	Logger      *zerolog.Logger
	Config      *config.Config
	AuthService port.AuthServiceInterface
}

// NewUseCaseFactory creates a new UseCaseFactory
func NewUseCaseFactory(logger *zerolog.Logger, config *config.Config, authService port.AuthServiceInterface) *UseCaseFactory {
	return &UseCaseFactory{
		Logger:      logger,
		Config:      config,
		AuthService: authService,
	}
}

// BuildProtectedUserUseCase creates a new instance of the ProtectedUserUseCase.
// Note: It returns the ports.ProtectedUserService interface type, as defined by the use case constructor.
func (f *UseCaseFactory) BuildProtectedUserUseCase() portservice.ProtectedUserService {
	return usecase.NewProtectedUserUseCase(
f.AuthService,
)
}

// BuildAuthService returns the AuthService instance.
// This is a simple pass-through method to provide the AuthService to handlers.
func (f *UseCaseFactory) BuildAuthService() port.AuthServiceInterface {
	return f.AuthService
}

// BuildUserService creates a new instance of the UserService.
func (f *UseCaseFactory) BuildUserService() port.UserServiceInterface {
	// Create a memory repository for now
	userRepo := memory.NewUserRepository(f.Logger)

	// Create the user service
	return service.NewUserService(
f.Logger,
f.Config,
userRepo,
f.AuthService,
)
}

// TODO: Add factory methods for other potential use cases

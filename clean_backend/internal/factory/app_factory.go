package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
)

// AppFactory is a centralized factory for creating application components
type AppFactory struct {
	config *config.Config
	logger *zerolog.Logger

	// Cache for created components
	components map[string]interface{}
}

// NewAppFactory creates a new AppFactory
func NewAppFactory(config *config.Config, logger *zerolog.Logger) *AppFactory {
	return &AppFactory{
		config:     config,
		logger:     logger,
		components: make(map[string]interface{}),
	}
}

// IsMockAllowed checks if mock implementations are allowed
func (f *AppFactory) IsMockAllowed() bool {
	// Never allow mocks in production
	if f.config.Environment == "production" {
		return false
	}
	return f.config.Mock.Enabled
}

// ShouldUseMock checks if a specific component should use a mock implementation
func (f *AppFactory) ShouldUseMock(componentKey string) bool {
	// First check if mocks are globally allowed
	if !f.IsMockAllowed() {
		return false
	}

	// Then check component-specific setting
	switch componentKey {
	case "mexc_client":
		return f.config.Mock.MEXCClient
	case "wallet_service":
		return f.config.Mock.WalletService
	case "auth_service":
		return f.config.Mock.AuthService
	case "ai_service":
		return f.config.Mock.AIService
	default:
		// If we don't have a specific setting, default to false
		return false
	}
}

// LogMockUsage logs when a mock implementation is being used
func (f *AppFactory) LogMockUsage(componentName string) {
	f.logger.Warn().
		Str("component", componentName).
		Str("environment", f.config.Environment).
		Msg("USING MOCK IMPLEMENTATION - NOT FOR PRODUCTION USE")
}

// GetComponent retrieves a component from the cache or creates it if it doesn't exist
func (f *AppFactory) GetComponent(key string, creator func() interface{}) interface{} {
	if component, ok := f.components[key]; ok {
		return component
	}

	component := creator()
	f.components[key] = component
	return component
}

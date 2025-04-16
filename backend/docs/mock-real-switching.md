# Centralized Mock/Real Implementation Switching

## Overview

This document outlines the standardized approach for switching between mock and real implementations in the application. The goal is to establish a consistent pattern for component substitution while ensuring that testing configurations cannot accidentally be enabled in production environments.

## Current Issues

1. **Inconsistent Mocking Approaches**:
   - Some components use interface substitution
   - Others use environment variables
   - Some use conditional compilation
   - No standardized approach for mock detection

2. **Production Safety Concerns**:
   - Lack of safeguards to prevent test/mock implementations in production
   - No centralized logging of mock usage
   - No documentation of which components support mocking

3. **Developer Experience**:
   - Unclear how to enable mocks for testing
   - Different components require different approaches

## Standardized Approach

### 1. Configuration-Based Component Selection

We will standardize on a configuration-based approach for component selection:

```go
// internal/config/config.go

type MockConfig struct {
    // Global mock setting (overrides all other settings if false)
    Enabled bool `mapstructure:"enabled" env:"MOCK_ENABLED"`
    
    // Component-specific mock settings
    MEXCClient     bool `mapstructure:"mexc_client" env:"MOCK_MEXC_CLIENT"`
    WalletService  bool `mapstructure:"wallet_service" env:"MOCK_WALLET_SERVICE"`
    AIService      bool `mapstructure:"ai_service" env:"MOCK_AI_SERVICE"`
    AuthMiddleware bool `mapstructure:"auth_middleware" env:"MOCK_AUTH_MIDDLEWARE"`
    // Add other components as needed
}

// Add to main Config struct
type Config struct {
    // Existing fields
    // ...
    
    // Mock configuration
    Mock MockConfig `mapstructure:"mock"`
    
    // Environment indicator
    Environment string `mapstructure:"environment" env:"APP_ENV" default:"development"`
}
```

### 2. Centralized Factory Logic

The unified `AppFactory` (from the refactoring plan) will include standardized mock detection and safety checks:

```go
// internal/factory/app_factory.go

// isMockAllowed checks if mock implementations are allowed
func (f *AppFactory) isMockAllowed() bool {
    // Never allow mocks in production
    if f.config.Environment == "production" {
        return false
    }
    return f.config.Mock.Enabled
}

// shouldUseMock checks if a specific component should use a mock implementation
func (f *AppFactory) shouldUseMock(componentKey string) bool {
    // First check if mocks are globally allowed
    if !f.isMockAllowed() {
        return false
    }
    
    // Then check component-specific setting
    switch componentKey {
    case "mexc_client":
        return f.config.Mock.MEXCClient
    case "wallet_service":
        return f.config.Mock.WalletService
    case "ai_service":
        return f.config.Mock.AIService
    case "auth_middleware":
        return f.config.Mock.AuthMiddleware
    default:
        // If we don't have a specific setting, default to false
        return false
    }
}

// logMockUsage logs when a mock implementation is being used
func (f *AppFactory) logMockUsage(componentName string) {
    f.logger.Warn().
        Str("component", componentName).
        Msg("USING MOCK IMPLEMENTATION - NOT FOR PRODUCTION USE")
}
```

### 3. Example Implementation

The factory methods will use the standardized mock detection:

```go
// GetMEXCClient returns either a real or mock MEXC client
func (f *AppFactory) GetMEXCClient() port.MEXCClient {
    if f.mexcClient != nil {
        return f.mexcClient
    }
    
    if f.shouldUseMock("mexc_client") {
        f.logMockUsage("MEXCClient")
        f.mexcClient = mock.NewMEXCClient(f.logger)
    } else {
        f.mexcClient = mexc.NewClient(
            f.config.MEXC.APIKey,
            f.config.MEXC.APISecret,
            f.logger,
        )
    }
    
    return f.mexcClient
}

// GetAIService returns either a real or mock AI service
func (f *AppFactory) GetAIService() port.AIService {
    serviceKey := "ai_service"
    
    if service, ok := f.services[serviceKey]; ok {
        return service.(port.AIService)
    }
    
    var service port.AIService
    var err error
    
    if f.shouldUseMock("ai_service") {
        f.logMockUsage("AIService")
        service = ai.NewMockGeminiAIService()
    } else {
        service, err = ai.NewGeminiAIService(f.config, f.logger)
        if err != nil {
            f.logger.Error().Err(err).Msg("Failed to create AI service, falling back to mock")
            service = ai.NewMockGeminiAIService()
            f.logMockUsage("AIService (fallback)")
        }
    }
    
    f.services[serviceKey] = service
    return service
}
```

## Configuration Examples

### .env Files

**.env.development**
```
# Enable mocks globally
MOCK_ENABLED=true

# Enable specific mocks
MOCK_MEXC_CLIENT=true
MOCK_AI_SERVICE=true

# Disable specific mocks
MOCK_WALLET_SERVICE=false
```

**.env.test**
```
# Enable all mocks for testing
MOCK_ENABLED=true
```

**.env.production**
```
# Disable all mocks for production (redundant as the code prevents them anyway)
MOCK_ENABLED=false
```

## Mock Implementation Guidelines

1. **Naming Convention**
   - All mock implementations should have a name starting with "Mock" (e.g., `MockMEXCClient`)
   - Factory methods for mocks should start with "NewMock" (e.g., `NewMockMEXCClient`)

2. **Behavior Requirements**
   - Mocks should return predictable, deterministic results
   - They should log operations at DEBUG level for traceability
   - Complex operations should be configurable (success/failure scenarios)

3. **Test Support**
   - Mocks should provide methods to inspect calls and set behaviors
   - Consider using testify/mock for more complex mocks

## Testing with Mocks

### Unit Tests
For unit tests, use direct mock instantiation rather than going through the factory:

```go
func TestWalletService(t *testing.T) {
    // Create mock dependencies
    mockMEXCClient := mock.NewMEXCClient(zerolog.Nop())
    mockWalletRepo := mock.NewWalletRepository(zerolog.Nop())
    
    // Create service with mocks
    walletService := usecase.NewWalletService(mockWalletRepo, mockMEXCClient, zerolog.Nop())
    
    // Test service
    // ...
}
```

### Integration Tests
For integration tests, configure the factory to use mocks:

```go
func TestIntegration(t *testing.T) {
    // Create test config with mocks enabled
    cfg := &config.Config{
        Environment: "test",
        Mock: config.MockConfig{
            Enabled:     true,
            MEXCClient:  true,
            AIService:   true,
        },
    }
    
    // Create factory with test config
    factory := factory.NewAppFactory(cfg, zerolog.Nop(), testDB)
    
    // Use factory to get components with mocks
    mexcClient := factory.GetMEXCClient() // Will be a mock
    aiService := factory.GetAIService()   // Will be a mock
    
    // Test with mocked components
    // ...
}
```

## Conclusion

By centralizing the mock/real switching logic and enforcing a consistent pattern, we will:

1. Improve code maintainability through standardization
2. Prevent accidental use of test implementations in production
3. Make the developer experience more consistent and predictable
4. Provide clear visibility into which components are using mock implementations

This approach balances flexibility for testing with safety for production deployments. 
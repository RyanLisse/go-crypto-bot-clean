# Authentication Middleware

This package provides a standardized authentication middleware for the application. It supports multiple authentication providers, including Clerk and a test provider for development and testing.

## Components

### Interfaces

- `AuthMiddleware`: Interface for authentication middleware
  - `Middleware()`: Returns a middleware function that validates authentication
  - `RequireAuthentication()`: Middleware that requires authentication
  - `RequireRole()`: Middleware that requires a specific role

### Implementations

- `ClerkMiddleware`: Authentication middleware using Clerk
- `TestMiddleware`: Authentication middleware for testing
- `DisabledMiddleware`: Authentication middleware that bypasses authentication

### Factory

- `AuthFactory`: Factory for creating authentication middleware
  - `CreateMiddleware()`: Creates an authentication middleware based on the configuration
  - `CreateDefaultMiddleware()`: Creates the default authentication middleware based on the environment

### Error Handling

- `UnifiedErrorMiddleware`: Middleware that combines error handling, recovery, logging, and tracing
  - `Middleware()`: Returns a middleware function that handles errors, recovers from panics, and logs requests

## Usage

### Basic Usage

```go
// Create a logger
logger := zerolog.New(os.Stdout)

// Create a config
cfg := &config.Config{
    Auth: config.AuthConfig{
        Provider:       "clerk",
        Disabled:       false,
        ClerkSecretKey: "your-clerk-secret-key",
    },
}

// Create an auth service
authService := service.NewClerkAuthService(cfg, &logger)

// Create an auth factory
factory := middleware.NewAuthFactory(authService, cfg, &logger)

// Create the default middleware
authMiddleware := factory.CreateDefaultMiddleware()

// Use the middleware in your router
router := chi.NewRouter()
router.Use(authMiddleware.Middleware())

// Protect a route
router.Group(func(r chi.Router) {
    r.Use(authMiddleware.RequireAuthentication)
    r.Get("/protected", protectedHandler)
})

// Require a specific role
router.Group(func(r chi.Router) {
    r.Use(authMiddleware.RequireRole("admin"))
    r.Get("/admin", adminHandler)
})
```

### Using with Dependency Injection

The middleware is designed to be used with the application's dependency injection container. See the `internal/di/providers_middleware.go` file for details.

## Configuration

The middleware can be configured using the `config.AuthConfig` struct:

```go
type AuthConfig struct {
    Provider          string // "clerk", "jwt", etc.
    Disabled          bool   // Whether authentication is disabled
    ClerkSecretKey    string // Clerk secret key
    ClerkJWTPublicKey string // Clerk JWT public key
    ClerkJWTTemplate  string // Clerk JWT template
    JWTSecret         string // JWT secret for testing
    UseEnhanced       bool   // Whether to use enhanced authentication
}
```

## Environment Variables

The middleware can be configured using the following environment variables:

- `AUTH_PROVIDER`: Authentication provider (default: "clerk")
- `DISABLE_AUTH`: Whether authentication is disabled (default: false)
- `CLERK_SECRET_KEY`: Clerk secret key
- `CLERK_JWT_PUBLIC_KEY`: Clerk JWT public key
- `CLERK_JWT_TEMPLATE`: Clerk JWT template
- `JWT_SECRET`: JWT secret for testing
- `USE_ENHANCED_AUTH`: Whether to use enhanced authentication (default: true)

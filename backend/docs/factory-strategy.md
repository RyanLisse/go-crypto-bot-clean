# Factory Standardization Strategy

This document outlines our approach to standardizing and consolidating factory implementations across the application, focusing on eliminating redundancy and ensuring consistent initialization of components.

## Core Factory Pattern

The application will use a consolidated factory approach based on dependency injection:

```go
// Example factory implementation
type ServiceFactory struct {
    // Dependencies
    config     *config.Config
    logger     *zap.Logger
    db         *gorm.DB
    repositories *RepositoryFactory
    // ... other dependencies
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(
    config *config.Config,
    logger *zap.Logger,
    db *gorm.DB,
    repositories *RepositoryFactory,
) *ServiceFactory {
    return &ServiceFactory{
        config:     config,
        logger:     logger,
        db:         db,
        repositories: repositories,
    }
}
```

## Factory Hierarchy

We will establish a clear hierarchy of factory components to avoid circular dependencies:

1. **ConfigFactory** - Creates configuration objects
2. **InfrastructureFactory** - Creates database connections, cache clients, etc.
3. **RepositoryFactory** - Creates repository implementations
4. **ServiceFactory** - Creates service implementations
5. **MiddlewareFactory** - Creates HTTP middleware
6. **HandlerFactory** - Creates HTTP handlers
7. **RouterFactory** - Creates HTTP routers

Each factory should only depend on factories that are earlier in the hierarchy.

## Consolidated Factory Components

### 1. Configuration Factory

```go
// internal/factory/config_factory.go

type ConfigFactory struct {
    logger *zap.Logger
}

func NewConfigFactory(logger *zap.Logger) *ConfigFactory {
    return &ConfigFactory{
        logger: logger,
    }
}

// CreateAppConfig creates the application configuration
func (f *ConfigFactory) CreateAppConfig() (*config.Config, error) {
    // Load configuration from environment variables and/or config files
    cfg, err := config.LoadConfig()
    if err != nil {
        f.logger.Error("Failed to load configuration", zap.Error(err))
        return nil, err
    }
    
    // Validate configuration
    if err := cfg.Validate(); err != nil {
        f.logger.Error("Invalid configuration", zap.Error(err))
        return nil, err
    }
    
    return cfg, nil
}
```

### 2. Infrastructure Factory

```go
// internal/factory/infrastructure_factory.go

type InfrastructureFactory struct {
    config *config.Config
    logger *zap.Logger
}

func NewInfrastructureFactory(
    config *config.Config,
    logger *zap.Logger,
) *InfrastructureFactory {
    return &InfrastructureFactory{
        config: config,
        logger: logger,
    }
}

// CreateDBConnection creates a database connection
func (f *InfrastructureFactory) CreateDBConnection() (*gorm.DB, error) {
    // Create database connection based on configuration
    db, err := persistence.NewGormDB(
        f.config.Database.Driver,
        f.config.Database.DSN,
        f.config.Database.MaxOpenConns,
        f.config.Database.MaxIdleConns,
        f.config.Database.ConnMaxLifetime,
    )
    if err != nil {
        f.logger.Error("Failed to create database connection", zap.Error(err))
        return nil, err
    }
    
    return db, nil
}

// CreateRedisClient creates a Redis client
func (f *InfrastructureFactory) CreateRedisClient() (*redis.Client, error) {
    // Create Redis client based on configuration
    redisOpts := &redis.Options{
        Addr:     f.config.Redis.Address,
        Password: f.config.Redis.Password,
        DB:       f.config.Redis.DB,
    }
    
    client := redis.NewClient(redisOpts)
    
    // Test the connection
    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        f.logger.Error("Failed to connect to Redis", zap.Error(err))
        return nil, err
    }
    
    return client, nil
}

// CreateHTTPClient creates an HTTP client
func (f *InfrastructureFactory) CreateHTTPClient() *http.Client {
    // Create HTTP client with configuration
    return &http.Client{
        Timeout: f.config.HTTP.Timeout,
        Transport: &http.Transport{
            MaxIdleConns:        f.config.HTTP.MaxIdleConns,
            MaxIdleConnsPerHost: f.config.HTTP.MaxIdleConnsPerHost,
            IdleConnTimeout:     f.config.HTTP.IdleConnTimeout,
        },
    }
}

// CreateRateLimitStore creates a rate limit store
func (f *InfrastructureFactory) CreateRateLimitStore() (port.RateLimitStore, error) {
    // Choose the appropriate implementation based on configuration
    switch f.config.RateLimit.StoreType {
    case "memory":
        return cache.NewInMemoryRateLimitStore(), nil
    case "redis":
        redisClient, err := f.CreateRedisClient()
        if err != nil {
            return nil, err
        }
        return cache.NewRedisRateLimitStore(redisClient), nil
    default:
        return nil, fmt.Errorf("unsupported rate limit store type: %s", f.config.RateLimit.StoreType)
    }
}
```

### 3. Repository Factory

```go
// internal/factory/repository_factory.go

type RepositoryFactory struct {
    config *config.Config
    logger *zap.Logger
    db     *gorm.DB
}

func NewRepositoryFactory(
    config *config.Config,
    logger *zap.Logger,
    db *gorm.DB,
) *RepositoryFactory {
    return &RepositoryFactory{
        config: config,
        logger: logger,
        db:     db,
    }
}

// CreateUserRepository creates a user repository
func (f *RepositoryFactory) CreateUserRepository() port.UserRepository {
    return gorm.NewUserRepository(f.db)
}

// CreateMarketRepository creates a market repository
func (f *RepositoryFactory) CreateMarketRepository() port.MarketRepository {
    return gorm.NewMarketRepository(f.db)
}

// CreateOrderRepository creates an order repository
func (f *RepositoryFactory) CreateOrderRepository() port.OrderRepository {
    return gorm.NewOrderRepository(f.db)
}

// CreateWalletRepository creates a wallet repository
func (f *RepositoryFactory) CreateWalletRepository() port.WalletRepository {
    return gorm.NewWalletRepository(f.db)
}

// CreateAutoBuyRepository creates an auto buy repository
func (f *RepositoryFactory) CreateAutoBuyRepository() port.AutoBuyRepository {
    return gorm.NewAutoBuyRepository(f.db)
}

// CreateAPIKeyRepository creates an API key repository
func (f *RepositoryFactory) CreateAPIKeyRepository() port.APIKeyRepository {
    return gorm.NewAPIKeyRepository(f.db)
}
```

### 4. Service Factory

```go
// internal/factory/service_factory.go

type ServiceFactory struct {
    config       *config.Config
    logger       *zap.Logger
    repositories *RepositoryFactory
    httpClient   *http.Client
}

func NewServiceFactory(
    config *config.Config,
    logger *zap.Logger,
    repositories *RepositoryFactory,
    httpClient *http.Client,
) *ServiceFactory {
    return &ServiceFactory{
        config:       config,
        logger:       logger,
        repositories: repositories,
        httpClient:   httpClient,
    }
}

// CreateUserService creates a user service
func (f *ServiceFactory) CreateUserService() port.UserService {
    userRepo := f.repositories.CreateUserRepository()
    return service.NewUserService(userRepo, f.logger)
}

// CreateAuthService creates an authentication service
func (f *ServiceFactory) CreateAuthService() port.AuthService {
    // Choose the appropriate implementation based on configuration
    switch f.config.Auth.Provider {
    case "clerk":
        return service.NewClerkAuthService(
            f.config.Auth.ClerkPublishableKey,
            f.config.Auth.ClerkSecretKey,
            f.httpClient,
            f.logger,
        )
    case "simple":
        userRepo := f.repositories.CreateUserRepository()
        return service.NewSimpleAuthService(userRepo, f.logger)
    case "mock":
        // Only allow mock in non-production environments
        if f.config.Environment == "production" {
            f.logger.Fatal("Attempted to use mock auth service in production")
        }
        return service.NewMockAuthService(f.logger)
    default:
        f.logger.Fatal("Unsupported auth provider", zap.String("provider", f.config.Auth.Provider))
        return nil // This will never be reached due to Fatal, but satisfies the compiler
    }
}

// CreateMarketService creates a market service
func (f *ServiceFactory) CreateMarketService() port.MarketService {
    marketRepo := f.repositories.CreateMarketRepository()
    return service.NewMarketService(marketRepo, f.logger)
}

// CreateOrderService creates an order service
func (f *ServiceFactory) CreateOrderService() port.OrderService {
    orderRepo := f.repositories.CreateOrderRepository()
    userRepo := f.repositories.CreateUserRepository()
    return service.NewOrderService(orderRepo, userRepo, f.logger)
}

// CreateWalletService creates a wallet service
func (f *ServiceFactory) CreateWalletService() port.WalletService {
    walletRepo := f.repositories.CreateWalletRepository()
    return service.NewWalletService(walletRepo, f.logger)
}

// CreateExchangeService creates an exchange service
func (f *ServiceFactory) CreateExchangeService() port.ExchangeService {
    switch f.config.Exchange.Provider {
    case "mexc":
        return exchange.NewMEXCExchangeService(
            f.config.Exchange.MEXC.APIKey,
            f.config.Exchange.MEXC.APISecret,
            f.httpClient,
            f.logger,
        )
    case "mock":
        // Only allow mock in non-production environments
        if f.config.Environment == "production" {
            f.logger.Fatal("Attempted to use mock exchange service in production")
        }
        return exchange.NewMockExchangeService(f.logger)
    default:
        f.logger.Fatal("Unsupported exchange provider", zap.String("provider", f.config.Exchange.Provider))
        return nil // This will never be reached due to Fatal, but satisfies the compiler
    }
}

// CreateAutoBuyService creates an auto buy service
func (f *ServiceFactory) CreateAutoBuyService() port.AutoBuyService {
    autoBuyRepo := f.repositories.CreateAutoBuyRepository()
    orderService := f.CreateOrderService()
    marketService := f.CreateMarketService()
    exchangeService := f.CreateExchangeService()
    
    return service.NewAutoBuyService(
        autoBuyRepo,
        orderService,
        marketService,
        exchangeService,
        f.logger,
    )
}
```

### 5. Handler Factory

```go
// internal/factory/handler_factory.go

type HandlerFactory struct {
    config   *config.Config
    logger   *zap.Logger
    services *ServiceFactory
}

func NewHandlerFactory(
    config *config.Config,
    logger *zap.Logger,
    services *ServiceFactory,
) *HandlerFactory {
    return &HandlerFactory{
        config:   config,
        logger:   logger,
        services: services,
    }
}

// CreateUserHandler creates a user handler
func (f *HandlerFactory) CreateUserHandler() *handler.UserHandler {
    userService := f.services.CreateUserService()
    return handler.NewUserHandler(userService, f.logger)
}

// CreateAuthHandler creates an authentication handler
func (f *HandlerFactory) CreateAuthHandler() *handler.AuthHandler {
    authService := f.services.CreateAuthService()
    userService := f.services.CreateUserService()
    return handler.NewAuthHandler(authService, userService, f.logger)
}

// CreateMarketHandler creates a market handler
func (f *HandlerFactory) CreateMarketHandler() *handler.MarketHandler {
    marketService := f.services.CreateMarketService()
    return handler.NewMarketHandler(marketService, f.logger)
}

// CreateOrderHandler creates an order handler
func (f *HandlerFactory) CreateOrderHandler() *handler.OrderHandler {
    orderService := f.services.CreateOrderService()
    return handler.NewOrderHandler(orderService, f.logger)
}

// CreateWalletHandler creates a wallet handler
func (f *HandlerFactory) CreateWalletHandler() *handler.WalletHandler {
    walletService := f.services.CreateWalletService()
    return handler.NewWalletHandler(walletService, f.logger)
}

// CreateAutoBuyHandler creates an auto buy handler
func (f *HandlerFactory) CreateAutoBuyHandler() *handler.AutoBuyHandler {
    autoBuyService := f.services.CreateAutoBuyService()
    return handler.NewAutoBuyHandler(autoBuyService, f.logger)
}
```

### 6. Middleware Factory

```go
// internal/factory/middleware_factory.go

type MiddlewareFactory struct {
    config        *config.Config
    logger        *zap.Logger
    services      *ServiceFactory
    rateLimitStore port.RateLimitStore
}

func NewMiddlewareFactory(
    config *config.Config,
    logger *zap.Logger,
    services *ServiceFactory,
    rateLimitStore port.RateLimitStore,
) *MiddlewareFactory {
    return &MiddlewareFactory{
        config:        config,
        logger:        logger,
        services:      services,
        rateLimitStore: rateLimitStore,
    }
}

// CreateAuthMiddleware creates an authentication middleware
func (f *MiddlewareFactory) CreateAuthMiddleware() *middleware.ConsolidatedAuthMiddleware {
    authService := f.services.CreateAuthService()
    userService := f.services.CreateUserService()
    
    // Test mode is only allowed in non-production environments
    enableDummyAuth := f.config.Auth.EnableDummyAuth && f.config.Environment != "production"
    
    return middleware.NewConsolidatedAuthMiddleware(
        authService,
        userService,
        f.logger,
        enableDummyAuth,
        f.config.Auth.DummyUserID,
    )
}

// CreateRateLimiterMiddleware creates a rate limiter middleware
func (f *MiddlewareFactory) CreateRateLimiterMiddleware() *middleware.RateLimiterMiddleware {
    return middleware.NewRateLimiterMiddleware(
        f.rateLimitStore,
        f.logger,
    )
}

// Additional middleware factory methods as defined in the middleware strategy document
```

### 7. Router Factory

```go
// internal/factory/router_factory.go

type RouterFactory struct {
    config      *config.Config
    logger      *zap.Logger
    handlers    *HandlerFactory
    middlewares *MiddlewareFactory
}

func NewRouterFactory(
    config *config.Config,
    logger *zap.Logger,
    handlers *HandlerFactory,
    middlewares *MiddlewareFactory,
) *RouterFactory {
    return &RouterFactory{
        config:      config,
        logger:      logger,
        handlers:    handlers,
        middlewares: middlewares,
    }
}

// CreateRouter creates the main application router
func (f *RouterFactory) CreateRouter() *chi.Mux {
    router := chi.NewRouter()
    
    // Apply global middleware
    router.Use(middleware.RequestID)
    router.Use(f.middlewares.CreateErrorMiddleware().Middleware())
    router.Use(f.middlewares.CreateLoggingMiddleware().Middleware())
    router.Use(f.middlewares.CreateCORSMiddleware().Middleware())
    router.Use(f.middlewares.CreateRateLimiterMiddleware().Middleware())
    router.Use(f.middlewares.CreateAuthMiddleware().Middleware())
    router.Use(middleware.Recoverer)
    
    // Mount API routes
    router.Route("/api", func(r chi.Router) {
        // User routes
        r.Route("/users", f.CreateUserRoutes)
        
        // Auth routes
        r.Route("/auth", f.CreateAuthRoutes)
        
        // Market routes
        r.Route("/markets", f.CreateMarketRoutes)
        
        // Order routes
        r.Route("/orders", f.CreateOrderRoutes)
        
        // Wallet routes
        r.Route("/wallets", f.CreateWalletRoutes)
        
        // Auto buy routes
        r.Route("/auto-buys", f.CreateAutoBuyRoutes)
    })
    
    return router
}

// CreateUserRoutes creates user routes
func (f *RouterFactory) CreateUserRoutes(r chi.Router) {
    userHandler := f.handlers.CreateUserHandler()
    authMiddleware := f.middlewares.CreateAuthMiddleware()
    
    r.Get("/", userHandler.List)
    r.Post("/", userHandler.Create)
    
    r.Route("/{id}", func(r chi.Router) {
        r.Use(authMiddleware.RequireAuthentication)
        r.Get("/", userHandler.Get)
        r.Put("/", userHandler.Update)
        r.Delete("/", userHandler.Delete)
    })
}

// Additional route creation methods for other resources
```

## App Factory

To tie everything together, we'll create an app factory:

```go
// internal/factory/app_factory.go

type AppFactory struct {
    logger *zap.Logger
}

func NewAppFactory(logger *zap.Logger) *AppFactory {
    return &AppFactory{
        logger: logger,
    }
}

// CreateApp creates the complete application
func (f *AppFactory) CreateApp() (*app.App, error) {
    // Create configuration
    configFactory := NewConfigFactory(f.logger)
    config, err := configFactory.CreateAppConfig()
    if err != nil {
        return nil, err
    }
    
    // Create infrastructure
    infraFactory := NewInfrastructureFactory(config, f.logger)
    
    db, err := infraFactory.CreateDBConnection()
    if err != nil {
        return nil, err
    }
    
    httpClient := infraFactory.CreateHTTPClient()
    
    rateLimitStore, err := infraFactory.CreateRateLimitStore()
    if err != nil {
        return nil, err
    }
    
    // Create repositories
    repoFactory := NewRepositoryFactory(config, f.logger, db)
    
    // Create services
    serviceFactory := NewServiceFactory(config, f.logger, repoFactory, httpClient)
    
    // Create handlers
    handlerFactory := NewHandlerFactory(config, f.logger, serviceFactory)
    
    // Create middlewares
    middlewareFactory := NewMiddlewareFactory(config, f.logger, serviceFactory, rateLimitStore)
    
    // Create router
    routerFactory := NewRouterFactory(config, f.logger, handlerFactory, middlewareFactory)
    router := routerFactory.CreateRouter()
    
    // Create app
    app := app.NewApp(
        config,
        f.logger,
        router,
        db,
    )
    
    return app, nil
}
```

## Testing Support

For testing, we'll create special factory methods to create mock implementations:

```go
// internal/factory/test_factory.go

// CreateTestUserRepository creates a test user repository
func (f *RepositoryFactory) CreateTestUserRepository() port.UserRepository {
    // Only allow test repositories in non-production environments
    if f.config.Environment == "production" {
        f.logger.Fatal("Attempted to create test repository in production")
    }
    return mock.NewMockUserRepository()
}

// CreateTestAuthService creates a test auth service
func (f *ServiceFactory) CreateTestAuthService() port.AuthService {
    // Only allow test services in non-production environments
    if f.config.Environment == "production" {
        f.logger.Fatal("Attempted to create test service in production")
    }
    return mock.NewMockAuthService(f.logger)
}

// Additional test factory methods for other components
```

## Safety Mechanisms

To prevent test/mock implementations from being used in production:

1. All factory methods that create test/mock implementations should check the environment:

```go
if f.config.Environment == "production" {
    f.logger.Fatal("Attempted to create test/mock implementation in production")
}
```

2. Configuration validation should enforce safety constraints:

```go
func (c *Config) Validate() error {
    // Ensure no test/mock implementations in production
    if c.Environment == "production" {
        if c.Auth.Provider == "mock" {
            return errors.New("mock auth provider not allowed in production")
        }
        if c.Exchange.Provider == "mock" {
            return errors.New("mock exchange provider not allowed in production")
        }
        if c.Auth.EnableDummyAuth {
            return errors.New("dummy auth not allowed in production")
        }
    }
    
    return nil
}
```

## Implementation Selection Mechanism

To support switching between real and mock implementations:

```go
// internal/config/config.go

type Config struct {
    // Other fields...
    
    // Environment type (development, testing, production)
    Environment string `yaml:"environment" env:"ENVIRONMENT"`
    
    // Authentication configuration
    Auth struct {
        Provider           string `yaml:"provider" env:"AUTH_PROVIDER"`
        ClerkPublishableKey string `yaml:"clerk_publishable_key" env:"CLERK_PUBLISHABLE_KEY"`
        ClerkSecretKey     string `yaml:"clerk_secret_key" env:"CLERK_SECRET_KEY"`
        EnableDummyAuth    bool   `yaml:"enable_dummy_auth" env:"ENABLE_DUMMY_AUTH"`
        DummyUserID        string `yaml:"dummy_user_id" env:"DUMMY_USER_ID"`
    } `yaml:"auth"`
    
    // Exchange configuration
    Exchange struct {
        Provider string `yaml:"provider" env:"EXCHANGE_PROVIDER"`
        MEXC     struct {
            APIKey    string `yaml:"api_key" env:"MEXC_API_KEY"`
            APISecret string `yaml:"api_secret" env:"MEXC_API_SECRET"`
        } `yaml:"mexc"`
    } `yaml:"exchange"`
    
    // Database configuration
    Database struct {
        Driver          string        `yaml:"driver" env:"DB_DRIVER"`
        DSN             string        `yaml:"dsn" env:"DB_DSN"`
        MaxOpenConns    int           `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
        MaxIdleConns    int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
        ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME"`
        UseMock         bool          `yaml:"use_mock" env:"DB_USE_MOCK"`
    } `yaml:"database"`
    
    // Other configuration sections...
}
```

## Migration Plan

To migrate from multiple factory implementations to the consolidated approach:

1. Create the new consolidated factory components
2. Implement the hierarchy of factories
3. Update application initialization to use the new factory structure
4. Replace direct component instantiation with factory method calls
5. Remove deprecated factory implementations

## Conclusion

This standardized factory approach provides several benefits:

1. **Clarity**: Clear separation of concerns with a defined hierarchy
2. **Consistency**: Uniform approach to component creation
3. **Maintainability**: Centralized configuration and dependency management
4. **Security**: Prevention of test/mock implementations in production
5. **Testability**: Easy substitution of real implementations with mocks
6. **Flexibility**: Environment-specific configuration and implementation selection 
package api

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"go-crypto-bot-clean/backend/internal/api/controllers"
	"go-crypto-bot-clean/backend/internal/api/handlers"
	apirepository "go-crypto-bot-clean/backend/internal/api/repository"
	"go-crypto-bot-clean/backend/internal/api/service"
	"go-crypto-bot-clean/backend/internal/api/websocket"
	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/config"
	"go-crypto-bot-clean/backend/internal/core/account"
	aiservice "go-crypto-bot-clean/backend/internal/domain/ai/service"
	"go-crypto-bot-clean/backend/internal/domain/repositories"
	"go-crypto-bot-clean/backend/internal/domain/repository"
	domainservice "go-crypto-bot-clean/backend/internal/domain/service"
	"go-crypto-bot-clean/backend/internal/logging"
	"go-crypto-bot-clean/backend/internal/platform/database/gorm"
	gormrepo "go-crypto-bot-clean/backend/internal/platform/database/gorm/repositories"
	"go-crypto-bot-clean/backend/internal/platform/mexc"
)

// Dependencies contains all the dependencies for the API.
type Dependencies struct {
	// router field is used in SetupConsolidatedRouter
	router                 chi.Router
	logger                 *zap.Logger
	mexcClient             domainservice.ExchangeService
	HealthHandler          *handlers.HealthHandler
	StatusHandler          *handlers.StatusHandler
	PortfolioHandler       *handlers.PortfolioHandler
	TradeHandler           *handlers.TradeHandler
	NewCoinHandler         *handlers.NewCoinsHandler
	CoinHandler            *handlers.CoinHandler
	ConfigHandler          *handlers.ConfigHandler
	WebSocketHandler       *websocket.Handler
	AuthHandler            *handlers.AuthHandler
	AnalyticsHandler       *handlers.AnalyticsHandler
	EnhancedAccountHandler *handlers.EnhancedAccountHandler
	BacktestHandler        *handlers.BacktestHandler
	AccountController      *controllers.AccountController

	// AI Service
	AIService aiservice.AIService

	// Repositories
	BoughtCoinRepository  repositories.BoughtCoinRepository
	NewCoinRepository     repositories.NewCoinRepository
	UserRepository        apirepository.UserRepository
	WalletRepository      repository.WalletRepository
	TransactionRepository repository.TransactionRepository

	// Authentication
	ValidAPIKeys map[string]struct{}
	Config       *config.Config
	AuthService  auth.AuthProvider

	// Rate limiting
	RateLimit struct {
		Rate     float64
		Capacity int
	}

	// Services
	AccountService   account.AccountService
	UserService      *service.UserService
	StrategyService  *service.StrategyService
	NewCoinService   interface{} // Using interface{} to allow for both real and mock implementations
	AnalyticsService interface{} // Using interface{} to allow for both real and mock implementations
}

// NewDependencies creates a new Dependencies instance.
func NewDependencies(cfg *config.Config, logger *zap.Logger) (*Dependencies, error) {
	deps := &Dependencies{
		Config: cfg,
		logger: logger,
	}

	// Initialize MEXC client
	stdLogger := zap.NewStdLog(logger) // Create standard logger from Zap logger
	client, err := mexc.NewClient(cfg, mexc.WithLogger(stdLogger))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MEXC client: %w", err)
	}
	deps.mexcClient = client

	// Initialize database and repositories
	if err := deps.initializeDatabaseAndRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize database and repositories: %w", err)
	}

	var authService auth.AuthProvider
	if cfg.Auth.Enabled {
		// Initialize auth service with Clerk
		authService = auth.NewService(cfg.Auth.ClerkSecretKey)
	} else {
		// Use disabled auth service when auth is not enabled
		authService = auth.NewDisabledService()
	}

	deps.AuthService = authService

	// Initialize the real account service
	deps.logger.Info("Initializing real account service with MEXC API")
	if err := deps.initializeRealAccountService(); err != nil {
		// Log detailed error information
		deps.logger.Error("Failed to initialize real account service",
			zap.Error(err),
			zap.String("api_key_set", boolToString(deps.Config.Mexc.APIKey != "")),
			zap.String("secret_key_set", boolToString(deps.Config.Mexc.SecretKey != "")))

		// Instead of returning error, fall back to mock service when in development mode
		if deps.Config.App.Environment == "development" {
			deps.logger.Warn("Falling back to mock account service for development mode")
			// Use the mock account service already defined in the package
			deps.AccountService = &MockAccountService{}

			// Create account service adapter for the handler
			accountServiceAdapter := NewRealAccountServiceAdapter(deps.AccountService, deps.logger)

			// Initialize the EnhancedAccountHandler with the account service adapter
			deps.EnhancedAccountHandler = handlers.NewEnhancedAccountHandler(accountServiceAdapter)
		} else {
			// Only return error in production mode
			return nil, fmt.Errorf("failed to initialize real account service: %w", err)
		}
	}

	return deps, nil
}

// initializeDatabaseAndRepositories initializes the database connection and repositories
func (deps *Dependencies) initializeDatabaseAndRepositories() error {
	// Initialize GORM database
	gormDB, err := gorm.NewDatabase(gorm.Config{
		Path:            deps.Config.Database.Path,
		Debug:           deps.Config.App.Debug,
		Logger:          deps.logger,
		MaxIdleConns:    deps.Config.Database.MaxIdleConns,
		MaxOpenConns:    deps.Config.Database.MaxOpenConns,
		ConnMaxLifetime: 0, // Use default
	})
	if err != nil {
		deps.logger.Error("Failed to initialize GORM database", zap.Error(err))
		return err
	}

	// Run migrations
	if err := gorm.RunMigrations(gormDB, deps.logger); err != nil {
		deps.logger.Error("Failed to run database migrations", zap.Error(err))
		return err
	}
	// Create repositories
	boughtCoinRepo := gormrepo.NewGORMBoughtCoinRepository(gormDB, deps.logger)
	newCoinRepo := gormrepo.NewGORMNewCoinRepository(gormDB, deps.logger)

	// Get repository factory for transaction repository
	factory := gormrepo.NewRepositoryFactory(gormDB, &logging.LoggerWrapper{Logger: deps.logger})

	// Store repositories in dependencies
	deps.BoughtCoinRepository = boughtCoinRepo
	deps.NewCoinRepository = newCoinRepo

	// Create mock wallet repository since there's no actual implementation
	deps.WalletRepository = NewMockWalletRepository(deps.logger)

	// Get transaction repository from factory and create adapter
	transactionRepo := factory.GetTransactionRepository()
	deps.TransactionRepository = NewTransactionRepositoryAdapter(transactionRepo, deps.logger)

	deps.logger.Info("Database and repositories initialized successfully")
	return nil
}

// initializeRealAccountService initializes the real account service with the MEXC client and repositories
func (deps *Dependencies) initializeRealAccountService() error {
	// Validate configuration parameters first
	if err := deps.validateMEXCConfiguration(); err != nil {
		return fmt.Errorf("invalid MEXC configuration: %w", err)
	}

	// Type assert the mexcClient to get access to the internal clients
	mexcClientConcrete, ok := deps.mexcClient.(*mexc.Client)
	if !ok {
		deps.logger.Error("Failed to type assert mexcClient to *mexc.Client")
		return fmt.Errorf("mexcClient is not of type *mexc.Client")
	}

	// Log API key information (without revealing the actual keys)
	deps.logger.Debug("MEXC API configuration",
		zap.String("base_url", deps.Config.Mexc.BaseURL),
		zap.String("websocket_url", deps.Config.Mexc.WebsocketURL),
		zap.String("api_key_set", boolToString(deps.Config.Mexc.APIKey != "")),
		zap.String("secret_key_set", boolToString(deps.Config.Mexc.SecretKey != "")))

	// Get REST and WebSocket clients
	restClient := mexcClientConcrete.GetRestClient()
	wsClient := mexcClientConcrete.GetWsClient()

	// No need to check for nil as GetRestClient never returns nil
	// Just log the client info for debugging
	deps.logger.Debug("Using REST client", zap.String("client_type", fmt.Sprintf("%T", restClient)))

	// Validate API keys before proceeding
	deps.logger.Debug("Validating MEXC API keys before initializing services")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to validate API keys - forcing real API usage
	valid, err := restClient.ValidateKeys(ctx)
	if err != nil {
		// Log the detailed error
		deps.logger.Error("Failed to validate MEXC API keys",
			zap.Error(err),
			zap.String("api_key_length", fmt.Sprintf("%d", len(deps.Config.Mexc.APIKey))),
			zap.String("secret_key_length", fmt.Sprintf("%d", len(deps.Config.Mexc.SecretKey))))

		// Return error without falling back to mock
		return fmt.Errorf("failed to validate MEXC API keys: %w", err)
	}

	if !valid {
		deps.logger.Error("MEXC API keys are invalid")
		return fmt.Errorf("MEXC API keys are invalid")
	}

	deps.logger.Info("Successfully validated MEXC API keys")

	// Create adapter for BoughtCoinRepository
	boughtCoinAdapter := NewBoughtCoinRepositoryAdapter(deps.BoughtCoinRepository, deps.logger)

	// Create config adapter
	configAdapter := NewConfigAdapter(deps.Config)

	// Create the real account service with REST and WebSocket clients
	realAccountSvc := account.NewRealAccountServiceWithLogger(
		restClient,
		wsClient,
		boughtCoinAdapter,
		deps.WalletRepository,
		deps.TransactionRepository,
		configAdapter,
		deps.logger,
	)

	deps.AccountService = realAccountSvc

	// Create account service adapter for the handler
	accountServiceAdapter := NewRealAccountServiceAdapter(deps.AccountService, deps.logger)

	// Initialize the EnhancedAccountHandler with the account service adapter
	deps.EnhancedAccountHandler = handlers.NewEnhancedAccountHandler(accountServiceAdapter)

	deps.logger.Info("Real account service created successfully")

	return nil
}

// validateMEXCConfiguration validates the MEXC API configuration
func (deps *Dependencies) validateMEXCConfiguration() error {
	// Check if API keys are present
	if deps.Config.Mexc.APIKey == "" {
		// In development mode, log a warning but don't return an error
		if deps.Config.App.Environment == "development" {
			deps.logger.Warn("MEXC API key is missing, will use mock services in development mode")
			return nil
		}
		// In production, return an error
		deps.logger.Error("MEXC API key is missing")
		return fmt.Errorf("MEXC API key is missing or empty")
	}

	if deps.Config.Mexc.SecretKey == "" {
		// In development mode, log a warning but don't return an error
		if deps.Config.App.Environment == "development" {
			deps.logger.Warn("MEXC Secret key is missing, will use mock services in development mode")
			return nil
		}
		// In production, return an error
		deps.logger.Error("MEXC Secret key is missing")
		return fmt.Errorf("MEXC Secret key is missing or empty")
	}

	// Validate API key format (basic validation)
	if len(deps.Config.Mexc.APIKey) < 10 {
		deps.logger.Error("MEXC API key is too short",
			zap.Int("length", len(deps.Config.Mexc.APIKey)),
			zap.String("expected", "at least 10 characters"))
		return fmt.Errorf("MEXC API key is too short (length: %d)", len(deps.Config.Mexc.APIKey))
	}

	if len(deps.Config.Mexc.SecretKey) < 10 {
		deps.logger.Error("MEXC Secret key is too short",
			zap.Int("length", len(deps.Config.Mexc.SecretKey)),
			zap.String("expected", "at least 10 characters"))
		return fmt.Errorf("MEXC Secret key is too short (length: %d)", len(deps.Config.Mexc.SecretKey))
	}

	// Validate base URL
	if deps.Config.Mexc.BaseURL == "" {
		deps.logger.Warn("MEXC Base URL is empty, using default")
		deps.Config.Mexc.BaseURL = "https://api.mexc.com"
	}

	// Validate websocket URL
	if deps.Config.Mexc.WebsocketURL == "" {
		deps.logger.Warn("MEXC Websocket URL is empty, using default")
		deps.Config.Mexc.WebsocketURL = "wss://wbs.mexc.com/ws"
	}

	return nil
}

// initializeWithRealServices is a reference implementation for switching from mock to real services at runtime.
// It's currently not used but kept for future implementation of dynamic service switching.
func (deps *Dependencies) initializeWithRealServices() {
	deps.logger.Info("Initializing with real services")

	// Create account service adapter
	accountServiceAdapter := NewRealAccountServiceAdapter(deps.AccountService, deps.logger)

	deps.EnhancedAccountHandler = handlers.NewEnhancedAccountHandler(accountServiceAdapter)
}

// Initialize the status handler with a mock service
func (deps *Dependencies) InitializeStatusHandler() {
	deps.StatusHandler = handlers.NewStatusHandler(&MockStatusService{})
}

// Initialize the portfolio handler with a mock service
func (deps *Dependencies) InitializePortfolioHandler() {
	deps.PortfolioHandler = handlers.NewPortfolioHandler(&MockPortfolioService{})
}

// boolToString converts a boolean to a "yes" or "no" string
func boolToString(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

package api

import (
	"context"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ryanlisse/go-crypto-bot/internal/api/handlers"
	"github.com/ryanlisse/go-crypto-bot/internal/api/websocket"
	"github.com/ryanlisse/go-crypto-bot/internal/config"
	"github.com/ryanlisse/go-crypto-bot/internal/core/account"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/service"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repositories"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/mexc/rest"
	"go.uber.org/zap"
)

// Dependencies contains all the dependencies for the API.
type Dependencies struct {
	// Handlers
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

	// AI Service
	AIService service.AIService

	// Repositories
	BoughtCoinRepository repositories.BoughtCoinRepository
	NewCoinRepository    repositories.NewCoinRepository

	// Authentication
	ValidAPIKeys map[string]struct{}
	Config       *config.Config

	// Rate limiting
	RateLimit struct {
		Rate     float64
		Capacity int
	}
}

// NewDependencies creates a new Dependencies instance.
func NewDependencies(cfg *config.Config) *Dependencies {
	deps := &Dependencies{
		Config: cfg,
	}

	// Initialize API keys
	deps.ValidAPIKeys = make(map[string]struct{})
	for _, key := range cfg.Auth.APIKeys {
		deps.ValidAPIKeys[key] = struct{}{}
	}

	// Initialize rate limiting
	deps.RateLimit.Rate = 10     // Default: 10 requests per second
	deps.RateLimit.Capacity = 20 // Default: burst capacity of 20 requests

	// Initialize handlers
	deps.HealthHandler = handlers.NewHealthHandler()
	deps.AuthHandler = handlers.NewAuthHandler(
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTExpiry,
		cfg.Auth.CookieName,
	)

	// Initialize status handler with mock service
	deps.StatusHandler = handlers.NewStatusHandler(&MockStatusService{})

	// Initialize logger and MEXC client for all services
	logger, err := zap.NewProduction()
	if err != nil {
		// Fall back to basic logger if we can't create the production logger
		logger = zap.NewExample()
	}

	// Validate API keys before proceeding
	if cfg.Mexc.APIKey == "" || cfg.Mexc.SecretKey == "" {
		logger.Error("MEXC API keys are not configured properly. Using mock services.")
		err = errors.New("missing MEXC API keys")
	}

	// Create MEXC client with the configured API keys
	var mexcClient *rest.Client
	if err == nil {
		mexcClient, err = rest.NewClient(cfg.Mexc.APIKey, cfg.Mexc.SecretKey, rest.WithLogger(logger))
		if err != nil {
			// Log the error but continue with mock services
			logger.Error("Failed to create MEXC client, will fall back to mock services", zap.Error(err))
		} else {
			// Validate the API keys
			valid, validateErr := mexcClient.ValidateAPIKeys(context.Background())
			if validateErr != nil || !valid {
				logger.Error("MEXC API keys validation failed", zap.Error(validateErr))
				err = errors.New("invalid MEXC API keys")
				mexcClient = nil
			} else {
				logger.Info("MEXC API keys validated successfully")
			}
		}
	}

	// Initialize Database dependencies (must be first)
	deps.InitializeDatabaseDependencies()

	// Initialize portfolio handler
	if err != nil || mexcClient == nil {
		// Fall back to mock service if we can't create the MEXC client
		logger.Error("Using mock portfolio service due to MEXC client initialization failure")
		deps.PortfolioHandler = handlers.NewPortfolioHandler(&MockPortfolioService{})
	} else {
		// Create real portfolio service adapter
		logger.Info("Using real portfolio service with MEXC client")
		// We don't pass the repository here as it has a different interface
		portfolioAdapter := NewRealPortfolioServiceAdapter(mexcClient, nil, logger)
		deps.PortfolioHandler = handlers.NewPortfolioHandler(portfolioAdapter)
	}

	// Initialize enhanced account handler with real account service
	if err != nil || mexcClient == nil {
		// Fall back to mock service if we can't create the MEXC client
		logger.Error("Using mock account service due to MEXC client initialization failure")
		deps.EnhancedAccountHandler = handlers.NewEnhancedAccountHandler(&MockAccountService{})
	} else {
		// Create mock config for the account service
		logger.Info("Using real account service with MEXC client")
		mockConfig := &MockAccountConfig{}

		// Create real account service - pass nil for the repositories that don't match the interface
		accountService := account.NewRealAccountService(
			mexcClient,  // MexcRESTClient
			nil,         // MexcWebSocketClient - we don't have a compatible implementation
			nil,         // BoughtCoinRepository - interface mismatch
			nil,         // WalletRepository
			nil,         // TransactionRepository
			mockConfig,  // Config
		)
		// Create adapter to make it compatible with the AccountServiceInterface
		accountAdapter := NewRealAccountServiceAdapter(accountService)
		deps.EnhancedAccountHandler = handlers.NewEnhancedAccountHandler(accountAdapter)
	}

	// Initialize NewCoin dependencies with the MEXC client
	deps.InitializeNewCoinDependencies()

	// Initialize Analytics dependencies
	deps.InitializeAnalyticsDependencies()

	// Initialize Trade dependencies
	deps.InitializeTradeDependencies()

	// Initialize Config dependencies
	deps.InitializeConfigDependencies()

	// Initialize WebSocket dependencies
	deps.InitializeWebSocketDependencies()

	return deps
}

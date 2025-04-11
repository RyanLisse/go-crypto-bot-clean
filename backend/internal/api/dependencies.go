package api

import (
	"context"
	"errors"
	"fmt"

	"go-crypto-bot-clean/backend/internal/api/handlers"
	apirepository "go-crypto-bot-clean/backend/internal/api/repository" // API layer repos
	"go-crypto-bot-clean/backend/internal/api/service"
	"go-crypto-bot-clean/backend/internal/api/websocket"
	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/config"
	aiservice "go-crypto-bot-clean/backend/internal/domain/ai/service"
	"go-crypto-bot-clean/backend/internal/domain/repositories" // Domain repos
	"go-crypto-bot-clean/backend/internal/platform/mexc/rest"

	_ "github.com/mattn/go-sqlite3"
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
	BacktestHandler        *handlers.BacktestHandler

	// AI Service
	AIService aiservice.AIService

	// Repositories
	BoughtCoinRepository repositories.BoughtCoinRepository
	NewCoinRepository    repositories.NewCoinRepository
	UserRepository       apirepository.UserRepository // Use the one from internal/api/repository

	// Authentication
	ValidAPIKeys map[string]struct{}
	Config       *config.Config
	AuthService  auth.AuthProvider

	// Rate limiting
	RateLimit struct {
		Rate     float64
		Capacity int
	}

	// Services for Huma integration
	BacktestService *service.BacktestService
	StrategyService *service.StrategyService
	UserService     *service.UserService

	// Logger
	logger *zap.Logger
}

// NewDependencies creates a new Dependencies instance.
func NewDependencies(cfg *config.Config) (*Dependencies, error) {
	deps := &Dependencies{
		Config: cfg,
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	deps.logger = logger

	// Initialize auth service
	authService := auth.NewService(cfg.Auth.ClerkSecretKey)
	deps.AuthService = authService

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
	// Auth handler no longer needs arguments
	deps.AuthHandler = handlers.NewAuthHandler()

	// Initialize status handler
	deps.InitializeStatusHandler()

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
			valid, validateErr := mexcClient.ValidateKeys(context.Background())
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
	// deps.InitializeDatabaseDependencies() // Method doesn't exist, remove call
	// TODO: Initialize repositories here directly once DB connection is available

	// Initialize portfolio handler
	deps.InitializePortfolioHandler()

	// Initialize enhanced account handler
	deps.InitializeAccountHandler()

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

	// Initialize Backtest dependencies
	deps.InitializeBacktestDependencies()

	return deps, nil
}

// Initialize the status handler with a mock service
func (deps *Dependencies) InitializeStatusHandler() {
	deps.StatusHandler = handlers.NewStatusHandler(&MockStatusService{})
}

// Initialize the account handler with a mock service
func (deps *Dependencies) InitializeAccountHandler() {
	deps.EnhancedAccountHandler = handlers.NewEnhancedAccountHandler(&MockAccountService{})
}

// Initialize the portfolio handler with a mock service
func (deps *Dependencies) InitializePortfolioHandler() {
	deps.PortfolioHandler = handlers.NewPortfolioHandler(&MockPortfolioService{})
}

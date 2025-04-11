package api

import (
	// Added import for standard logger

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"go-crypto-bot-clean/backend/internal/api/handlers"
	apirepository "go-crypto-bot-clean/backend/internal/api/repository"
	"go-crypto-bot-clean/backend/internal/api/service"
	"go-crypto-bot-clean/backend/internal/api/websocket"
	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/config"
	"go-crypto-bot-clean/backend/internal/core/account"
	aiservice "go-crypto-bot-clean/backend/internal/domain/ai/service"
	"go-crypto-bot-clean/backend/internal/domain/repositories"
	domainservice "go-crypto-bot-clean/backend/internal/domain/service"
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

	// AI Service
	AIService aiservice.AIService

	// Repositories
	BoughtCoinRepository repositories.BoughtCoinRepository
	NewCoinRepository    repositories.NewCoinRepository
	UserRepository       apirepository.UserRepository

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
	AccountService  account.AccountService
	UserService     *service.UserService
	StrategyService *service.StrategyService
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
		return nil, err
	}
	deps.mexcClient = client

	var authService auth.AuthProvider
	if cfg.Auth.Enabled {
		// Initialize auth service with Clerk
		authService = auth.NewService(cfg.Auth.ClerkSecretKey)
	} else {
		// Use disabled auth service when auth is not enabled
		authService = auth.NewDisabledService()
	}

	// Create account service
	accountService := account.NewSimpleAccountService()

	deps.AuthService = authService
	deps.AccountService = accountService

	// Initialize the EnhancedAccountHandler with the account service adapter
	if accountService != nil {
		deps.EnhancedAccountHandler = handlers.NewEnhancedAccountHandler(NewRealAccountServiceAdapter(accountService))
	} else {
		// Use mock account service if real service is not available
		deps.logger.Warn("Account service is nil, using mock service")
		deps.initializeWithMockServices()
	}

	return deps, nil
}

// Initialize with mock services when real services are not available
func (deps *Dependencies) initializeWithMockServices() {
	deps.logger.Info("Initializing with mock services")
	mockAccountService := account.NewSimpleAccountService()
	deps.EnhancedAccountHandler = handlers.NewEnhancedAccountHandler(NewRealAccountServiceAdapter(mockAccountService))
}

// Initialize with real services when available
// This method is kept for future use when dynamic service switching is implemented
func (deps *Dependencies) initializeWithRealServices() {
	deps.logger.Info("Initializing with real services")
	deps.EnhancedAccountHandler = handlers.NewEnhancedAccountHandler(NewRealAccountServiceAdapter(deps.AccountService))
}

// Initialize the status handler with a mock service
func (deps *Dependencies) InitializeStatusHandler() {
	deps.StatusHandler = handlers.NewStatusHandler(&MockStatusService{})
}

// Initialize the portfolio handler with a mock service
func (deps *Dependencies) InitializePortfolioHandler() {
	deps.PortfolioHandler = handlers.NewPortfolioHandler(&MockPortfolioService{})
}

package factory

import (
	"fmt"
	"net/http"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	gormdb "gorm.io/gorm"
)

// ConsolidatedFactory provides a unified factory for creating all components
type ConsolidatedFactory struct {
	db         *gormdb.DB
	logger     *zerolog.Logger
	cfg        *config.Config
	mexcClient port.MEXCClient
	txManager  port.TransactionManager
}

// NewConsolidatedFactory creates a new ConsolidatedFactory
func NewConsolidatedFactory(db *gormdb.DB, logger *zerolog.Logger, cfg *config.Config) *ConsolidatedFactory {
	// Create MEXC client
	mexcClient := NewMEXCClient(cfg, logger)

	// Create transaction manager
	txManager := gorm.NewTransactionManager(db, logger)

	return &ConsolidatedFactory{
		db:         db,
		logger:     logger,
		cfg:        cfg,
		mexcClient: mexcClient,
		txManager:  txManager,
	}
}

// GetMEXCClient returns the MEXC client
func (f *ConsolidatedFactory) GetMEXCClient() port.MEXCClient {
	return f.mexcClient
}

// GetWalletRepository returns a wallet repository
func (f *ConsolidatedFactory) GetWalletRepository() port.WalletRepository {
	return repo.NewConsolidatedWalletRepository(f.db, f.logger)
}

// GetWalletService returns a wallet service
func (f *ConsolidatedFactory) GetWalletService() usecase.WalletService {
	walletRepo := f.GetWalletRepository()
	return usecase.NewWalletService(walletRepo, f.mexcClient, f.logger)
}

// GetOrderRepository returns an order repository
func (f *ConsolidatedFactory) GetOrderRepository() port.OrderRepository {
	// TODO: implement when needed
	return nil
}

// GetNewCoinRepository returns a new coin repository
func (f *ConsolidatedFactory) GetNewCoinRepository() port.NewCoinRepository {
	// TODO: implement when needed
	return nil
}

// GetEventRepository returns an event repository
func (f *ConsolidatedFactory) GetEventRepository() port.EventRepository {
	// TODO: implement when needed
	return nil
}

// GetTickerRepository returns a ticker repository
func (f *ConsolidatedFactory) GetTickerRepository() port.TickerRepository {
	// TODO: implement when needed
	return nil
}

// GetConversationMemoryRepository returns an AI conversation repository
func (f *ConsolidatedFactory) GetConversationMemoryRepository() port.ConversationMemoryRepository {
	// TODO: implement when needed
	return nil
}

// GetEmbeddingRepository returns an embedding repository
func (f *ConsolidatedFactory) GetEmbeddingRepository() port.EmbeddingRepository {
	// TODO: implement when needed
	return nil
}

// GetStrategyRepository returns a strategy repository
func (f *ConsolidatedFactory) GetStrategyRepository() port.StrategyRepository {
	// TODO: implement when needed
	return nil
}

// GetNotificationRepository returns a notification repository
func (f *ConsolidatedFactory) GetNotificationRepository() port.NotificationRepository {
	// TODO: implement when needed
	return nil
}

// GetAnalyticsRepository returns an analytics repository
func (f *ConsolidatedFactory) GetAnalyticsRepository() port.AnalyticsRepository {
	// TODO: implement when needed
	return nil
}

// GetSystemStatusRepository returns a system status repository
func (f *ConsolidatedFactory) GetSystemStatusRepository() port.SystemStatusRepository {
	// TODO: implement when needed
	return nil
}

// GetSymbolRepository returns a symbol repository
func (f *ConsolidatedFactory) GetSymbolRepository() port.SymbolRepository {
	// TODO: implement when needed
	return nil
}

// GetMarketDataRepository returns a market data repository
func (f *ConsolidatedFactory) GetMarketDataRepository() port.MarketDataRepository {
	// TODO: implement when needed
	return nil
}

// GetAPICredentialRepository returns an API credential repository
func (f *ConsolidatedFactory) GetAPICredentialRepository() port.APICredentialRepository {
	// TODO: implement when needed
	return nil
}

// --- Middleware Factory Methods ---

// GetUserRepository returns a user repository for auth services
func (f *ConsolidatedFactory) GetUserRepository() port.UserRepository {
	return repo.NewUserRepository(f.db, f.logger)
}

// GetUserService returns a user service
func (f *ConsolidatedFactory) GetUserService() service.UserServiceInterface {
	userRepo := f.GetUserRepository()
	return service.NewUserService(userRepo)
}

// GetAuthService returns an authentication service
func (f *ConsolidatedFactory) GetAuthService() (service.AuthServiceInterface, error) {
	userService := f.GetUserService()
	userServiceImpl, ok := userService.(*service.UserService)
	if !ok {
		return nil, fmt.Errorf("failed to cast UserServiceInterface to *UserService")
	}
	return service.NewAuthService(userServiceImpl, f.cfg.Auth.ClerkSecretKey)
}

// GetAuthMiddleware returns the authentication middleware
func (f *ConsolidatedFactory) GetAuthMiddleware() (middleware.AuthMiddleware, error) {
	authService, err := f.GetAuthService()
	if err != nil {
		return nil, err
	}
	return middleware.NewAuthMiddleware(authService, f.logger), nil
}

// GetTestAuthMiddleware returns the test authentication middleware
func (f *ConsolidatedFactory) GetTestAuthMiddleware() middleware.AuthMiddleware {
	return middleware.NewTestAuthMiddleware(f.logger)
}

// GetDisabledAuthMiddleware returns the disabled authentication middleware
func (f *ConsolidatedFactory) GetDisabledAuthMiddleware() middleware.AuthMiddleware {
	return middleware.NewDisabledAuthMiddleware(f.logger)
}

// GetRateLimiter returns an advanced rate limiter
func (f *ConsolidatedFactory) GetRateLimiter() *middleware.AdvancedRateLimiter {
	return middleware.NewAdvancedRateLimiter(&f.cfg.RateLimit, f.logger)
}

// GetRateLimiterMiddleware returns the rate limiter middleware
func (f *ConsolidatedFactory) GetRateLimiterMiddleware() func(http.Handler) http.Handler {
	limiter := f.GetRateLimiter()
	return middleware.AdvancedRateLimiterMiddleware(limiter)
}

// GetCSRFMiddleware returns a CSRF middleware
func (f *ConsolidatedFactory) GetCSRFMiddleware() *middleware.CSRFMiddleware {
	return middleware.NewCSRFMiddleware(&f.cfg.CSRF, f.logger)
}

// GetCSRFProtectionMiddleware returns the CSRF protection middleware
func (f *ConsolidatedFactory) GetCSRFProtectionMiddleware() func(http.Handler) http.Handler {
	csrfMiddleware := f.GetCSRFMiddleware()
	return csrfMiddleware.Middleware()
}

// GetSecureHeadersMiddleware returns a secure headers middleware
func (f *ConsolidatedFactory) GetSecureHeadersMiddleware() *middleware.SecureHeadersMiddleware {
	return middleware.NewSecureHeadersMiddleware(&f.cfg.SecureHeaders, f.logger)
}

// GetSecureHeadersHandler returns the secure headers handler
func (f *ConsolidatedFactory) GetSecureHeadersHandler() func(http.Handler) http.Handler {
	secureHeadersMiddleware := f.GetSecureHeadersMiddleware()
	return secureHeadersMiddleware.Middleware()
}

// GetUnifiedErrorMiddleware returns the unified error middleware
func (f *ConsolidatedFactory) GetUnifiedErrorMiddleware() *middleware.UnifiedErrorMiddleware {
	return middleware.NewUnifiedErrorMiddleware(f.logger)
}

// GetUnifiedErrorHandler returns the unified error handler middleware function
func (f *ConsolidatedFactory) GetUnifiedErrorHandler() func(http.Handler) http.Handler {
	errorMiddleware := f.GetUnifiedErrorMiddleware()
	return errorMiddleware.Middleware()
}

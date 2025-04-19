package di

import (
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/server"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/persistence/gorm/migrations"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/scheduler"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	portgateway "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/gateway"
	portservice "github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port/service"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/service"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/factory"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/usecase"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/util/crypto"
)

// Container holds all the application dependencies.
type Container struct {
	Config *config.Config
	Logger *zerolog.Logger
	DB     *gorm.DB

	// Factories (Pluralizing for consistency)
	AppFactory        *factory.AppFactory
	RepositoryFactory port.RepositoryFactory
	UseCaseFactory    *factory.UseCaseFactory

	// Crypto
	KeyManager        *crypto.KeyManager
	EncryptionService port.EnhancedEncryptionService

	// Repositories
	WalletRepo        port.WalletRepository
	NewCoinRepo       port.NewCoinRepository
	EventRepo         port.EventRepository
	SymbolRepo        port.SymbolRepository
	OrderRepo         port.OrderRepository
	MarketRepo        port.MarketRepository
	APICredentialRepo port.APICredentialRepository

	// Services (Consider renaming MexcSniperService if it's now a UseCase)
	WalletService     service.WalletService
	AuthService       port.AuthServiceInterface
	TriggerService    port.TriggerService
	MarketDataService port.MarketDataService
	TradeExecutor     port.TradeExecutor
	MexcSniperService *service.MexcSniperService
	SniperShotUseCase *usecase.SniperShotService
	WalletUseCase     port.WalletService
	ProtectedUseCase  portservice.ProtectedUserService

	// Event Bus
	EventBus port.EventBus

	// Gateways
	MEXCGateway  portgateway.MEXCGateway
	ClerkGateway portgateway.ClerkGateway
	AIGateway    portgateway.AIGateway

	// Middleware
	AuthFactory     *middleware.AuthFactory
	AuthMiddleware  middleware.AuthMiddleware
	ErrorMiddleware *middleware.UnifiedErrorMiddleware

	// Server
	Server *server.Server

	// Workers
	NewCoinWorker *scheduler.NewCoinWorker
}

// NewContainer initializes and returns a new dependency container.
func NewContainer() (*Container, error) {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".") // Assuming LoadConfig exists
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Initialize Logger
	logger := factory.NewLogger(cfg.Log) // Assuming NewLogger exists

	// 3. Initialize Database Connection
	db, err := factory.NewDBConnection(cfg, logger) // Assuming NewDBConnection exists
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 4. Run Database Migrations
	// Assuming a consolidated migration runner exists
	if err := migrations.RunAllMigrations(db, logger); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	// --- Initialize Crypto Services ---
	keyManager, err := crypto.NewKeyManager(cfg.Encryption.CurrentKeyID, cfg.Encryption.Keys, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize key manager: %w", err)
	}

	var encryptionService port.EnhancedEncryptionService = crypto.NewEnhancedEncryptionService(keyManager, logger)

	// --- Initialize Repository Factory ---
	repoFactory := provideRepositoryFactory(db, logger).WithEncryptionService(encryptionService)

	// --- Initialize Repositories using Factory ---
	walletRepo := provideWalletRepository(repoFactory)
	newCoinRepo := provideNewCoinRepository(repoFactory)
	eventRepo := provideEventRepository(repoFactory)
	symbolRepo := provideSymbolRepository(repoFactory)
	orderRepo := provideOrderRepository(repoFactory)
	marketRepo := provideMarketRepository(repoFactory)
	apiCredentialRepo := provideAPICredentialRepository(repoFactory)

	// --- Initialize Event Bus ---
	eventBus := provideEventBus(logger)

	// --- Initialize Gateways ---
	// Pass container itself to providers
	containerForProviders := &Container{Config: cfg, Logger: logger, DB: db}
	containerForProviders.RepositoryFactory = repoFactory
	containerForProviders.EncryptionService = encryptionService

	mexcGateway, err := ProvideMEXCGateway(containerForProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MEXC gateway: %w", err)
	}
	containerForProviders.MEXCGateway = mexcGateway

	clerkGateway, err := ProvideClerkGateway(containerForProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Clerk gateway: %w", err)
	}
	containerForProviders.ClerkGateway = clerkGateway

	aiGateway, err := ProvideAIGateway(containerForProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AI gateway: %w", err)
	}
	containerForProviders.AIGateway = aiGateway

	// --- Initialize Domain/Infra Services (dependencies for Use Cases/Factory) ---
	authService, err := provideAuthService(containerForProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AuthService: %w", err)
	}
	triggerService := provideTriggerService(logger)
	marketDataService, err := provideMarketDataService(containerForProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Market Data Service: %w", err)
	}
	tradeExecutor, err := provideTradeExecutor(containerForProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Trade Executor: %w", err)
	}

	// --- Initialize UseCase Factory ---
	useCaseFactory := factory.NewUseCaseFactory(
		logger,
		tradeExecutor,
		orderRepo,
		marketDataService,
		triggerService,
		walletRepo,
		authService,
	)

	// --- Build Use Cases using Factory ---
	sniperShotUseCase := useCaseFactory.BuildSniperShotService()
	walletUseCase := useCaseFactory.BuildWalletUseCase()
	protectedUseCase := useCaseFactory.BuildProtectedUserUseCase()

	// --- Initialize MexcSniperService (Domain Service) - Keep for now if still needed elsewhere ---
	mexcSniperService, err := provideMexcSniperService(containerForProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Mexc Sniper Service: %w", err)
	}

	// --- Initialize Middleware ---
	middlewareProviders := provideMiddlewares(cfg, logger, authService)

	// --- Initialize Workers ---
	newCoinWorker, err := provideNewCoinWorker(containerForProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NewCoinWorker: %w", err)
	}

	// --- Initialize Server ---
	httpServer := server.NewServer(db, cfg, logger, useCaseFactory)
	httpServer.SetMiddleware(middlewareProviders.AuthMiddleware, middlewareProviders.ErrorMiddleware)
	if err := httpServer.SetupRoutes(); err != nil {
		return nil, fmt.Errorf("failed to set up server routes: %w", err)
	}

	// --- Initialize AppFactory ---
	appFactory := factory.NewAppFactory(cfg, logger)

	// --- Assemble Final Container ---
	container := &Container{
		Config: cfg,
		Logger: logger,
		DB:     db,

		// Factories
		AppFactory:        appFactory,
		RepositoryFactory: repoFactory,
		UseCaseFactory:    useCaseFactory,

		// Crypto
		KeyManager:        keyManager,
		EncryptionService: encryptionService,

		// Repositories
		WalletRepo:        walletRepo,
		NewCoinRepo:       newCoinRepo,
		EventRepo:         eventRepo,
		SymbolRepo:        symbolRepo,
		OrderRepo:         orderRepo,
		MarketRepo:        marketRepo,
		APICredentialRepo: apiCredentialRepo,

		// Services & Use Cases
		AuthService:       authService,
		TriggerService:    triggerService,
		MarketDataService: marketDataService,
		TradeExecutor:     tradeExecutor,
		MexcSniperService: mexcSniperService,
		SniperShotUseCase: sniperShotUseCase,
		WalletUseCase:     walletUseCase,
		ProtectedUseCase:  protectedUseCase,

		// Event Bus
		EventBus: eventBus,

		// Gateways
		MEXCGateway:  mexcGateway,
		ClerkGateway: clerkGateway,
		AIGateway:    aiGateway,

		// Middleware
		AuthFactory:     middlewareProviders.AuthFactory,
		AuthMiddleware:  middlewareProviders.AuthMiddleware,
		ErrorMiddleware: middlewareProviders.ErrorMiddleware,

		// Server
		Server: httpServer,

		// Workers
		NewCoinWorker: newCoinWorker,
	}

	// Populate containerForProviders fully before returning the final container
	// This ensures all dependencies are available for any remaining provider calls
	containerForProviders.PopulateFrom(container)

	logger.Info().Msg("Dependency container initialized successfully")
	return container, nil
}

// PopulateFrom copies essential references from a fully built container.
// Used to ensure the temporary container used during provider calls has access to all needed services.
func (c *Container) PopulateFrom(source *Container) {
	if c.Config == nil {
		c.Config = source.Config
	}
	if c.Logger == nil {
		c.Logger = source.Logger
	}
	if c.DB == nil {
		c.DB = source.DB
	}
	if c.RepositoryFactory == nil {
		c.RepositoryFactory = source.RepositoryFactory
	}
	if c.EncryptionService == nil {
		c.EncryptionService = source.EncryptionService
	}
	if c.EventBus == nil {
		c.EventBus = source.EventBus
	}
	if c.MEXCGateway == nil {
		c.MEXCGateway = source.MEXCGateway
	}
	if c.ClerkGateway == nil {
		c.ClerkGateway = source.ClerkGateway
	}
	if c.AIGateway == nil {
		c.AIGateway = source.AIGateway
	}
	if c.MarketDataService == nil {
		c.MarketDataService = source.MarketDataService
	}
	if c.TradeExecutor == nil {
		c.TradeExecutor = source.TradeExecutor
	}
	if c.TriggerService == nil {
		c.TriggerService = source.TriggerService
	}
}

// --- Getter Methods ---
// (Add getters for dependencies needed externally, e.g., by main.go)

func (c *Container) GetConfig() *config.Config {
	return c.Config
}

func (c *Container) GetLogger() *zerolog.Logger {
	return c.Logger
}

func (c *Container) GetDB() *gorm.DB {
	return c.DB
}

func (c *Container) GetWalletService() service.WalletService {
	return c.WalletService
}

func (c *Container) GetTriggerService() port.TriggerService {
	return c.TriggerService
}

func (c *Container) GetEncryptionService() port.EnhancedEncryptionService {
	return c.EncryptionService
}

func (c *Container) GetAuthService() port.AuthServiceInterface {
	return c.AuthService
}

func (c *Container) GetAuthMiddleware() middleware.AuthMiddleware {
	return c.AuthMiddleware
}

func (c *Container) GetErrorMiddleware() *middleware.UnifiedErrorMiddleware {
	return c.ErrorMiddleware
}

func (c *Container) GetMEXCGateway() portgateway.MEXCGateway {
	return c.MEXCGateway
}

func (c *Container) GetClerkGateway() portgateway.ClerkGateway {
	return c.ClerkGateway
}

func (c *Container) GetAIGateway() portgateway.AIGateway {
	return c.AIGateway
}

func (c *Container) GetNewCoinRepository() port.NewCoinRepository {
	return c.NewCoinRepo
}

func (c *Container) GetEventRepository() port.EventRepository {
	return c.EventRepo
}

func (c *Container) GetEventBus() port.EventBus {
	return c.EventBus
}

func (c *Container) GetMexcSniperService() *service.MexcSniperService {
	return c.MexcSniperService
}

func (c *Container) GetNewCoinWorker() *scheduler.NewCoinWorker {
	return c.NewCoinWorker
}

func (c *Container) GetTradeExecutor() port.TradeExecutor {
	return c.TradeExecutor
}

// Getter for UseCaseFactory
func (c *Container) GetUseCaseFactory() *factory.UseCaseFactory {
	return c.UseCaseFactory
}

// Getter for SniperShotUseCase
func (c *Container) GetSniperShotUseCase() *usecase.SniperShotService {
	return c.SniperShotUseCase
}

// Getter for WalletUseCase
func (c *Container) GetWalletUseCase() port.WalletService {
	return c.WalletUseCase
}

// Getter for ProtectedUserUseCase
func (c *Container) GetProtectedUserUseCase() portservice.ProtectedUserService {
	return c.ProtectedUseCase
}

// Add other getters as needed...

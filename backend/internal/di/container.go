package di

import (
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/factory"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/service"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	gormdb "gorm.io/gorm"
)

// Container provides dependency injection for the application
type Container struct {
	config            *config.Config
	logger            *zerolog.Logger
	db                *gormdb.DB
	txManager         port.TransactionManager
	repositoryFactory *factory.RepositoryFactory
	useCaseFactory    *factory.UseCaseFactory

	// External services
	mexcClient    port.MEXCClient
	eventBus      port.EventBus
	sniperService port.SniperService

	// Repositories
	orderRepository          port.OrderRepository
	walletRepository         port.WalletRepository
	newCoinRepository        port.NewCoinRepository
	eventRepository          port.EventRepository
	tickerRepository         port.TickerRepository
	aiConversationRepository port.ConversationMemoryRepository
	embeddingRepository      port.EmbeddingRepository
	strategyRepository       port.StrategyRepository
	notificationRepository   port.NotificationRepository
	analyticsRepository      port.AnalyticsRepository
	statusRepository         port.SystemStatusRepository

	// Use cases
	tradeUseCase               usecase.TradeUseCase
	positionUseCase            usecase.PositionUseCase
	newCoinUseCase             usecase.NewCoinUseCase
	aiUseCase                  *usecase.AIUsecase
	statusUseCase              usecase.StatusUseCase
	sniperUseCase              usecase.SniperUseCase
	newListingDetectionService *service.NewListingDetectionService
}

// NewContainer creates a new dependency injection container
func NewContainer(cfg *config.Config, logger *zerolog.Logger, db *gormdb.DB) *Container {
	return &Container{
		config: cfg,
		logger: logger,
		db:     db,
	}
}

// Initialize initializes all dependencies
func (c *Container) Initialize() error {
	// Create transaction manager
	c.txManager = gorm.NewTransactionManager(c.db, c.logger)

	// Create repository factory
	c.repositoryFactory = factory.NewRepositoryFactory(c.db, c.logger, c.config)

	// Initialize repositories
	c.initializeRepositories()

	// Initialize external services
	c.initializeExternalServices()

	// Create use case factory
	c.useCaseFactory = factory.NewUseCaseFactory(
		c.config,
		c.logger,
		c.orderRepository,
		c.walletRepository,
		c.newCoinRepository,
		c.eventRepository,
		c.tickerRepository,
		c.aiConversationRepository,
		c.embeddingRepository,
		c.strategyRepository,
		c.notificationRepository,
		c.analyticsRepository,
		c.statusRepository,
		c.mexcClient,
		c.eventBus,
		c.txManager,
		c.sniperService,
	)

	// Initialize use cases
	c.initializeUseCases()

	return nil
}

// initializeRepositories initializes all repositories
func (c *Container) initializeRepositories() {
	c.orderRepository = c.repositoryFactory.CreateOrderRepository()
	c.walletRepository = c.repositoryFactory.CreateWalletRepository()
	c.newCoinRepository = c.repositoryFactory.CreateNewCoinRepository()
	c.eventRepository = c.repositoryFactory.CreateEventRepository()
	c.tickerRepository = c.repositoryFactory.CreateTickerRepository()
	c.aiConversationRepository = c.repositoryFactory.CreateAIConversationRepository()
	c.embeddingRepository = c.repositoryFactory.CreateEmbeddingRepository()
	c.strategyRepository = c.repositoryFactory.CreateStrategyRepository()
	c.notificationRepository = c.repositoryFactory.CreateNotificationRepository()
	c.analyticsRepository = c.repositoryFactory.CreateAnalyticsRepository()
	c.statusRepository = c.repositoryFactory.CreateStatusRepository()
}

// initializeExternalServices initializes external services
func (c *Container) initializeExternalServices() {
	// Create MEXC client
	mexcFactory := factory.NewMEXCFactory(c.config, c.logger)
	c.mexcClient = mexcFactory.CreateMEXCClient()

	// Create event bus
	c.eventBus = delivery.NewInMemoryEventBus(*c.logger)

	// Create sniper service
	// Create market factory
	marketFactory := factory.NewMarketFactory(c.config, c.logger, c.db)

	// Create market data service
	marketDataService := marketFactory.CreateMarketDataService()

	// Create new listing detection service first
	interval := 10 * time.Second
	if c.config.MEXC.RateLimit.RequestsPerMinute > 0 {
		interval = time.Minute / time.Duration(c.config.MEXC.RateLimit.RequestsPerMinute)
	}
	c.newListingDetectionService = service.NewNewListingDetectionService(
		c.newCoinRepository,
		c.eventRepository,
		c.eventBus,
		c.mexcClient,
		c.logger,
		service.NewListingDetectionConfig{
			RESTPollingInterval: interval,
			WebSocketEnabled:    c.config.MEXC.WSBaseURL != "",
			MaxQueueSize:        1000,
		},
	)
	if err := c.newListingDetectionService.Start(); err != nil {
		c.logger.Error().Err(err).Msg("Failed to start new listing detection service")
	}

	// Create sniper factory
	sniperFactory := factory.NewSniperFactory(
		c.mexcClient,
		c.repositoryFactory.CreateSymbolRepository(),
		c.orderRepository,
		marketDataService,
		c.newListingDetectionService,
		c.logger,
	)

	// Create real sniper service
	c.sniperService = sniperFactory.CreateSniperService()
}

// initializeUseCases initializes all use cases
func (c *Container) initializeUseCases() {
	c.tradeUseCase = c.useCaseFactory.CreateTradeUseCase()
	c.positionUseCase = c.useCaseFactory.CreatePositionUseCase()
	c.newCoinUseCase = c.useCaseFactory.CreateNewCoinUseCase()
	c.aiUseCase = c.useCaseFactory.CreateAIUseCase()
	c.statusUseCase = c.useCaseFactory.CreateStatusUseCase()
	c.sniperUseCase = c.useCaseFactory.CreateSniperUseCase()
}

// GetOrderRepository returns the order repository
func (c *Container) GetOrderRepository() port.OrderRepository {
	return c.orderRepository
}

// GetWalletRepository returns the wallet repository
func (c *Container) GetWalletRepository() port.WalletRepository {
	return c.walletRepository
}

// GetNewCoinRepository returns the new coin repository
func (c *Container) GetNewCoinRepository() port.NewCoinRepository {
	return c.newCoinRepository
}

// GetEventRepository returns the event repository
func (c *Container) GetEventRepository() port.EventRepository {
	return c.eventRepository
}

// GetTickerRepository returns the ticker repository
func (c *Container) GetTickerRepository() port.TickerRepository {
	return c.tickerRepository
}

// GetAIConversationRepository returns the AI conversation repository
func (c *Container) GetAIConversationRepository() port.ConversationMemoryRepository {
	return c.aiConversationRepository
}

// GetEmbeddingRepository returns the embedding repository
func (c *Container) GetEmbeddingRepository() port.EmbeddingRepository {
	return c.embeddingRepository
}

// GetStrategyRepository returns the strategy repository
func (c *Container) GetStrategyRepository() port.StrategyRepository {
	return c.strategyRepository
}

// GetNotificationRepository returns the notification repository
func (c *Container) GetNotificationRepository() port.NotificationRepository {
	return c.notificationRepository
}

// GetAnalyticsRepository returns the analytics repository
func (c *Container) GetAnalyticsRepository() port.AnalyticsRepository {
	return c.analyticsRepository
}

// GetStatusRepository returns the status repository
func (c *Container) GetStatusRepository() port.SystemStatusRepository {
	return c.statusRepository
}

// GetConfig returns the application configuration
func (c *Container) GetConfig() *config.Config {
	return c.config
}

// GetLogger returns the application logger
func (c *Container) GetLogger() *zerolog.Logger {
	return c.logger
}

// GetDB returns the database connection
func (c *Container) GetDB() *gormdb.DB {
	return c.db
}

// GetTransactionManager returns the transaction manager
func (c *Container) GetTransactionManager() port.TransactionManager {
	return c.txManager
}

// GetTradeUseCase returns the trade use case
func (c *Container) GetTradeUseCase() usecase.TradeUseCase {
	return c.tradeUseCase
}

// GetPositionUseCase returns the position use case
func (c *Container) GetPositionUseCase() usecase.PositionUseCase {
	return c.positionUseCase
}

// GetNewCoinUseCase returns the new coin use case
func (c *Container) GetNewCoinUseCase() usecase.NewCoinUseCase {
	return c.newCoinUseCase
}

// GetAIUseCase returns the AI use case
func (c *Container) GetAIUseCase() *usecase.AIUsecase {
	return c.aiUseCase
}

// GetStatusUseCase returns the status use case
func (c *Container) GetStatusUseCase() usecase.StatusUseCase {
	return c.statusUseCase
}

// GetMEXCClient returns the MEXC client
func (c *Container) GetMEXCClient() port.MEXCClient {
	return c.mexcClient
}

// GetEventBus returns the event bus
func (c *Container) GetEventBus() port.EventBus {
	return c.eventBus
}

// GetNewListingDetectionService returns the new listing detection service
func (c *Container) GetNewListingDetectionService() *service.NewListingDetectionService {
	return c.newListingDetectionService
}

// GetSniperUseCase returns the sniper use case
func (c *Container) GetSniperUseCase() usecase.SniperUseCase {
	return c.sniperUseCase
}

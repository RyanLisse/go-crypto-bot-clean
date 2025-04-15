package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
)

// UseCaseFactory creates use case instances
type UseCaseFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger

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

	// External services
	mexcClient port.MEXCClient
	eventBus   port.EventBus
	txManager  port.TransactionManager
}

// NewUseCaseFactory creates a new UseCaseFactory
func NewUseCaseFactory(
	cfg *config.Config,
	logger *zerolog.Logger,
	orderRepo port.OrderRepository,
	walletRepo port.WalletRepository,
	newCoinRepo port.NewCoinRepository,
	eventRepo port.EventRepository,
	tickerRepo port.TickerRepository,
	aiConversationRepo port.ConversationMemoryRepository,
	embeddingRepo port.EmbeddingRepository,
	strategyRepo port.StrategyRepository,
	notificationRepo port.NotificationRepository,
	analyticsRepo port.AnalyticsRepository,
	statusRepo port.SystemStatusRepository,
	mexcClient port.MEXCClient,
	eventBus port.EventBus,
	txManager port.TransactionManager,
) *UseCaseFactory {
	return &UseCaseFactory{
		cfg:                      cfg,
		logger:                   logger,
		orderRepository:          orderRepo,
		walletRepository:         walletRepo,
		newCoinRepository:        newCoinRepo,
		eventRepository:          eventRepo,
		tickerRepository:         tickerRepo,
		aiConversationRepository: aiConversationRepo,
		embeddingRepository:      embeddingRepo,
		strategyRepository:       strategyRepo,
		notificationRepository:   notificationRepo,
		analyticsRepository:      analyticsRepo,
		statusRepository:         statusRepo,
		mexcClient:               mexcClient,
		eventBus:                 eventBus,
		txManager:                txManager,
	}
}

// CreateTradeUseCase creates a trade use case
func (f *UseCaseFactory) CreateTradeUseCase() usecase.TradeUseCase {
	// Return a mock implementation for now to avoid dependency issues
	return &usecase.MockTradeUseCase{}
}

// CreatePositionUseCase creates a position use case
func (f *UseCaseFactory) CreatePositionUseCase() usecase.PositionUseCase {
	// Using a mock implementation for now
	return &usecase.MockPositionUseCase{}
}

// CreateNewCoinUseCase creates a new coin use case
func (f *UseCaseFactory) CreateNewCoinUseCase() *usecase.NewCoinUseCase {
	return usecase.NewNewCoinUseCase(
		f.newCoinRepository,
		f.eventRepository,
		f.eventBus,
		f.mexcClient,
		f.logger,
	)
}

// CreateAIUseCase creates an AI use case
func (f *UseCaseFactory) CreateAIUseCase() *usecase.AIUsecase {
	return usecase.NewAIUsecase(
		nil, // TODO: Replace with AI service
		f.aiConversationRepository,
		f.embeddingRepository,
		*f.logger,
	)
}

// CreateStatusUseCase creates a status use case
func (f *UseCaseFactory) CreateStatusUseCase() usecase.StatusUseCase {
	// Using a mock implementation for now
	return &usecase.MockStatusUseCase{}
}

package factory

import (
	gormadapter "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// RepositoryFactory creates repository instances
type RepositoryFactory struct {
	db     *gorm.DB
	logger *zerolog.Logger
	cfg    *config.Config
}

// NewRepositoryFactory creates a new RepositoryFactory
func NewRepositoryFactory(db *gorm.DB, logger *zerolog.Logger, cfg *config.Config) *RepositoryFactory {
	return &RepositoryFactory{
		db:     db,
		logger: logger,
		cfg:    cfg,
	}
}

// CreateOrderRepository creates an OrderRepository
func (f *RepositoryFactory) CreateOrderRepository() port.OrderRepository {
	// TODO: implement actual repository when needed
	return nil
}

// CreateWalletRepository creates a WalletRepository
func (f *RepositoryFactory) CreateWalletRepository() port.WalletRepository {
	return repo.NewConsolidatedWalletRepository(f.db, f.logger)
}

// CreateNewCoinRepository creates a NewCoinRepository
func (f *RepositoryFactory) CreateNewCoinRepository() port.NewCoinRepository {
   // Create a GORM-based NewCoinRepository
   return repo.NewGormNewCoinRepository(f.db, f.logger)
}

// CreateEventRepository creates an EventRepository
func (f *RepositoryFactory) CreateEventRepository() port.EventRepository {
   // Create a GORM-based EventRepository
   return repo.NewGormEventRepository(f.db, f.logger)
}

// CreateTickerRepository creates a TickerRepository
func (f *RepositoryFactory) CreateTickerRepository() port.TickerRepository {
   // Create a GORM-based TickerRepository
   return repo.NewGormTickerRepository(f.db, f.logger)
}

// CreateAIConversationRepository creates an AIConversationRepository
func (f *RepositoryFactory) CreateAIConversationRepository() port.ConversationMemoryRepository {
	// TODO: implement actual repository when needed
	return nil
}

// CreateEmbeddingRepository creates an EmbeddingRepository
func (f *RepositoryFactory) CreateEmbeddingRepository() port.EmbeddingRepository {
	// TODO: implement actual repository when needed
	return nil
}

// CreateStrategyRepository creates a StrategyRepository
func (f *RepositoryFactory) CreateStrategyRepository() port.StrategyRepository {
	// TODO: implement actual repository when needed
	return nil
}

// CreateNotificationRepository creates a NotificationRepository
func (f *RepositoryFactory) CreateNotificationRepository() port.NotificationRepository {
	// TODO: implement actual repository when needed
	return nil
}

// CreateAnalyticsRepository creates an AnalyticsRepository
func (f *RepositoryFactory) CreateAnalyticsRepository() port.AnalyticsRepository {
	// TODO: implement actual repository when needed
	return nil
}

// CreateStatusRepository creates a SystemStatusRepository
func (f *RepositoryFactory) CreateStatusRepository() port.SystemStatusRepository {
	// TODO: implement actual repository when needed
	return nil
}

// CreateSymbolRepository creates a SymbolRepository
func (f *RepositoryFactory) CreateSymbolRepository() port.SymbolRepository {
	return gormadapter.NewSymbolRepository(f.db, f.logger)
}

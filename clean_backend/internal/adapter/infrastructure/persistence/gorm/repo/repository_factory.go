package repo

import (
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// RepositoryFactory creates repository instances
type RepositoryFactory struct {
	db            *gorm.DB
	logger        *zerolog.Logger
	encryptionSvc port.EncryptionService
}

// NewRepositoryFactory creates a new RepositoryFactory
func NewRepositoryFactory(db *gorm.DB, logger *zerolog.Logger) *RepositoryFactory {
	return &RepositoryFactory{
		db:            db,
		logger:        logger,
		encryptionSvc: nil, // Will be set later if needed
	}
}

// WithEncryptionService sets the encryption service for the factory
func (f *RepositoryFactory) WithEncryptionService(encryptionSvc port.EncryptionService) *RepositoryFactory {
	f.encryptionSvc = encryptionSvc
	return f
}

// CreateMarketRepository creates a new MarketRepository
func (f *RepositoryFactory) CreateMarketRepository() port.MarketRepository {
	return NewMarketRepository(f.db, f.logger)
}

// CreateSymbolRepository creates a new SymbolRepository
func (f *RepositoryFactory) CreateSymbolRepository() port.SymbolRepository {
	return NewSymbolRepository(f.db, f.logger)
}

// CreateUserRepository creates a new UserRepository
func (f *RepositoryFactory) CreateUserRepository() port.UserRepository {
	return NewGormUserRepository(f.db, f.logger)
}

// CreateWalletRepository creates a new WalletRepository
func (f *RepositoryFactory) CreateWalletRepository() port.WalletRepository {
	return &WalletRepository{
		DB:     f.db,
		Logger: f.logger,
	}
}

// CreateOrderRepository creates a new OrderRepository
func (f *RepositoryFactory) CreateOrderRepository() port.OrderRepository {
	return NewOrderRepository(f.db, f.logger)
}

// CreateAPICredentialRepository creates a new APICredentialRepository
func (f *RepositoryFactory) CreateAPICredentialRepository() port.APICredentialRepository {
	// For now, return nil if encryption service is not set
	// In a real implementation, this should be handled more gracefully
	if f.encryptionSvc == nil {
		f.logger.Warn().Msg("Encryption service not set, APICredentialRepository will not be created")
		return nil
	}
	return NewAPICredentialRepository(f.db, f.encryptionSvc, f.logger)
}

// CreateRiskRepository creates a new RiskRepository
// Uncomment when RiskRepository is implemented
// func (f *RepositoryFactory) CreateRiskRepository() port.RiskRepository {
// 	return NewRiskRepository(f.db, f.logger)
// }

package repositories

import (
	"go-crypto-bot-clean/backend/internal/domain/interfaces"
	"go-crypto-bot-clean/backend/internal/logging"

	"gorm.io/gorm"
)

// RepositoryFactory creates and provides access to all repositories
type RepositoryFactory struct {
	db     *gorm.DB
	logger *logging.LoggerWrapper

	positionRepository    interfaces.PositionRepository
	transactionRepository interfaces.TransactionRepository
	// Add other repositories as they're implemented
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db *gorm.DB, logger *logging.LoggerWrapper) *RepositoryFactory {
	factory := &RepositoryFactory{
		db:     db,
		logger: logger,
	}

	// Initialize repositories
	factory.initRepositories()

	return factory
}

// Initialize all repositories
func (f *RepositoryFactory) initRepositories() {
	f.positionRepository = NewGormPositionRepository(f.db, f.logger)
	f.transactionRepository = NewGormTransactionRepository(f.db, f.logger)
	// Initialize other repositories as they're implemented
}

// GetPositionRepository returns the position repository
func (f *RepositoryFactory) GetPositionRepository() interfaces.PositionRepository {
	return f.positionRepository
}

// GetTransactionRepository returns the transaction repository
func (f *RepositoryFactory) GetTransactionRepository() interfaces.TransactionRepository {
	return f.transactionRepository
}

// Add other getter methods for additional repositories

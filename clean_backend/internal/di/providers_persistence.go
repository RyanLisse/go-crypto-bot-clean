package di

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
)

// provideRepositoryFactory creates and returns a repository factory
func provideRepositoryFactory(db *gorm.DB, logger *zerolog.Logger) *repo.RepositoryFactory {
	return repo.NewRepositoryFactory(db, logger)
}

// provideWalletRepository creates and returns a new GORM wallet repository
func provideWalletRepository(factory *repo.RepositoryFactory) port.WalletRepository {
	return factory.CreateWalletRepository()
}

// provideUserRepository creates and returns a new GORM user repository
func provideUserRepository(factory *repo.RepositoryFactory) port.UserRepository {
	return factory.CreateUserRepository()
}

// provideMarketRepository creates and returns a new GORM market repository
func provideMarketRepository(factory *repo.RepositoryFactory) port.MarketRepository {
	return factory.CreateMarketRepository()
}

// provideSymbolRepository creates and returns a new GORM symbol repository
func provideSymbolRepository(factory *repo.RepositoryFactory) port.SymbolRepository {
	return factory.CreateSymbolRepository()
}

// provideOrderRepository creates and returns a new GORM order repository
func provideOrderRepository(factory *repo.RepositoryFactory) port.OrderRepository {
	return factory.CreateOrderRepository()
}

// provideAPICredentialRepository creates and returns a new GORM API credential repository
func provideAPICredentialRepository(factory *repo.RepositoryFactory) port.APICredentialRepository {
	return factory.CreateAPICredentialRepository()
}

// Uncomment when RiskRepository is implemented
// provideRiskRepository creates and returns a new GORM risk repository
// func provideRiskRepository(factory *repo.RepositoryFactory) port.RiskRepository {
// 	return factory.CreateRiskRepository()
// }

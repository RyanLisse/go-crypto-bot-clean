package di

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/infrastructure/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
)

// provideRepositoryFactory creates and returns a repository factory
func provideRepositoryFactory(db *gorm.DB, logger *zerolog.Logger) port.RepositoryFactory {
	return repo.NewRepositoryFactory(db, logger)
}

// provideWalletRepository creates and returns a new GORM wallet repository
func provideWalletRepository(factory port.RepositoryFactory) port.WalletRepository {
	return factory.CreateWalletRepository()
}

// provideUserRepository creates and returns a new GORM user repository
func provideUserRepository(factory port.RepositoryFactory) port.UserRepository {
	return factory.CreateUserRepository()
}

// provideMarketRepository creates and returns a new GORM market repository
func provideMarketRepository(factory port.RepositoryFactory) port.MarketRepository {
	return factory.CreateMarketRepository()
}

// provideSymbolRepository creates and returns a new GORM symbol repository
func provideSymbolRepository(factory port.RepositoryFactory) port.SymbolRepository {
	return factory.CreateSymbolRepository()
}

// provideOrderRepository creates and returns a new GORM order repository
func provideOrderRepository(factory port.RepositoryFactory) port.OrderRepository {
	return factory.CreateOrderRepository()
}

// provideNewCoinRepository creates and returns a new GORM new coin repository
func provideNewCoinRepository(factory port.RepositoryFactory) port.NewCoinRepository {
	return factory.CreateNewCoinRepository()
}

// provideEventRepository creates and returns a new GORM event repository
func provideEventRepository(factory port.RepositoryFactory) port.EventRepository {
	return factory.CreateEventRepository()
}

// provideAPICredentialRepository creates and returns a new GORM API credential repository
func provideAPICredentialRepository(factory port.RepositoryFactory) port.APICredentialRepository {
	return factory.CreateAPICredentialRepository()
}

// Uncomment when RiskRepository is implemented
// provideRiskRepository creates and returns a new GORM risk repository
// func provideRiskRepository(factory *repo.RepositoryFactory) port.RiskRepository {
// 	return factory.CreateRiskRepository()
// }

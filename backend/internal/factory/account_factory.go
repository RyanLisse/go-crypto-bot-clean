package factory

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// AccountFactory creates account-related components
type AccountFactory struct {
	cfg    *config.Config
	logger *zerolog.Logger
	db     *gorm.DB
}

// NewAccountFactory creates a new AccountFactory
func NewAccountFactory(cfg *config.Config, logger *zerolog.Logger, db *gorm.DB) *AccountFactory {
	return &AccountFactory{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// CreateAccountRepository creates an account repository
func (f *AccountFactory) CreateAccountRepository() port.WalletRepository {
	// Use the consolidated wallet repository implementation
	return repo.NewConsolidatedWalletRepository(f.db, f.logger)
}

// CreateAccountUseCase creates an account use case
func (f *AccountFactory) CreateAccountUseCase(mexcClient port.MEXCClient) usecase.AccountUsecase {
	// Create dependencies
	accountRepo := f.CreateAccountRepository()

	// Create use case
	return usecase.NewAccountUsecase(mexcClient, accountRepo, f.logger.With().Str("component", "account_usecase").Logger())
}

// CreateAccountHandler creates an account handler
func (f *AccountFactory) CreateAccountHandler(mexcClient port.MEXCClient) *handler.AccountHandler {
	// Create use case
	accountUseCase := f.CreateAccountUseCase(mexcClient)

	// Create handler
	return handler.NewAccountHandler(accountUseCase, f.logger)
}

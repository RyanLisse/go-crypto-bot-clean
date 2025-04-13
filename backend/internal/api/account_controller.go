package api

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/api/controllers"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/ports"
	"go-crypto-bot-clean/backend/internal/infrastructure/adapters"
	"go-crypto-bot-clean/backend/internal/platform/mexc"
	"go-crypto-bot-clean/backend/internal/platform/mexc/rest"

	"go.uber.org/zap"
)

// InitializeAccountController initializes the account controller and registers its routes
func (d *Dependencies) InitializeAccountController() {
	d.logger.Info("Initializing account controller")

	// Create an adapter for the account service
	var accountServiceAdapter ports.AccountService
	if d.AccountService != nil {
		// Create a simple adapter that converts the Balance return type to float64
		accountServiceAdapter = &accountServicePortAdapter{
			accountService: d.AccountService,
		}
	} else {
		d.logger.Error("Account service is nil")
		return
	}

	// Type assert the mexcClient to get access to the REST client
	mexcClientConcrete, ok := d.mexcClient.(*mexc.Client)
	if !ok {
		d.logger.Error("Failed to type assert mexcClient to *mexc.Client")
		return
	}

	// Get the REST client
	restClientInterface := mexcClientConcrete.GetRestClient()
	// No need to check for nil as GetRestClient never returns nil
	d.logger.Debug("Using REST client", zap.String("client_type", fmt.Sprintf("%T", restClientInterface)))

	// Type assert the REST client
	restClient, ok := restClientInterface.(*rest.Client)
	if !ok {
		d.logger.Error("Failed to type assert restClient to *rest.Client")
		return
	}

	// Create MEXC client adapter using the real client
	mexcClientAdapter := adapters.NewMEXCClientAdapter(restClient, d.logger)
	d.logger.Info("Using real MEXC client with API keys from configuration")

	// Create account controller
	d.AccountController = controllers.NewAccountController(accountServiceAdapter, mexcClientAdapter, d.logger)

	// Register account controller routes in the consolidated router
	if d.router != nil {
		d.logger.Info("Registering account controller routes")
		d.AccountController.RegisterRoutes(d.router)
	} else {
		d.logger.Warn("Router is nil, cannot register account controller routes")
	}
}

// accountServicePortAdapter adapts account.AccountService to ports.AccountService
type accountServicePortAdapter struct {
	accountService interface {
		GetAccountBalance(ctx context.Context) (models.Balance, error)
		GetWallet(ctx context.Context) (*models.Wallet, error)
		ValidateAPIKeys(ctx context.Context) (bool, error)
	}
}

// GetAccountBalance adapts the return type from models.Balance to float64
func (a *accountServicePortAdapter) GetAccountBalance(ctx context.Context) (float64, error) {
	balance, err := a.accountService.GetAccountBalance(ctx)
	if err != nil {
		return 0, err
	}
	return balance.Fiat, nil
}

// GetWallet passes through the GetWallet method
func (a *accountServicePortAdapter) GetWallet(ctx context.Context) (*models.Wallet, error) {
	return a.accountService.GetWallet(ctx)
}

// ValidateAPIKeys passes through the ValidateAPIKeys method
func (a *accountServicePortAdapter) ValidateAPIKeys(ctx context.Context) (bool, error) {
	return a.accountService.ValidateAPIKeys(ctx)
}

// GetBalance implements the GetBalance method required by the AccountService interface
func (a *accountServicePortAdapter) GetBalance(ctx context.Context) (models.Balance, error) {
	// Get the balance directly from accountService
	return a.accountService.GetAccountBalance(ctx)
}

// FetchBalances is not used but required by the interface
func (a *accountServicePortAdapter) FetchBalances(ctx context.Context) (models.Balance, error) {
	return a.accountService.GetAccountBalance(ctx)
}

// GetPortfolioValue implements the GetPortfolioValue method required by the AccountService interface
func (a *accountServicePortAdapter) GetPortfolioValue(ctx context.Context) (float64, error) {
	// Get wallet and calculate total value
	wallet, err := a.GetWallet(ctx)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, balance := range wallet.Balances {
		total += balance.Total * balance.Price
	}

	return total, nil
}

// GetBalanceSummary implements the GetBalanceSummary method required by the AccountService interface
func (a *accountServicePortAdapter) GetBalanceSummary(ctx context.Context, days int) (*models.BalanceSummary, error) {
	// Return a simple balance summary with current data
	balance, err := a.GetBalance(ctx)
	if err != nil {
		return nil, err
	}

	return &models.BalanceSummary{
		CurrentBalance:   balance.Fiat,
		Period:           days,
		GeneratedAt:      time.Now(),
		TransactionCount: 0,
		Deposits:         0,
		Withdrawals:      0,
		NetChange:        0,
	}, nil
}

// GetTransactionHistory implements the GetTransactionHistory method required by the AccountService interface
func (a *accountServicePortAdapter) GetTransactionHistory(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	// Return empty transaction history for now
	return []*models.Transaction{}, nil
}

// UpdateBalance implements the UpdateBalance method required by the AccountService interface
func (a *accountServicePortAdapter) UpdateBalance(ctx context.Context, amount float64, reason string) error {
	// No-op implementation
	return nil
}

// SyncWithExchange implements the SyncWithExchange method required by the AccountService interface
func (a *accountServicePortAdapter) SyncWithExchange(ctx context.Context) error {
	// No-op implementation
	return nil
}

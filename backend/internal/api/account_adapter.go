package api

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/core/account"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// AccountServiceAdapter adapts a *account.SimpleAccountService to the account.AccountService interface
type AccountServiceAdapter struct {
	service *account.SimpleAccountService
}

// NewAccountServiceAdapter creates a new adapter for the SimpleAccountService
func NewAccountServiceAdapter(service *account.SimpleAccountService) *AccountServiceAdapter {
	return &AccountServiceAdapter{
		service: service,
	}
}

// GetAccountBalance converts account.Balance to models.Balance
func (a *AccountServiceAdapter) GetAccountBalance(ctx context.Context) (models.Balance, error) {
	balance, err := a.service.GetAccountBalance(ctx)
	if err != nil {
		return models.Balance{}, err
	}

	// Convert from account.Balance to models.Balance
	return models.Balance{
		Fiat:      balance.Fiat,
		Assets:    make(map[string]float64), // Initialize empty map if not available in source
		Available: balance.Available,
		Locked:    make(map[string]float64), // Initialize empty map if not available in source
		UpdatedAt: time.Now(),
	}, nil
}

// GetWallet adapts the account.Wallet to models.Wallet
func (a *AccountServiceAdapter) GetWallet(ctx context.Context) (*models.Wallet, error) {
	accWallet, err := a.service.GetWallet(ctx)
	if err != nil {
		return nil, err
	}

	// Convert from account.Wallet to models.Wallet
	wallet := &models.Wallet{
		Balances:  make(map[string]*models.AssetBalance),
		UpdatedAt: accWallet.UpdatedAt,
	}

	// Copy balances
	for asset, balance := range accWallet.Balances {
		wallet.Balances[asset] = &models.AssetBalance{
			Asset:  balance.Asset,
			Free:   balance.Free,
			Locked: balance.Locked,
			Total:  balance.Total,
		}
	}

	return wallet, nil
}

// GetPortfolioValue delegates to the underlying service
func (a *AccountServiceAdapter) GetPortfolioValue(ctx context.Context) (float64, error) {
	return a.service.GetPortfolioValue(ctx)
}

// GetPositionRisk is not implemented in SimpleAccountService, so we return a default value
func (a *AccountServiceAdapter) GetPositionRisk(ctx context.Context, symbol string) (models.PositionRisk, error) {
	return models.PositionRisk{
		Symbol:    symbol,
		RiskLevel: "LOW",
	}, nil
}

// GetAllPositionRisks is not implemented in SimpleAccountService, so we return a default value
func (a *AccountServiceAdapter) GetAllPositionRisks(ctx context.Context) (map[string]models.PositionRisk, error) {
	return map[string]models.PositionRisk{}, nil
}

// GetCurrentExposure delegates to the underlying service
func (a *AccountServiceAdapter) GetCurrentExposure(ctx context.Context) (float64, error) {
	return a.service.GetCurrentExposure(ctx)
}

// ValidateAPIKeys delegates to the underlying service
func (a *AccountServiceAdapter) ValidateAPIKeys(ctx context.Context) (bool, error) {
	return a.service.ValidateAPIKeys(ctx)
}

// UpdateBalance is not implemented in SimpleAccountService, so we provide a stub
func (a *AccountServiceAdapter) UpdateBalance(ctx context.Context, amount float64, reason string) error {
	return nil
}

// SyncWithExchange is not implemented in SimpleAccountService, so we provide a stub
func (a *AccountServiceAdapter) SyncWithExchange(ctx context.Context) error {
	return nil
}

// GetBalanceSummary is not implemented in SimpleAccountService, so we provide a stub
func (a *AccountServiceAdapter) GetBalanceSummary(ctx context.Context, days int) (*models.BalanceSummary, error) {
	return &models.BalanceSummary{
		CurrentBalance:   0,
		Deposits:         0,
		Withdrawals:      0,
		NetChange:        0,
		TransactionCount: 0,
		Period:           days,
	}, nil
}

// GetTransactionHistory is not implemented in SimpleAccountService, so we provide a stub
func (a *AccountServiceAdapter) GetTransactionHistory(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	return []*models.Transaction{}, nil
}

// AnalyzeTransactions creates a new models.TransactionAnalysis
// We don't need to convert since SimpleAccountService already returns a models.TransactionAnalysis
func (a *AccountServiceAdapter) AnalyzeTransactions(ctx context.Context, startTime, endTime time.Time) (*models.TransactionAnalysis, error) {
	// SimpleAccountService.AnalyzeTransactions already returns the correct type
	return a.service.AnalyzeTransactions(ctx, startTime, endTime)
}

// SubscribeToBalanceUpdates is not implemented in SimpleAccountService, so we provide a stub
func (a *AccountServiceAdapter) SubscribeToBalanceUpdates(ctx context.Context, callback func(*models.Wallet)) error {
	return nil
}

// GetListenKey delegates to the underlying service
func (a *AccountServiceAdapter) GetListenKey(ctx context.Context) (string, error) {
	return a.service.GetListenKey(ctx)
}

// RenewListenKey delegates to the underlying service
func (a *AccountServiceAdapter) RenewListenKey(ctx context.Context, listenKey string) error {
	return a.service.RenewListenKey(ctx, listenKey)
}

// CloseListenKey delegates to the underlying service
func (a *AccountServiceAdapter) CloseListenKey(ctx context.Context, listenKey string) error {
	return a.service.CloseListenKey(ctx, listenKey)
}

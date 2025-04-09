package account

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
)

// Default cache time-to-live
const defaultCacheTTL = 5 * time.Minute

// AccountService defines the interface for account management operations
type AccountService interface {
	// Core account operations
	GetAccountBalance(ctx context.Context) (models.Balance, error)
	GetWallet(ctx context.Context) (*models.Wallet, error)
	GetPortfolioValue(ctx context.Context) (float64, error)

	// Risk assessment
	GetPositionRisk(ctx context.Context, symbol string) (models.PositionRisk, error)
	GetAllPositionRisks(ctx context.Context) (map[string]models.PositionRisk, error)
	GetCurrentExposure(ctx context.Context) (float64, error)

	// Validation
	ValidateAPIKeys(ctx context.Context) (bool, error)

	// Enhanced account management
	UpdateBalance(ctx context.Context, amount float64, reason string) error
	SyncWithExchange(ctx context.Context) error
	GetBalanceSummary(ctx context.Context, days int) (*models.BalanceSummary, error)
	GetTransactionHistory(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error)
	AnalyzeTransactions(ctx context.Context, startTime, endTime time.Time) (*models.TransactionAnalysis, error)

	// Real-time updates
	SubscribeToBalanceUpdates(ctx context.Context, callback func(*models.Wallet)) error
}

// MexcClient defines the interface for interacting with the MEXC exchange API
type MexcClient interface {
	FetchBalances(ctx context.Context) (models.Balance, error)
	ValidateKeys(ctx context.Context) (bool, error)
	GetAccountBalance(ctx context.Context) (float64, error)
}

// BoughtCoinRepository defines the interface for accessing bought coin data
type BoughtCoinRepository interface {
	GetAll(ctx context.Context) ([]string, error)
	GetPosition(ctx context.Context, symbol string) (float64, error)
}

// Config defines the interface for accessing configuration values
type Config interface {
	GetRiskThreshold() float64
	GetCacheTTL() time.Duration
}

// BalanceService defines the interface for the balance service
type BalanceService interface {
	GetLatestBalance(ctx context.Context) (*models.Wallet, error)
	SyncWithExchange(ctx context.Context) error
	SubscribeToUpdates(callback func(*models.Wallet))
	Start(ctx context.Context)
	Stop()
}

// accountService implements the AccountService interface
type accountService struct {
	mexcClient      MexcClient
	coinRepo        BoughtCoinRepository
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	config          Config
	balanceService  BalanceService

	// Cache mechanism
	cacheTTL        time.Duration
	mutex           sync.RWMutex
	balanceCache    *models.Balance
	balanceCacheExp time.Time
	walletCache     *models.Wallet
	walletCacheExp  time.Time

	// Subscribers
	balanceSubscribers []func(*models.Wallet)
	subMutex           sync.RWMutex
}

// NewAccountService creates a new instance of AccountService
func NewAccountService(
	mexcClient MexcClient,
	coinRepo BoughtCoinRepository,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	config Config,
) AccountService {
	cacheTTL := config.GetCacheTTL()
	if cacheTTL <= 0 {
		cacheTTL = defaultCacheTTL
	}

	return &accountService{
		mexcClient:         mexcClient,
		coinRepo:           coinRepo,
		walletRepo:         walletRepo,
		transactionRepo:    transactionRepo,
		config:             config,
		cacheTTL:           cacheTTL,
		balanceCache:       nil,
		balanceCacheExp:    time.Time{},
		walletCache:        nil,
		walletCacheExp:     time.Time{},
		balanceSubscribers: make([]func(*models.Wallet), 0),
	}
}

// GetAccountBalance retrieves the current account balance from the exchange
func (s *accountService) GetAccountBalance(ctx context.Context) (models.Balance, error) {
	// Check cache first for better performance
	s.mutex.RLock()
	if s.balanceCache != nil && time.Now().Before(s.balanceCacheExp) {
		balance := *s.balanceCache
		s.mutex.RUnlock()
		return balance, nil
	}
	s.mutex.RUnlock()

	// Get balance from exchange if cache is invalid
	balance, err := s.mexcClient.FetchBalances(ctx)
	if err != nil {
		return models.Balance{}, err
	}

	// Update cache
	s.mutex.Lock()
	s.balanceCache = &balance
	s.balanceCacheExp = time.Now().Add(s.cacheTTL)
	s.mutex.Unlock()

	return balance, nil
}

// GetWallet retrieves the current wallet from the repository
func (s *accountService) GetWallet(ctx context.Context) (*models.Wallet, error) {
	// If balance service is available, use it for real-time data
	if s.balanceService != nil {
		return s.balanceService.GetLatestBalance(ctx)
	}

	// Otherwise, use the traditional approach
	// Check cache first for better performance
	s.mutex.RLock()
	if s.walletCache != nil && time.Now().Before(s.walletCacheExp) {
		wallet := s.walletCache
		s.mutex.RUnlock()
		return wallet, nil
	}
	s.mutex.RUnlock()

	// Get wallet from repository if cache is invalid
	wallet, err := s.walletRepo.GetWallet(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// Initialize wallet if not exists
			wallet = &models.Wallet{
				Balances:  make(map[string]*models.AssetBalance),
				UpdatedAt: time.Now(),
			}

			if _, err := s.walletRepo.SaveWallet(ctx, wallet); err != nil {
				return nil, fmt.Errorf("error initializing wallet: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error getting wallet: %w", err)
		}
	}

	// Update cache
	s.mutex.Lock()
	s.walletCache = wallet
	s.walletCacheExp = time.Now().Add(s.cacheTTL)
	s.mutex.Unlock()

	return wallet, nil
}

// GetPortfolioValue calculates the total portfolio value including fiat and crypto
func (s *accountService) GetPortfolioValue(ctx context.Context) (float64, error) {
	balance, err := s.GetAccountBalance(ctx)
	if err != nil {
		return 0, err
	}

	total := balance.Fiat

	symbols, err := s.coinRepo.GetAll(ctx)
	if err != nil {
		return 0, err
	}

	for _, sym := range symbols {
		val, err := s.coinRepo.GetPosition(ctx, sym)
		if err != nil {
			return 0, err
		}
		total += val
	}

	return total, nil
}

// GetPositionRisk assesses the risk level of a specific position
func (s *accountService) GetPositionRisk(ctx context.Context, symbol string) (models.PositionRisk, error) {
	exposure, err := s.coinRepo.GetPosition(ctx, symbol)
	if err != nil {
		return models.PositionRisk{}, err
	}

	riskLevel := "LOW"
	if exposure > s.config.GetRiskThreshold() {
		riskLevel = "HIGH"
	}

	return models.PositionRisk{
		Symbol:      symbol,
		ExposureUSD: exposure,
		RiskLevel:   riskLevel,
	}, nil
}

// GetAllPositionRisks retrieves risk assessments for all positions
func (s *accountService) GetAllPositionRisks(ctx context.Context) (map[string]models.PositionRisk, error) {
	symbols, err := s.coinRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]models.PositionRisk)
	for _, sym := range symbols {
		risk, err := s.GetPositionRisk(ctx, sym)
		if err != nil {
			return nil, err
		}
		result[sym] = risk
	}

	return result, nil
}

// ValidateAPIKeys verifies that the API keys are valid
func (s *accountService) ValidateAPIKeys(ctx context.Context) (bool, error) {
	return s.mexcClient.ValidateKeys(ctx)
}

// GetCurrentExposure calculates the total current exposure across all positions
func (s *accountService) GetCurrentExposure(ctx context.Context) (float64, error) {
	symbols, err := s.coinRepo.GetAll(ctx)
	if err != nil {
		return 0, err
	}

	totalExposure := 0.0
	for _, sym := range symbols {
		exposure, err := s.coinRepo.GetPosition(ctx, sym)
		if err != nil {
			return 0, err
		}
		totalExposure += exposure
	}

	return totalExposure, nil
}

// UpdateBalance updates the wallet balance and records a transaction
func (s *accountService) UpdateBalance(ctx context.Context, amount float64, reason string) error {
	// Get current wallet
	wallet, err := s.GetWallet(ctx)
	if err != nil {
		return fmt.Errorf("error getting wallet: %w", err)
	}

	// Get current balance from exchange
	balance, err := s.GetAccountBalance(ctx)
	if err != nil {
		return fmt.Errorf("error getting account balance: %w", err)
	}

	// Update wallet with new balance
	// In a real implementation, we would update specific assets based on the transaction
	// For now, we'll just update the USDT balance as an example
	usdtBalance := balance.Fiat + amount

	// Update wallet
	if wallet.Balances == nil {
		wallet.Balances = make(map[string]*models.AssetBalance)
	}

	now := time.Now()
	wallet.Balances["USDT"] = &models.AssetBalance{
		Asset:  "USDT",
		Free:   usdtBalance,
		Locked: 0,
		Total:  usdtBalance,
	}

	wallet.UpdatedAt = now

	// Save updated wallet
	updatedWallet, err := s.walletRepo.SaveWallet(ctx, wallet)
	if err != nil {
		return fmt.Errorf("error saving wallet: %w", err)
	}

	// Record transaction
	transaction := &models.Transaction{
		Amount:    amount,
		Balance:   usdtBalance,
		Reason:    reason,
		Timestamp: time.Now(),
	}

	if _, err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("error recording transaction: %w", err)
	}

	// Update cache
	s.mutex.Lock()
	s.walletCache = updatedWallet
	s.walletCacheExp = time.Now().Add(s.cacheTTL)
	s.mutex.Unlock()

	return nil
}

// SyncWithExchange synchronizes the local wallet balance with the exchange
func (s *accountService) SyncWithExchange(ctx context.Context) error {
	// If balance service is available, use it
	if s.balanceService != nil {
		return s.balanceService.SyncWithExchange(ctx)
	}

	// Otherwise, use the traditional approach
	// Get balance from exchange
	exchangeBalance, err := s.mexcClient.GetAccountBalance(ctx)
	if err != nil {
		return fmt.Errorf("error getting exchange balance: %w", err)
	}

	// Get local wallet
	wallet, err := s.GetWallet(ctx)
	if err != nil {
		return fmt.Errorf("error getting local wallet: %w", err)
	}

	// Calculate current local balance
	var currentBalance float64
	if usdtBalance, ok := wallet.Balances["USDT"]; ok {
		currentBalance = usdtBalance.Free
	}

	// If there's a discrepancy, update local wallet and record transaction
	if currentBalance != exchangeBalance {
		difference := exchangeBalance - currentBalance
		reason := "Balance sync with exchange"

		if err := s.UpdateBalance(ctx, difference, reason); err != nil {
			return fmt.Errorf("error updating balance: %w", err)
		}
	}

	return nil
}

// GetBalanceSummary generates a summary of the wallet and transactions
func (s *accountService) GetBalanceSummary(ctx context.Context, days int) (*models.BalanceSummary, error) {
	// Get current balance
	balance, err := s.GetAccountBalance(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %w", err)
	}

	// Calculate time period
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)

	// Get transactions for period
	transactions, err := s.GetTransactionHistory(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions: %w", err)
	}

	// Calculate metrics
	var deposits, withdrawals float64

	for _, tx := range transactions {
		if tx.Amount > 0 {
			deposits += tx.Amount
		} else {
			withdrawals += -tx.Amount
		}
	}

	// Construct summary
	summary := &models.BalanceSummary{
		CurrentBalance:   balance.Fiat,
		Deposits:         deposits,
		Withdrawals:      withdrawals,
		NetChange:        deposits - withdrawals,
		TransactionCount: len(transactions),
		Period:           days,
		GeneratedAt:      time.Now(),
	}

	return summary, nil
}

// GetTransactionHistory retrieves transaction history for a specified period
func (s *accountService) GetTransactionHistory(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	if endTime.IsZero() {
		endTime = time.Now()
	}

	if startTime.IsZero() || startTime.After(endTime) {
		return nil, errors.New("invalid time range")
	}

	return s.transactionRepo.FindByTimeRange(ctx, startTime, endTime)
}

// AnalyzeTransactions performs analysis on transaction data
func (s *accountService) AnalyzeTransactions(ctx context.Context, startTime, endTime time.Time) (*models.TransactionAnalysis, error) {
	transactions, err := s.GetTransactionHistory(ctx, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions: %w", err)
	}

	if len(transactions) == 0 {
		return &models.TransactionAnalysis{
			StartTime:   startTime,
			EndTime:     endTime,
			TotalCount:  0,
			TotalVolume: 0,
		}, nil
	}

	var buys, sells int
	var buyVolume, sellVolume float64

	for _, tx := range transactions {
		if isBuyTransaction(tx.Reason) {
			buys++
			buyVolume += tx.Amount
		} else if isSellTransaction(tx.Reason) {
			sells++
			sellVolume += -tx.Amount
		}
	}

	// Create analysis result
	analysis := &models.TransactionAnalysis{
		StartTime:   startTime,
		EndTime:     endTime,
		TotalCount:  len(transactions),
		BuyCount:    buys,
		SellCount:   sells,
		TotalVolume: buyVolume + sellVolume,
		BuyVolume:   buyVolume,
		SellVolume:  sellVolume,
	}

	return analysis, nil
}

// Helper to determine transaction type from reason
func isBuyTransaction(reason string) bool {
	return strings.Contains(strings.ToLower(reason), "purchase") ||
		strings.Contains(strings.ToLower(reason), "buy") ||
		strings.Contains(strings.ToLower(reason), "deposit")
}

func isSellTransaction(reason string) bool {
	return strings.Contains(strings.ToLower(reason), "sale") ||
		strings.Contains(strings.ToLower(reason), "sell") ||
		strings.Contains(strings.ToLower(reason), "withdrawal")
}

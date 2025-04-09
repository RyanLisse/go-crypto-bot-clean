package account

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
	"go.uber.org/zap"
)

// MexcRESTClient defines the interface for interacting with the MEXC REST API
type MexcRESTClient interface {
	FetchBalances(ctx context.Context) (models.Balance, error)
	ValidateKeys(ctx context.Context) (bool, error)
	GetAccountBalance(ctx context.Context) (float64, error)
	GetWallet(ctx context.Context) (*models.Wallet, error)
}

// MexcWebSocketClient defines the interface for interacting with the MEXC WebSocket API
type MexcWebSocketClient interface {
	Connect(ctx context.Context) error
	Disconnect() error
	IsConnected() bool
	SubscribeToAccountUpdates(ctx context.Context, callback func(*models.Wallet)) error
	UnsubscribeFromAccountUpdates(ctx context.Context) error
	Authenticate(ctx context.Context) error
	SetReconnectHandler(handler func() error)
}

// realAccountService implements the AccountService interface with real MEXC API integration
type realAccountService struct {
	restClient      MexcRESTClient
	wsClient        MexcWebSocketClient
	coinRepo        BoughtCoinRepository
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	config          Config
	logger          *zap.Logger

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

// NewRealAccountService creates a new instance of the real account service
func NewRealAccountService(
	restClient MexcRESTClient,
	wsClient MexcWebSocketClient,
	coinRepo BoughtCoinRepository,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	config Config,
) AccountService {
	cacheTTL := config.GetCacheTTL()
	if cacheTTL <= 0 {
		cacheTTL = defaultCacheTTL
	}

	logger, _ := zap.NewProduction()

	return NewRealAccountServiceWithLogger(
		restClient,
		wsClient,
		coinRepo,
		walletRepo,
		transactionRepo,
		config,
		logger,
	)
}

// NewRealAccountServiceWithLogger creates a new instance of the real account service with a custom logger
func NewRealAccountServiceWithLogger(
	restClient MexcRESTClient,
	wsClient MexcWebSocketClient,
	coinRepo BoughtCoinRepository,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	config Config,
	logger *zap.Logger,
) AccountService {
	cacheTTL := config.GetCacheTTL()
	if cacheTTL <= 0 {
		cacheTTL = defaultCacheTTL
	}

	svc := &realAccountService{
		restClient:         restClient,
		wsClient:           wsClient,
		coinRepo:           coinRepo,
		walletRepo:         walletRepo,
		transactionRepo:    transactionRepo,
		config:             config,
		logger:             logger,
		cacheTTL:           cacheTTL,
		balanceCache:       nil,
		balanceCacheExp:    time.Time{},
		walletCache:        nil,
		walletCacheExp:     time.Time{},
		balanceSubscribers: make([]func(*models.Wallet), 0),
	}

	// Set up WebSocket reconnect handler
	if wsClient != nil {
		wsClient.SetReconnectHandler(func() error {
			return svc.handleReconnect()
		})
	}

	return svc
}

// GetAccountBalance retrieves the current account balance from the exchange
func (s *realAccountService) GetAccountBalance(ctx context.Context) (models.Balance, error) {
	// Check cache first for better performance
	s.mutex.RLock()
	if s.balanceCache != nil && time.Now().Before(s.balanceCacheExp) {
		balance := *s.balanceCache
		s.mutex.RUnlock()
		return balance, nil
	}
	s.mutex.RUnlock()

	// Get balance from exchange if cache is invalid
	s.logger.Debug("Fetching balances from MEXC API")
	balance, err := s.restClient.FetchBalances(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch balances from MEXC", zap.Error(err))
		return models.Balance{}, fmt.Errorf("failed to fetch balances from MEXC: %w", err)
	}

	// Update cache
	s.mutex.Lock()
	s.balanceCache = &balance
	s.balanceCacheExp = time.Now().Add(s.cacheTTL)
	s.mutex.Unlock()

	s.logger.Debug("Successfully fetched balances from MEXC",
		zap.Float64("fiat_balance", balance.Fiat),
		zap.Int("asset_count", len(balance.Available)))

	return balance, nil
}

// GetWallet retrieves the current wallet from the repository
func (s *realAccountService) GetWallet(ctx context.Context) (*models.Wallet, error) {
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
func (s *realAccountService) GetPortfolioValue(ctx context.Context) (float64, error) {
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
func (s *realAccountService) GetPositionRisk(ctx context.Context, symbol string) (models.PositionRisk, error) {
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
func (s *realAccountService) GetAllPositionRisks(ctx context.Context) (map[string]models.PositionRisk, error) {
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
func (s *realAccountService) ValidateAPIKeys(ctx context.Context) (bool, error) {
	return s.restClient.ValidateKeys(ctx)
}

// GetCurrentExposure calculates the total current exposure across all positions
func (s *realAccountService) GetCurrentExposure(ctx context.Context) (float64, error) {
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
func (s *realAccountService) UpdateBalance(ctx context.Context, amount float64, reason string) error {
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

	// Notify subscribers
	s.notifyBalanceSubscribers(updatedWallet)

	return nil
}

// SyncWithExchange synchronizes the local wallet balance with the exchange
func (s *realAccountService) SyncWithExchange(ctx context.Context) error {
	// Get balance from exchange
	balance, err := s.GetAccountBalance(ctx)
	if err != nil {
		return fmt.Errorf("error getting exchange balance: %w", err)
	}

	// Convert balance to wallet format
	wallet := &models.Wallet{
		Balances:  make(map[string]*models.AssetBalance),
		UpdatedAt: time.Now(),
	}

	// Add USDT balance
	wallet.Balances["USDT"] = &models.AssetBalance{
		Asset:  "USDT",
		Free:   balance.Fiat,
		Locked: 0,
		Total:  balance.Fiat,
	}

	// Add other assets
	for asset, amount := range balance.Available {
		wallet.Balances[asset] = &models.AssetBalance{
			Asset:  asset,
			Free:   amount,
			Locked: 0,
			Total:  amount,
		}
	}

	// Save to repository
	updatedWallet, err := s.walletRepo.SaveWallet(ctx, wallet)
	if err != nil {
		return fmt.Errorf("error saving wallet: %w", err)
	}

	// Update cache
	s.mutex.Lock()
	s.walletCache = updatedWallet
	s.walletCacheExp = time.Now().Add(s.cacheTTL)
	s.mutex.Unlock()

	// Notify subscribers
	s.notifyBalanceSubscribers(updatedWallet)

	return nil
}

// GetBalanceSummary generates a summary of the wallet and transactions
func (s *realAccountService) GetBalanceSummary(ctx context.Context, days int) (*models.BalanceSummary, error) {
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
func (s *realAccountService) GetTransactionHistory(ctx context.Context, startTime, endTime time.Time) ([]*models.Transaction, error) {
	if endTime.IsZero() {
		endTime = time.Now()
	}

	if startTime.IsZero() || startTime.After(endTime) {
		return nil, errors.New("invalid time range")
	}

	return s.transactionRepo.FindByTimeRange(ctx, startTime, endTime)
}

// AnalyzeTransactions performs analysis on transaction data
func (s *realAccountService) AnalyzeTransactions(ctx context.Context, startTime, endTime time.Time) (*models.TransactionAnalysis, error) {
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

// SubscribeToBalanceUpdates subscribes to real-time balance updates via WebSocket
func (s *realAccountService) SubscribeToBalanceUpdates(ctx context.Context, callback func(*models.Wallet)) error {
	// Add the callback to subscribers
	s.subMutex.Lock()
	s.balanceSubscribers = append(s.balanceSubscribers, callback)
	s.subMutex.Unlock()

	s.logger.Info("Subscribing to balance updates")

	// If WebSocket client is not connected, connect it
	if !s.wsClient.IsConnected() {
		s.logger.Debug("WebSocket not connected, connecting now")
		if err := s.wsClient.Connect(ctx); err != nil {
			s.logger.Error("Failed to connect WebSocket client", zap.Error(err))
			return fmt.Errorf("failed to connect WebSocket client: %w", err)
		}

		// Authenticate the WebSocket connection
		s.logger.Debug("Authenticating WebSocket connection")
		if err := s.wsClient.Authenticate(ctx); err != nil {
			s.logger.Error("Failed to authenticate WebSocket connection", zap.Error(err))
			return fmt.Errorf("failed to authenticate WebSocket connection: %w", err)
		}
	}

	// Subscribe to account updates
	s.logger.Debug("Subscribing to account updates via WebSocket")
	if err := s.wsClient.SubscribeToAccountUpdates(ctx, s.handleWalletUpdate); err != nil {
		s.logger.Error("Failed to subscribe to account updates", zap.Error(err))
		return fmt.Errorf("failed to subscribe to account updates: %w", err)
	}

	s.logger.Info("Successfully subscribed to balance updates")
	return nil
}

// handleWalletUpdate processes wallet updates from WebSocket
func (s *realAccountService) handleWalletUpdate(wallet *models.Wallet) {
	s.logger.Debug("Received wallet update from WebSocket",
		zap.Time("update_time", wallet.UpdatedAt),
		zap.Int("asset_count", len(wallet.Balances)))

	// Update cache
	s.mutex.Lock()
	s.walletCache = wallet
	s.walletCacheExp = time.Now().Add(s.cacheTTL)
	s.mutex.Unlock()

	// Save to repository
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := s.walletRepo.SaveWallet(ctx, wallet); err != nil {
		s.logger.Error("Failed to save wallet update", zap.Error(err))
	} else {
		s.logger.Debug("Successfully saved wallet update to repository")
	}

	// Notify subscribers
	s.notifyBalanceSubscribers(wallet)
}

// notifyBalanceSubscribers notifies all subscribers of a wallet update
func (s *realAccountService) notifyBalanceSubscribers(wallet *models.Wallet) {
	s.subMutex.RLock()
	subscribers := make([]func(*models.Wallet), len(s.balanceSubscribers))
	copy(subscribers, s.balanceSubscribers)
	subCount := len(subscribers)
	s.subMutex.RUnlock()

	s.logger.Debug("Notifying balance subscribers", zap.Int("subscriber_count", subCount))

	for _, callback := range subscribers {
		go callback(wallet)
	}
}

// handleReconnect handles WebSocket reconnection
func (s *realAccountService) handleReconnect() error {
	s.logger.Info("Handling WebSocket reconnection")

	// Create a new context for reconnection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to WebSocket
	s.logger.Debug("Reconnecting to WebSocket")
	if err := s.wsClient.Connect(ctx); err != nil {
		s.logger.Error("Failed to reconnect to WebSocket", zap.Error(err))
		return fmt.Errorf("failed to reconnect to WebSocket: %w", err)
	}

	// Authenticate the WebSocket connection
	s.logger.Debug("Re-authenticating WebSocket connection")
	if err := s.wsClient.Authenticate(ctx); err != nil {
		s.logger.Error("Failed to re-authenticate WebSocket connection", zap.Error(err))
		return fmt.Errorf("failed to re-authenticate WebSocket connection: %w", err)
	}

	// Re-subscribe to account updates
	s.logger.Debug("Re-subscribing to account updates")
	if err := s.wsClient.SubscribeToAccountUpdates(ctx, s.handleWalletUpdate); err != nil {
		s.logger.Error("Failed to re-subscribe to account updates", zap.Error(err))
		return fmt.Errorf("failed to re-subscribe to account updates: %w", err)
	}

	s.logger.Info("Successfully reconnected to WebSocket")
	return nil
}

// GetListenKey retrieves a listen key for WebSocket authentication
func (s *realAccountService) GetListenKey(ctx context.Context) (string, error) {
	// This would typically call an API endpoint to get a listen key
	// For now, we'll just return a placeholder error
	return "", errors.New("not implemented")
}

// RenewListenKey renews a listen key to keep it active
func (s *realAccountService) RenewListenKey(ctx context.Context, listenKey string) error {
	// This would typically call an API endpoint to renew a listen key
	// For now, we'll just return a placeholder error
	return errors.New("not implemented")
}

// CloseListenKey closes a listen key when it's no longer needed
func (s *realAccountService) CloseListenKey(ctx context.Context, listenKey string) error {
	// This would typically call an API endpoint to close a listen key
	// For now, we'll just return a placeholder error
	return errors.New("not implemented")
}

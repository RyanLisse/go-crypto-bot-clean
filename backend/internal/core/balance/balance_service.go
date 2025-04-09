package balance

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
	"go.uber.org/zap"
)

// Config defines the configuration for the balance service
type Config interface {
	GetSyncInterval() time.Duration
	GetCacheTTL() time.Duration
}

// ExchangeClient defines the interface for interacting with the exchange API
type ExchangeClient interface {
	FetchBalances(ctx context.Context) (models.Balance, error)
	GetWallet(ctx context.Context) (*models.Wallet, error)
}

// WebSocketClient defines the interface for receiving real-time updates
type WebSocketClient interface {
	SubscribeToAccountUpdates(ctx context.Context, callback func(*models.Wallet))
	IsConnected() bool
}

// BalanceUpdateCallback is a function that is called when a balance update is received
type BalanceUpdateCallback func(*models.Wallet)

// BalanceService manages account balances with real-time updates and fallback to stored data
type BalanceService interface {
	// GetLatestBalance returns the most up-to-date balance information
	GetLatestBalance(ctx context.Context) (*models.Wallet, error)
	
	// SyncWithExchange forces a sync with the exchange
	SyncWithExchange(ctx context.Context) error
	
	// SubscribeToUpdates registers a callback for balance updates
	SubscribeToUpdates(callback BalanceUpdateCallback)
	
	// Start starts the background sync process
	Start(ctx context.Context)
	
	// Stop stops the background sync process
	Stop()
}

type balanceService struct {
	exchangeClient  ExchangeClient
	wsClient        WebSocketClient
	walletRepo      repository.WalletRepository
	logger          *zap.Logger
	config          Config
	
	// Synchronization
	mutex           sync.RWMutex
	syncInterval    time.Duration
	cacheTTL        time.Duration
	
	// State
	latestWallet    *models.Wallet
	lastSyncTime    time.Time
	subscribers     []BalanceUpdateCallback
	
	// Control channels
	stopCh          chan struct{}
	syncCh          chan struct{}
}

// NewBalanceService creates a new balance service
func NewBalanceService(
	exchangeClient ExchangeClient,
	wsClient WebSocketClient,
	walletRepo repository.WalletRepository,
	logger *zap.Logger,
	config Config,
) BalanceService {
	syncInterval := config.GetSyncInterval()
	if syncInterval <= 0 {
		syncInterval = 5 * time.Minute // Default to 5 minutes
	}
	
	cacheTTL := config.GetCacheTTL()
	if cacheTTL <= 0 {
		cacheTTL = 1 * time.Minute // Default to 1 minute
	}
	
	return &balanceService{
		exchangeClient: exchangeClient,
		wsClient:       wsClient,
		walletRepo:     walletRepo,
		logger:         logger,
		config:         config,
		syncInterval:   syncInterval,
		cacheTTL:       cacheTTL,
		subscribers:    make([]BalanceUpdateCallback, 0),
		stopCh:         make(chan struct{}),
		syncCh:         make(chan struct{}, 1),
	}
}

// Start begins the background sync process
func (s *balanceService) Start(ctx context.Context) {
	// Initial sync
	if err := s.SyncWithExchange(ctx); err != nil {
		s.logger.Error("Failed initial balance sync", zap.Error(err))
	}
	
	// Subscribe to real-time updates if WebSocket is available
	if s.wsClient != nil {
		s.wsClient.SubscribeToAccountUpdates(ctx, s.handleRealTimeUpdate)
	}
	
	// Start background sync ticker
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Balance service stopping due to context cancellation")
			return
		case <-s.stopCh:
			s.logger.Info("Balance service stopping")
			return
		case <-ticker.C:
			if err := s.SyncWithExchange(ctx); err != nil {
				s.logger.Error("Failed periodic balance sync", zap.Error(err))
			}
		case <-s.syncCh:
			if err := s.SyncWithExchange(ctx); err != nil {
				s.logger.Error("Failed manual balance sync", zap.Error(err))
			}
		}
	}
}

// Stop halts the background sync process
func (s *balanceService) Stop() {
	close(s.stopCh)
}

// GetLatestBalance returns the most up-to-date balance information
func (s *balanceService) GetLatestBalance(ctx context.Context) (*models.Wallet, error) {
	s.mutex.RLock()
	
	// If we have a recent wallet and it's still valid, return it
	if s.latestWallet != nil && time.Since(s.lastSyncTime) < s.cacheTTL {
		wallet := s.latestWallet
		s.mutex.RUnlock()
		return wallet, nil
	}
	s.mutex.RUnlock()
	
	// Try to get real-time data if WebSocket is connected
	if s.wsClient != nil && s.wsClient.IsConnected() {
		// We'll still return the latest wallet we have, but trigger a sync
		// to update it asynchronously
		select {
		case s.syncCh <- struct{}{}:
			// Sync request queued
		default:
			// Channel full, sync already pending
		}
	}
	
	// If we have any wallet data, return it even if expired
	s.mutex.RLock()
	if s.latestWallet != nil {
		wallet := s.latestWallet
		s.mutex.RUnlock()
		return wallet, nil
	}
	s.mutex.RUnlock()
	
	// No cached data, sync now
	if err := s.SyncWithExchange(ctx); err != nil {
		// Try to get from repository as last resort
		wallet, err := s.walletRepo.GetWallet(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get wallet: %w", err)
		}
		return wallet, nil
	}
	
	s.mutex.RLock()
	wallet := s.latestWallet
	s.mutex.RUnlock()
	
	return wallet, nil
}

// SyncWithExchange forces a sync with the exchange
func (s *balanceService) SyncWithExchange(ctx context.Context) error {
	// Get wallet from exchange
	wallet, err := s.exchangeClient.GetWallet(ctx)
	if err != nil {
		return fmt.Errorf("failed to get wallet from exchange: %w", err)
	}
	
	// Update last sync time
	now := time.Now()
	
	// Save to repository
	savedWallet, err := s.walletRepo.SaveWallet(ctx, wallet)
	if err != nil {
		s.logger.Error("Failed to save wallet to repository", zap.Error(err))
		// Continue with the wallet we got from the exchange
	} else {
		wallet = savedWallet
	}
	
	// Update cache
	s.mutex.Lock()
	s.latestWallet = wallet
	s.lastSyncTime = now
	s.mutex.Unlock()
	
	// Notify subscribers
	s.notifySubscribers(wallet)
	
	return nil
}

// SubscribeToUpdates registers a callback for balance updates
func (s *balanceService) SubscribeToUpdates(callback BalanceUpdateCallback) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.subscribers = append(s.subscribers, callback)
}

// handleRealTimeUpdate processes real-time balance updates from WebSocket
func (s *balanceService) handleRealTimeUpdate(wallet *models.Wallet) {
	s.mutex.Lock()
	
	// Update our cached wallet with the real-time data
	s.latestWallet = wallet
	s.lastSyncTime = time.Now()
	
	// Make a copy for subscribers
	walletCopy := wallet
	
	s.mutex.Unlock()
	
	// Save to repository in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if _, err := s.walletRepo.SaveWallet(ctx, wallet); err != nil {
			s.logger.Error("Failed to save real-time wallet update", zap.Error(err))
		}
	}()
	
	// Notify subscribers
	s.notifySubscribers(walletCopy)
}

// notifySubscribers sends updates to all registered callbacks
func (s *balanceService) notifySubscribers(wallet *models.Wallet) {
	s.mutex.RLock()
	subscribers := make([]BalanceUpdateCallback, len(s.subscribers))
	copy(subscribers, s.subscribers)
	s.mutex.RUnlock()
	
	for _, callback := range subscribers {
		go callback(wallet)
	}
}

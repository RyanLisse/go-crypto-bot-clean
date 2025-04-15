package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// WalletDataSyncService defines the interface for wallet data synchronization
type WalletDataSyncService interface {
	// SyncWallet synchronizes wallet data for a specific wallet
	SyncWallet(ctx context.Context, walletID string) (*model.Wallet, error)

	// SyncWalletsByUserID synchronizes wallet data for all wallets of a user
	SyncWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error)

	// SyncAllWallets synchronizes wallet data for all wallets
	SyncAllWallets(ctx context.Context) (int, error)

	// ScheduleWalletSync schedules a wallet sync for a specific wallet
	ScheduleWalletSync(ctx context.Context, walletID string, interval time.Duration) error

	// CancelWalletSync cancels a scheduled wallet sync
	CancelWalletSync(ctx context.Context, walletID string) error

	// GetLastSyncTime gets the last sync time for a wallet
	GetLastSyncTime(ctx context.Context, walletID string) (*time.Time, error)

	// GetSyncStatus gets the sync status for a wallet
	GetSyncStatus(ctx context.Context, walletID string) (model.SyncStatus, error)

	// SaveBalanceHistory saves the balance history for a wallet
	SaveBalanceHistory(ctx context.Context, walletID string) error
}

// walletDataSyncService implements WalletDataSyncService
type walletDataSyncService struct {
	walletRepo           port.WalletRepository
	apiCredentialManager APICredentialManagerService
	providerRegistry     *wallet.ProviderRegistry
	logger               *zerolog.Logger
	syncJobs             map[string]*syncJob
	mu                   sync.RWMutex
}

type syncJob struct {
	walletID  string
	ticker    *time.Ticker
	done      chan bool
	lastSync  time.Time
	status    model.SyncStatus
	interval  time.Duration
	isRunning bool
}

// NewWalletDataSyncService creates a new wallet data sync service
func NewWalletDataSyncService(
	walletRepo port.WalletRepository,
	apiCredentialManager APICredentialManagerService,
	providerRegistry *wallet.ProviderRegistry,
	logger *zerolog.Logger,
) WalletDataSyncService {
	return &walletDataSyncService{
		walletRepo:           walletRepo,
		apiCredentialManager: apiCredentialManager,
		providerRegistry:     providerRegistry,
		logger:               logger,
		syncJobs:             make(map[string]*syncJob),
	}
}

// SyncWallet synchronizes wallet data for a specific wallet
func (s *walletDataSyncService) SyncWallet(ctx context.Context, walletID string) (*model.Wallet, error) {
	if walletID == "" {
		return nil, errors.New("wallet ID is required")
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to get wallet")
		return nil, err
	}

	// Update sync status
	s.updateSyncStatus(walletID, model.SyncStatusInProgress)

	// Sync wallet based on type
	var syncedWallet *model.Wallet
	switch wallet.Type {
	case model.WalletTypeExchange:
		syncedWallet, err = s.syncExchangeWallet(ctx, wallet)
	case model.WalletTypeWeb3:
		syncedWallet, err = s.syncWeb3Wallet(ctx, wallet)
	default:
		err = fmt.Errorf("unsupported wallet type: %s", wallet.Type)
	}

	if err != nil {
		s.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to sync wallet")
		s.updateSyncStatus(walletID, model.SyncStatusFailed)
		return nil, err
	}

	// Update last sync time
	now := time.Now()
	syncedWallet.LastSynced = &now
	syncedWallet.LastSyncAt = now

	// Save wallet
	if err := s.walletRepo.Save(ctx, syncedWallet); err != nil {
		s.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to save wallet after sync")
		s.updateSyncStatus(walletID, model.SyncStatusFailed)
		return nil, err
	}

	// Save balance history
	if err := s.SaveBalanceHistory(ctx, walletID); err != nil {
		s.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to save balance history")
		// Continue anyway, this is not critical
	}

	s.updateSyncStatus(walletID, model.SyncStatusSuccess)
	s.logger.Info().Str("walletID", walletID).Msg("Wallet synced successfully")
	return syncedWallet, nil
}

// SyncWalletsByUserID synchronizes wallet data for all wallets of a user
func (s *walletDataSyncService) SyncWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get wallets
	wallets, err := s.walletRepo.GetWalletsByUserID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallets")
		return nil, err
	}

	if len(wallets) == 0 {
		return nil, nil
	}

	// Sync wallets concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex
	syncedWallets := make([]*model.Wallet, 0, len(wallets))
	errCount := 0

	for _, wallet := range wallets {
		wg.Add(1)
		go func(w *model.Wallet) {
			defer wg.Done()

			syncedWallet, err := s.SyncWallet(ctx, w.ID)
			if err != nil {
				s.logger.Error().Err(err).Str("walletID", w.ID).Msg("Failed to sync wallet")
				mu.Lock()
				errCount++
				mu.Unlock()
				return
			}

			mu.Lock()
			syncedWallets = append(syncedWallets, syncedWallet)
			mu.Unlock()
		}(wallet)
	}

	wg.Wait()

	if errCount == len(wallets) {
		return nil, fmt.Errorf("failed to sync all wallets for user %s", userID)
	}

	s.logger.Info().Str("userID", userID).Int("total", len(wallets)).Int("synced", len(syncedWallets)).Int("failed", errCount).Msg("Wallets synced")
	return syncedWallets, nil
}

// SyncAllWallets synchronizes wallet data for all wallets
func (s *walletDataSyncService) SyncAllWallets(ctx context.Context) (int, error) {
	// Get all users with wallets
	// This is a simplified implementation that doesn't handle pagination
	// In a real implementation, we would need to get all users and then get their wallets

	// For now, we'll just return a placeholder
	s.logger.Warn().Msg("SyncAllWallets is not fully implemented yet")
	return 0, nil

	// This would be the implementation using existing methods
	/*
		// Get all users (this would require a UserRepository)
		// For each user, get their wallets
		// Collect all wallets
		var allWallets []*model.Wallet

		// For each user
		wallets, err := s.walletRepo.GetWalletsByUserID(ctx, userID)
		if err != nil {
			s.logger.Error().Err(err).Str("userID", userID).Msg("Failed to get wallets for user")
			continue
		}

		allWallets = append(allWallets, wallets...)

		if len(allWallets) == 0 {
			return 0, nil
		}

		// Sync wallets concurrently with rate limiting
		var wg sync.WaitGroup
		var mu sync.Mutex
		syncedCount := 0
		errCount := 0
		semaphore := make(chan struct{}, 10) // Limit to 10 concurrent syncs

		for _, wallet := range allWallets {
			wg.Add(1)
			semaphore <- struct{}{} // Acquire semaphore
			go func(w *model.Wallet) {
				defer func() {
					<-semaphore // Release semaphore
					wg.Done()
				}()

				_, err := s.SyncWallet(ctx, w.ID)
				if err != nil {
					s.logger.Error().Err(err).Str("walletID", w.ID).Msg("Failed to sync wallet")
					mu.Lock()
					errCount++
					mu.Unlock()
					return
				}

				mu.Lock()
				syncedCount++
				mu.Unlock()
			}(wallet)
		}

		wg.Wait()

		s.logger.Info().Int("total", len(allWallets)).Int("synced", syncedCount).Int("failed", errCount).Msg("All wallets synced")
		return syncedCount, nil
	*/
}

// ScheduleWalletSync schedules a wallet sync for a specific wallet
func (s *walletDataSyncService) ScheduleWalletSync(ctx context.Context, walletID string, interval time.Duration) error {
	if walletID == "" {
		return errors.New("wallet ID is required")
	}

	if interval < time.Minute {
		return errors.New("interval must be at least 1 minute")
	}

	// Check if wallet exists
	_, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to get wallet")
		return err
	}

	// Cancel existing job if any
	s.CancelWalletSync(ctx, walletID)

	// Create new job
	ticker := time.NewTicker(interval)
	done := make(chan bool)
	job := &syncJob{
		walletID:  walletID,
		ticker:    ticker,
		done:      done,
		lastSync:  time.Now(),
		status:    model.SyncStatusScheduled,
		interval:  interval,
		isRunning: true,
	}

	s.mu.Lock()
	s.syncJobs[walletID] = job
	s.mu.Unlock()

	// Start sync job
	go func() {
		for {
			select {
			case <-ticker.C:
				_, err := s.SyncWallet(ctx, walletID)
				if err != nil {
					s.logger.Error().Err(err).Str("walletID", walletID).Msg("Scheduled sync failed")
				}
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	s.logger.Info().Str("walletID", walletID).Dur("interval", interval).Msg("Wallet sync scheduled")
	return nil
}

// CancelWalletSync cancels a scheduled wallet sync
func (s *walletDataSyncService) CancelWalletSync(ctx context.Context, walletID string) error {
	if walletID == "" {
		return errors.New("wallet ID is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.syncJobs[walletID]
	if !exists {
		return nil
	}

	if job.isRunning {
		job.done <- true
		job.isRunning = false
	}

	delete(s.syncJobs, walletID)
	s.logger.Info().Str("walletID", walletID).Msg("Wallet sync cancelled")
	return nil
}

// GetLastSyncTime gets the last sync time for a wallet
func (s *walletDataSyncService) GetLastSyncTime(ctx context.Context, walletID string) (*time.Time, error) {
	if walletID == "" {
		return nil, errors.New("wallet ID is required")
	}

	// Check if job exists
	s.mu.RLock()
	job, exists := s.syncJobs[walletID]
	s.mu.RUnlock()

	if exists {
		lastSync := job.lastSync
		return &lastSync, nil
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to get wallet")
		return nil, err
	}

	return wallet.LastSynced, nil
}

// GetSyncStatus gets the sync status for a wallet
func (s *walletDataSyncService) GetSyncStatus(ctx context.Context, walletID string) (model.SyncStatus, error) {
	if walletID == "" {
		return "", errors.New("wallet ID is required")
	}

	// Check if job exists
	s.mu.RLock()
	job, exists := s.syncJobs[walletID]
	s.mu.RUnlock()

	if exists {
		return job.status, nil
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to get wallet")
		return "", err
	}

	return wallet.SyncStatus, nil
}

// SaveBalanceHistory saves the balance history for a wallet
func (s *walletDataSyncService) SaveBalanceHistory(ctx context.Context, walletID string) error {
	if walletID == "" {
		return errors.New("wallet ID is required")
	}

	// Get wallet
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		s.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to get wallet")
		return err
	}

	// Create balance history
	now := time.Now()

	history := &model.BalanceHistory{
		ID:            model.GenerateID(),
		UserID:        wallet.UserID,
		WalletID:      walletID,
		Balances:      wallet.Balances,
		TotalUSDValue: wallet.TotalUSDValue,
		Timestamp:     now,
	}

	// Save balance history
	if err := s.walletRepo.SaveBalanceHistory(ctx, history); err != nil {
		s.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to save balance history")
		return err
	}

	s.logger.Info().Str("walletID", walletID).Msg("Balance history saved")
	return nil
}

// syncExchangeWallet synchronizes an exchange wallet
func (s *walletDataSyncService) syncExchangeWallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	// Get API credentials
	credential, err := s.apiCredentialManager.GetCredentialForExchange(ctx, wallet.UserID, wallet.Exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to get API credentials: %w", err)
	}

	// Mark credential as used
	if err := s.apiCredentialManager.MarkCredentialAsUsed(ctx, credential.ID); err != nil {
		s.logger.Warn().Err(err).Str("credentialID", credential.ID).Msg("Failed to mark credential as used")
		// Continue anyway, this is not critical
	}

	// Get exchange provider
	provider, err := s.providerRegistry.GetExchangeProvider(wallet.Exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange provider: %w", err)
	}

	// Set API credentials
	if err := provider.SetAPICredentials(ctx, credential.APIKey, credential.APISecret); err != nil {
		return nil, fmt.Errorf("failed to set API credentials: %w", err)
	}

	// Get balance
	syncedWallet, err := provider.GetBalance(ctx, wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return syncedWallet, nil
}

// syncWeb3Wallet synchronizes a Web3 wallet
func (s *walletDataSyncService) syncWeb3Wallet(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	// Get Web3 provider
	provider, err := s.providerRegistry.GetWeb3Provider(wallet.Network)
	if err != nil {
		return nil, fmt.Errorf("failed to get Web3 provider: %w", err)
	}

	// Get balance
	syncedWallet, err := provider.GetBalance(ctx, wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return syncedWallet, nil
}

// updateSyncStatus updates the sync status for a wallet
func (s *walletDataSyncService) updateSyncStatus(walletID string, status model.SyncStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.syncJobs[walletID]
	if exists {
		job.status = status
		job.lastSync = time.Now()
	}
}

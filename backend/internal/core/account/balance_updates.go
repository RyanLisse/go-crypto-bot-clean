package account

import (
	"context"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// SubscribeToBalanceUpdates registers a callback for balance updates
func (s *accountService) SubscribeToBalanceUpdates(ctx context.Context, callback func(*models.Wallet)) error {
	s.subMutex.Lock()
	defer s.subMutex.Unlock()

	s.balanceSubscribers = append(s.balanceSubscribers, callback)

	// If balance service is available, subscribe there too
	if s.balanceService != nil {
		s.balanceService.SubscribeToUpdates(callback)
	}

	return nil
}

// handleBalanceUpdate processes balance updates from the balance service
func (s *accountService) handleBalanceUpdate(wallet *models.Wallet) {
	// Update cache
	s.mutex.Lock()
	s.walletCache = wallet
	s.walletCacheExp = wallet.UpdatedAt.Add(s.cacheTTL)
	s.mutex.Unlock()

	// Notify subscribers
	s.notifyBalanceSubscribers(wallet)
}

// notifyBalanceSubscribers sends updates to all registered callbacks
func (s *accountService) notifyBalanceSubscribers(wallet *models.Wallet) {
	s.subMutex.RLock()
	subscribers := make([]func(*models.Wallet), len(s.balanceSubscribers))
	copy(subscribers, s.balanceSubscribers)
	s.subMutex.RUnlock()

	for _, callback := range subscribers {
		go callback(wallet)
	}
}

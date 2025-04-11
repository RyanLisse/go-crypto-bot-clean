package account

import (
	"context"

	"go-crypto-bot-clean/backend/internal/domain/models"
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
// This method is called by the balance service when a wallet update is received
// It's currently not directly called but is used as a callback handler
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
// This method is called by handleBalanceUpdate to notify all subscribers of wallet changes
// It's part of the internal notification mechanism for balance updates
// NOTE: This function may be flagged as unused by static analysis tools because it's only
// called from handleBalanceUpdate which itself might be called through event-driven mechanisms.
// Do not remove this function as it's essential for the balance update notification system.
func (s *accountService) notifyBalanceSubscribers(wallet *models.Wallet) {
	s.subMutex.RLock()
	subscribers := make([]func(*models.Wallet), len(s.balanceSubscribers))
	copy(subscribers, s.balanceSubscribers)
	s.subMutex.RUnlock()

	for _, callback := range subscribers {
		go callback(wallet)
	}
}

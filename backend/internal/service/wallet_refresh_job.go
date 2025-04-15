package service

import (
	"context"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// WalletRefreshJob periodically refreshes wallet data from exchanges
// and persists new balance history snapshots

type WalletRefreshJob struct {
	interval    time.Duration
	stopCh      chan struct{}
	wg          sync.WaitGroup
	cache       port.WalletCache
	providerReg *wallet.ProviderRegistry
	historyRepo port.BalanceHistoryRepository
	logger      *zerolog.Logger
}

func NewWalletRefreshJob(
	interval time.Duration,
	cache port.WalletCache,
	providerReg *wallet.ProviderRegistry,
	historyRepo port.BalanceHistoryRepository,
	logger *zerolog.Logger,
) *WalletRefreshJob {
	return &WalletRefreshJob{
		interval:    interval,
		stopCh:      make(chan struct{}),
		cache:       cache,
		providerReg: providerReg,
		historyRepo: historyRepo,
		logger:      logger,
	}
}

func (j *WalletRefreshJob) Start(ctx context.Context, walletIDs []string) {
	j.logger.Info().Msg("Starting WalletRefreshJob")
	j.wg.Add(1)
	go func() {
		defer j.wg.Done()
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				j.logger.Info().Msg("WalletRefreshJob stopped by context")
				return
			case <-j.stopCh:
				j.logger.Info().Msg("WalletRefreshJob stopped by stopCh")
				return
			case <-ticker.C:
				j.refreshAll(ctx, walletIDs)
			}
		}
	}()
}

func (j *WalletRefreshJob) Stop() {
	close(j.stopCh)
	j.wg.Wait()
}

func (j *WalletRefreshJob) refreshAll(ctx context.Context, walletIDs []string) {
	for _, walletID := range walletIDs {
		wallet, found := j.cache.Get(walletID)
		if !found || wallet == nil {
			j.logger.Warn().Str("walletID", walletID).Msg("Wallet not found in cache, skipping")
			continue
		}
		provider, err := j.providerReg.GetProvider(wallet.Exchange)
		if err != nil {
			j.logger.Error().Err(err).Str("exchange", wallet.Exchange).Msg("Provider not found")
			continue
		}
		// Only refresh if provider supports ExchangeWalletProvider
		exProvider, ok := provider.(port.ExchangeWalletProvider)
		if !ok {
			j.logger.Warn().Str("exchange", wallet.Exchange).Msg("Provider does not support ExchangeWalletProvider")
			continue
		}
		// Fetch updated wallet data
		updatedWallet, err := exProvider.GetBalance(ctx, wallet)
		if err != nil {
			j.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to refresh wallet balance")
			continue
		}
		j.cache.Set(walletID, updatedWallet)
		// Persist balance history
		history := &model.BalanceHistory{
			ID:            model.GenerateID(),
			UserID:        updatedWallet.UserID,
			WalletID:      updatedWallet.ID,
			Balances:      updatedWallet.Balances,
			TotalUSDValue: updatedWallet.TotalUSDValue,
			Timestamp:     time.Now(),
		}
		if err := j.historyRepo.Save(ctx, history); err != nil {
			j.logger.Error().Err(err).Str("walletID", walletID).Msg("Failed to save balance history")
		}
		j.logger.Info().Str("walletID", walletID).Msg("Wallet refreshed and history saved")
	}
}

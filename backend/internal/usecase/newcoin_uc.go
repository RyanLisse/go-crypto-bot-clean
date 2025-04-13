package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/neo/crypto-bot/internal/domain/event"
	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog/log" // Assuming zerolog is used based on go.mod
)

// SymbolInfo is a struct for symbol information.
// TODO: Ensure this is properly defined in model or adjust as needed.
type SymbolInfo struct {
	Symbol string
	Status string // e.g., "TRADING", "AUCTION", "BREAK", etc. from the exchange API
}

// MarketDataServiceProvider defines the interface required from the market data service.
// This makes the dependency explicit for this use case.
// TODO: Define this properly, potentially reusing/refining existing MarketDataService port.
type MarketDataServiceProvider interface {
	GetSymbolInfo(ctx context.Context, symbol string) (*SymbolInfo, error) // Using local SymbolInfo
	// Add other methods if needed, e.g., GetTicker
}

// NewCoinUsecase handles the logic for detecting and processing new coin listings.
type NewCoinUsecase struct {
	repo      port.NewCoinRepository
	bus       port.EventBus
	marketSvc MarketDataServiceProvider // Use the specific interface
}

// NewNewCoinUsecase creates a new instance of NewCoinUsecase.
func NewNewCoinUsecase(
	repo port.NewCoinRepository,
	bus port.EventBus,
	marketSvc MarketDataServiceProvider, // Use the specific interface
) *NewCoinUsecase {
	return &NewCoinUsecase{
		repo:      repo,
		bus:       bus,
		marketSvc: marketSvc,
	}
}

// CheckNewListings polls for new coins, checks their status, and publishes events when they become tradable.
func (uc *NewCoinUsecase) CheckNewListings(ctx context.Context) error {
	// For now, just check coins listed/expected recently based on the test setup.
	// A real implementation would need more sophisticated time window logic.
	threshold := time.Now().Add(-24 * time.Hour) // Arbitrary threshold for FindRecentlyListed
	coinsToCheck, err := uc.repo.FindRecentlyListed(ctx, threshold)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find recently listed coins")
		return fmt.Errorf("failed to find recently listed coins: %w", err)
	}

	log.Info().Int("count", len(coinsToCheck)).Msg("Checking status for potential new coins")

	for _, coin := range coinsToCheck {
		// Only check coins that are expected but not yet trading
		if coin.Status == model.StatusExpected {
			log.Info().Str("symbol", coin.Symbol).Msg("Checking status for expected coin")
			symbolInfo, err := uc.marketSvc.GetSymbolInfo(ctx, coin.Symbol)
			if err != nil {
				// Log error but continue checking other coins
				log.Error().Err(err).Str("symbol", coin.Symbol).Msg("Failed to get symbol info from market service")
				continue
			}

			// Check if the status from the exchange indicates it's now tradable.
			// This comparison might need adjustment based on actual API response values.
			// Assuming "TRADING" status string for now, matching the test mock.
			if symbolInfo.Status == "TRADING" {
				log.Info().Str("symbol", coin.Symbol).Msg("Coin detected as TRADING")
				now := time.Now().UTC()
				coin.MarkAsTradable(now) // Update coin state

				// Persist the status change
				err = uc.repo.Update(ctx, coin)
				if err != nil {
					log.Error().Err(err).Str("symbol", coin.Symbol).Msg("Failed to update coin status in repository")
					// Decide on error handling: continue, return, retry? For now, continue.
					continue
				}
				log.Info().Str("symbol", coin.Symbol).Msg("Coin status updated in repository")

				// Publish the domain event
				// TODO: Get actual price/volume if available from GetSymbolInfo or another call
				tradableEvent := event.NewNewCoinTradable(coin, nil, nil) // Pass nil for price/volume for now
				err = uc.bus.Publish(ctx, tradableEvent)
				if err != nil {
					log.Error().Err(err).Str("symbol", coin.Symbol).Msg("Failed to publish NewCoinTradable event")
					// If publishing fails, should we revert the DB update? Depends on transaction strategy.
					// For now, log and continue.
					continue
				}
				log.Info().Str("symbol", coin.Symbol).Msg("NewCoinTradable event published")

				// Optionally, mark as processed immediately or let another process handle it
				// coin.MarkAsProcessed(now)
				// uc.repo.Update(ctx, coin) // Update again if marking processed here
			}
		}
	}

	return nil
}

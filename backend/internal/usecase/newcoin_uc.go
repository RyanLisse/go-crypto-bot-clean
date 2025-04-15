package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// NewCoinUseCase implements the new coin detection and management logic
type NewCoinUseCase struct {
	repo      port.NewCoinRepository
	eventRepo port.EventRepository
	eventBus  port.EventBus // Added EventBus
	mexc      port.MEXCClient
	logger    *zerolog.Logger
}

// NewNewCoinUseCase creates a new NewCoinUseCase instance
func NewNewCoinUseCase(repo port.NewCoinRepository, eventRepo port.EventRepository, eventBus port.EventBus, mexc port.MEXCClient, logger *zerolog.Logger) *NewCoinUseCase {
	return &NewCoinUseCase{
		repo:      repo,
		eventRepo: eventRepo,
		eventBus:  eventBus, // Initialize EventBus
		mexc:      mexc,
		logger:    logger,
	}
}

// DetectNewCoins checks for newly listed coins on MEXC
func (uc *NewCoinUseCase) DetectNewCoins() error {
	ctx := context.Background()
	coins, err := uc.mexc.GetNewListings(ctx)
	if err != nil {
		return fmt.Errorf("failed to get new listings: %w", err)
	}

	for _, coin := range coins {
		existing, err := uc.repo.GetBySymbol(ctx, coin.Symbol)
		if err != nil {
			uc.logger.Error().Err(err).Str("symbol", coin.Symbol).Msg("Failed to check existing coin")
			continue
		}

		if existing == nil {
			// New coin found
			coin.ID = uuid.New().String()
			if err := uc.repo.Save(ctx, coin); err != nil {
				uc.logger.Error().Err(err).Str("symbol", coin.Symbol).Msg("Failed to create new coin")
				continue
			}

			// Create and emit event
			event := &model.NewCoinEvent{
				ID:        uuid.New().String(),
				CoinID:    coin.ID,
				EventType: "new_coin_detected",
				NewStatus: coin.Status,
				CreatedAt: time.Now(),
			}
			if err := uc.eventRepo.SaveEvent(ctx, event); err != nil { // Use eventRepo.SaveEvent
				uc.logger.Error().Err(err).Str("symbol", coin.Symbol).Msg("Failed to save event")
			}
			uc.eventBus.Publish(event) // Publish event via EventBus
		} else {
			// Update existing coin if status changed
			if existing.Status != coin.Status {
				oldStatus := existing.Status
				existing.Status = coin.Status
				existing.UpdatedAt = time.Now()

				if coin.Status == model.StatusTrading && existing.BecameTradableAt == nil {
					now := time.Now()
					existing.BecameTradableAt = &now
				}

				if err := uc.repo.Update(ctx, existing); err != nil {
					uc.logger.Error().Err(err).Str("symbol", coin.Symbol).Msg("Failed to update coin")
					continue
				}

				// Create and emit status change event
				event := &model.NewCoinEvent{
					ID:        uuid.New().String(),
					CoinID:    existing.ID,
					EventType: "status_changed",
					OldStatus: oldStatus,
					NewStatus: coin.Status,
					CreatedAt: time.Now(),
				}
				if err := uc.eventRepo.SaveEvent(ctx, event); err != nil { // Use eventRepo.SaveEvent
					uc.logger.Error().Err(err).Str("symbol", coin.Symbol).Msg("Failed to save event")
				}
				uc.eventBus.Publish(event) // Publish event via EventBus
			}
		}
	}

	return nil
}

// UpdateCoinStatus updates a coin's status and creates an event
func (uc *NewCoinUseCase) UpdateCoinStatus(coinID string, newStatus model.Status) error {
	ctx := context.Background()
	coin, err := uc.repo.GetBySymbol(ctx, coinID)
	if err != nil {
		return fmt.Errorf("failed to get coin: %w", err)
	}
	if coin == nil {
		return fmt.Errorf("coin not found: %s", coinID)
	}

	oldStatus := coin.Status
	coin.Status = newStatus
	coin.UpdatedAt = time.Now()

	if newStatus == model.StatusTrading && coin.BecameTradableAt == nil {
		now := time.Now()
		coin.BecameTradableAt = &now
	}

	if err := uc.repo.Update(ctx, coin); err != nil {
		return fmt.Errorf("failed to update coin: %w", err)
	}

	event := &model.NewCoinEvent{
		ID:        uuid.New().String(),
		CoinID:    coin.ID,
		EventType: "status_changed",
		OldStatus: oldStatus,
		NewStatus: newStatus,
		CreatedAt: time.Now(),
	}
	if err := uc.eventRepo.SaveEvent(ctx, event); err != nil { // Use eventRepo.SaveEvent
		return fmt.Errorf("failed to save event: %w", err)
	}
	uc.eventBus.Publish(event) // Publish event via EventBus

	return nil
}

// GetCoinDetails retrieves detailed information about a coin
func (uc *NewCoinUseCase) GetCoinDetails(symbol string) (*model.NewCoin, error) {
	ctx := context.Background()
	return uc.repo.GetBySymbol(ctx, symbol)
}

// ListNewCoins retrieves a list of new coins with optional filtering
func (uc *NewCoinUseCase) ListNewCoins(status model.Status, limit, offset int) ([]*model.NewCoin, error) {
	ctx := context.Background()
	return uc.repo.GetByStatus(ctx, status)
}

// GetRecentTradableCoins retrieves recently listed coins that are now tradable
func (uc *NewCoinUseCase) GetRecentTradableCoins(limit int) ([]*model.NewCoin, error) {
	ctx := context.Background()
	return uc.repo.GetRecent(ctx, limit)
}

// SubscribeToEvents allows subscribing to new coin events
func (uc *NewCoinUseCase) SubscribeToEvents(callback func(*model.NewCoinEvent)) error {
	uc.eventBus.Subscribe(callback)
	return nil
}

// UnsubscribeFromEvents removes an event subscription
func (uc *NewCoinUseCase) UnsubscribeFromEvents(callback func(*model.NewCoinEvent)) error {
	uc.eventBus.Unsubscribe(callback)
	return nil
}

// Note: emitEvent is removed as publishing is handled by the eventBus directly.

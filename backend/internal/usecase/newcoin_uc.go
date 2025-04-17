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

// NewCoinUseCaseImpl implements the NewCoinUseCase interface
type NewCoinUseCaseImpl struct {
	repo      port.NewCoinRepository
	eventRepo port.EventRepository
	eventBus  port.EventBus // Added EventBus
	mexc      port.MEXCClient
	logger    *zerolog.Logger
}

// NewNewCoinUseCase creates a new NewCoinUseCase instance
func NewNewCoinUseCase(repo port.NewCoinRepository, eventRepo port.EventRepository, eventBus port.EventBus, mexc port.MEXCClient, logger *zerolog.Logger) NewCoinUseCase {
	return &NewCoinUseCaseImpl{
		repo:      repo,
		eventRepo: eventRepo,
		eventBus:  eventBus, // Initialize EventBus
		mexc:      mexc,
		logger:    logger,
	}
}

// DetectNewCoins checks for newly listed coins on MEXC
func (uc *NewCoinUseCaseImpl) DetectNewCoins() error {
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

				if coin.Status == model.CoinStatusTrading && existing.BecameTradableAt == nil {
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
func (uc *NewCoinUseCaseImpl) UpdateCoinStatus(coinID string, newStatus model.CoinStatus) error {
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

	if newStatus == model.CoinStatusTrading && coin.BecameTradableAt == nil {
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
func (uc *NewCoinUseCaseImpl) GetCoinDetails(symbol string) (*model.Coin, error) {
	// This is a temporary implementation that converts NewCoin to Coin
	ctx := context.Background()
	newCoin, err := uc.repo.GetBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	if newCoin == nil {
		return nil, nil
	}

	// Convert NewCoin to Coin
	return &model.Coin{
		ID:          newCoin.ID,
		Symbol:      newCoin.Symbol,
		Name:        newCoin.Name,
		Description: "",
		Status:      model.Status(string(newCoin.Status)),
		ListedAt:    newCoin.CreatedAt,
		UpdatedAt:   newCoin.UpdatedAt,
	}, nil
}

// ListNewCoins retrieves a list of new coins with optional filtering
func (uc *NewCoinUseCaseImpl) ListNewCoins(status model.CoinStatus, limit, offset int) ([]*model.NewCoin, error) {
	ctx := context.Background()
	// Retrieve coins by status
	return uc.repo.GetByStatus(ctx, status)
}

// GetRecentTradableCoins retrieves recently listed coins that are now tradable
func (uc *NewCoinUseCaseImpl) GetRecentTradableCoins(limit int) ([]*model.NewCoin, error) {
	ctx := context.Background()
	return uc.repo.GetRecent(ctx, limit)
}

// SubscribeToEvents allows subscribing to new coin events
func (uc *NewCoinUseCaseImpl) SubscribeToEvents(callback func(*model.CoinEvent)) error {
	// Convert CoinEvent callback to NewCoinEvent callback
	wrapper := func(event *model.NewCoinEvent) {
		// Convert NewCoinEvent to CoinEvent
		coinEvent := &model.CoinEvent{
			CoinID:     event.CoinID,
			EventType:  event.EventType,
			OldStatus:  model.CoinStatus(string(event.OldStatus)),
			NewStatus:  model.CoinStatus(string(event.NewStatus)),
			Timestamp:  event.CreatedAt,
			Exchange:   "mexc",
			Additional: make(map[string]interface{}),
		}
		callback(coinEvent)
	}
	uc.eventBus.Subscribe(wrapper)
	return nil
}

// UnsubscribeFromEvents removes an event subscription
func (uc *NewCoinUseCaseImpl) UnsubscribeFromEvents(callback func(*model.CoinEvent)) error {
	// This is a simplified implementation that doesn't actually unsubscribe
	// In a real implementation, we would need to keep track of the wrapper functions
	return nil
}

// Note: emitEvent is removed as publishing is handled by the eventBus directly.

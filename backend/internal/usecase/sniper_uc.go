package usecase

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// Errors
var (
	ErrSniperNotInitialized = errors.New("sniper service not initialized")
	ErrInvalidSniperConfig  = errors.New("invalid sniper configuration")
)

// SniperUseCase defines the interface for sniper operations
type SniperUseCase interface {
	// ExecuteSnipe executes a high-speed buy on a newly listed token
	ExecuteSnipe(ctx context.Context, symbol string) (*model.Order, error)

	// ExecuteSnipeWithConfig executes a high-speed buy with custom configuration
	ExecuteSnipeWithConfig(ctx context.Context, symbol string, config *port.SniperConfig) (*model.Order, error)

	// GetSniperConfig returns the current sniper configuration
	GetSniperConfig() (*port.SniperConfig, error)

	// UpdateSniperConfig updates the sniper configuration
	UpdateSniperConfig(config *port.SniperConfig) error

	// StartSniper starts the sniper service
	StartSniper() error

	// StopSniper stops the sniper service
	StopSniper() error

	// GetSniperStatus returns the current status of the sniper service
	GetSniperStatus() (string, error)

	// SetupAutoSnipe configures the sniper to automatically snipe new listings
	SetupAutoSnipe(enabled bool, config *port.SniperConfig) error
}

// sniperUseCase implements the SniperUseCase interface
type sniperUseCase struct {
	sniperService port.SniperService
	newCoinUC     NewCoinUseCase
	logger        zerolog.Logger

	// Auto-snipe configuration
	autoSnipeEnabled bool
	autoSnipeConfig  *port.SniperConfig
	autoSnipeMutex   sync.RWMutex
}

// NewSniperUseCase creates a new sniper use case
func NewSniperUseCase(
	sniperService port.SniperService,
	newCoinUC NewCoinUseCase,
	logger *zerolog.Logger,
) SniperUseCase {
	uc := &sniperUseCase{
		sniperService:    sniperService,
		newCoinUC:        newCoinUC,
		logger:           logger.With().Str("component", "sniper_usecase").Logger(),
		autoSnipeEnabled: false,
	}

	// Start the sniper service
	if err := uc.sniperService.Start(); err != nil {
		uc.logger.Error().Err(err).Msg("Failed to start sniper service")
	}

	// Subscribe to new coin events for auto-sniping
	if err := uc.setupNewCoinEventListener(); err != nil {
		uc.logger.Error().Err(err).Msg("Failed to set up new coin event listener")
	}

	return uc
}

// ExecuteSnipe executes a high-speed buy on a newly listed token
func (uc *sniperUseCase) ExecuteSnipe(ctx context.Context, symbol string) (*model.Order, error) {
	if uc.sniperService == nil {
		return nil, ErrSniperNotInitialized
	}

	uc.logger.Info().Str("symbol", symbol).Msg("Executing snipe")

	// Execute snipe
	order, err := uc.sniperService.ExecuteSnipe(ctx, symbol)
	if err != nil {
		uc.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to execute snipe")
		return nil, fmt.Errorf("failed to execute snipe: %w", err)
	}

	uc.logger.Info().
		Str("symbol", symbol).
		Str("orderID", order.OrderID).
		Float64("quantity", order.Quantity).
		Float64("price", order.Price).
		Msg("Snipe executed successfully")

	return order, nil
}

// ExecuteSnipeWithConfig executes a high-speed buy with custom configuration
func (uc *sniperUseCase) ExecuteSnipeWithConfig(ctx context.Context, symbol string, config *port.SniperConfig) (*model.Order, error) {
	if uc.sniperService == nil {
		return nil, ErrSniperNotInitialized
	}

	if config == nil {
		return nil, ErrInvalidSniperConfig
	}

	uc.logger.Info().
		Str("symbol", symbol).
		Float64("maxBuyAmount", config.MaxBuyAmount).
		Float64("maxPrice", config.MaxPricePerToken).
		Msg("Executing snipe with custom config")

	// Execute snipe with custom config
	order, err := uc.sniperService.ExecuteSnipeWithConfig(ctx, symbol, config)
	if err != nil {
		uc.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to execute snipe with custom config")
		return nil, fmt.Errorf("failed to execute snipe with custom config: %w", err)
	}

	uc.logger.Info().
		Str("symbol", symbol).
		Str("orderID", order.OrderID).
		Float64("quantity", order.Quantity).
		Float64("price", order.Price).
		Msg("Snipe with custom config executed successfully")

	return order, nil
}

// GetSniperConfig returns the current sniper configuration
func (uc *sniperUseCase) GetSniperConfig() (*port.SniperConfig, error) {
	if uc.sniperService == nil {
		return nil, ErrSniperNotInitialized
	}

	return uc.sniperService.GetConfig(), nil
}

// UpdateSniperConfig updates the sniper configuration
func (uc *sniperUseCase) UpdateSniperConfig(config *port.SniperConfig) error {
	if uc.sniperService == nil {
		return ErrSniperNotInitialized
	}

	if config == nil {
		return ErrInvalidSniperConfig
	}

	return uc.sniperService.UpdateConfig(config)
}

// StartSniper starts the sniper service
func (uc *sniperUseCase) StartSniper() error {
	if uc.sniperService == nil {
		return ErrSniperNotInitialized
	}

	return uc.sniperService.Start()
}

// StopSniper stops the sniper service
func (uc *sniperUseCase) StopSniper() error {
	if uc.sniperService == nil {
		return ErrSniperNotInitialized
	}

	return uc.sniperService.Stop()
}

// GetSniperStatus returns the current status of the sniper service
func (uc *sniperUseCase) GetSniperStatus() (string, error) {
	if uc.sniperService == nil {
		return "", ErrSniperNotInitialized
	}

	return uc.sniperService.GetStatus(), nil
}

// SetupAutoSnipe configures the sniper to automatically snipe new listings
func (uc *sniperUseCase) SetupAutoSnipe(enabled bool, config *port.SniperConfig) error {
	uc.autoSnipeMutex.Lock()
	defer uc.autoSnipeMutex.Unlock()

	uc.autoSnipeEnabled = enabled

	if enabled && config != nil {
		uc.autoSnipeConfig = config
	} else if enabled {
		// Use default config if none provided
		defaultConfig, err := uc.GetSniperConfig()
		if err != nil {
			return err
		}
		uc.autoSnipeConfig = defaultConfig
	}

	uc.logger.Info().
		Bool("enabled", enabled).
		Msg("Auto-snipe configuration updated")

	return nil
}

// setupNewCoinEventListener subscribes to new coin events for auto-sniping
func (uc *sniperUseCase) setupNewCoinEventListener() error {
	// Subscribe to new coin events
	err := uc.newCoinUC.SubscribeToEvents(func(event *model.CoinEvent) {
		// Check if auto-snipe is enabled
		uc.autoSnipeMutex.RLock()
		enabled := uc.autoSnipeEnabled
		config := uc.autoSnipeConfig
		uc.autoSnipeMutex.RUnlock()

		if !enabled || config == nil {
			return
		}

		// Only process status changes to trading
		if event.EventType == "status_change" && event.NewStatus == model.CoinStatusTrading {
			// Execute snipe in a new goroutine
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				uc.logger.Info().
					Str("coinID", event.CoinID).
					Str("eventType", event.EventType).
					Msg("Auto-sniping new listing")

				// Get coin details
				coin, err := uc.newCoinUC.GetCoinDetails(event.CoinID)
				if err != nil {
					uc.logger.Error().
						Err(err).
						Str("coinID", event.CoinID).
						Msg("Failed to get coin details for auto-snipe")
					return
				}

				// Execute snipe
				_, err = uc.sniperService.ExecuteSnipeWithConfig(ctx, coin.Symbol, config)
				if err != nil {
					uc.logger.Error().
						Err(err).
						Str("symbol", coin.Symbol).
						Msg("Failed to auto-snipe new listing")
					return
				}

				uc.logger.Info().
					Str("symbol", coin.Symbol).
					Msg("Auto-snipe executed successfully")
			}()
		}
	})

	return err
}

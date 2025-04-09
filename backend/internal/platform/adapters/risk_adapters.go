// Package adapters provides adapter implementations for various interfaces
package adapters

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/interfaces"
	"go-crypto-bot-clean/backend/internal/domain/repository"
	"go-crypto-bot-clean/backend/internal/domain/risk/controls"
	"go-crypto-bot-clean/backend/internal/platform/mexc/rest"
	"go.uber.org/zap"
)

// MexcPriceServiceAdapter adapts the MEXC client to the PriceService interface
type MexcPriceServiceAdapter struct {
	client *rest.Client
}

// NewMexcPriceServiceAdapter creates a new MexcPriceServiceAdapter
func NewMexcPriceServiceAdapter(client *rest.Client) *MexcPriceServiceAdapter {
	return &MexcPriceServiceAdapter{
		client: client,
	}
}

// GetPrice implements the PriceService interface
func (a *MexcPriceServiceAdapter) GetPrice(ctx context.Context, symbol string) (float64, error) {
	ticker, err := a.client.GetTicker(ctx, symbol)
	if err != nil {
		return 0, err
	}
	return ticker.Price, nil
}

// MexcAccountServiceAdapter adapts the MEXC client to the AccountService interface
type MexcAccountServiceAdapter struct {
	client *rest.Client
}

// NewMexcAccountServiceAdapter creates a new MexcAccountServiceAdapter
func NewMexcAccountServiceAdapter(client *rest.Client) *MexcAccountServiceAdapter {
	return &MexcAccountServiceAdapter{
		client: client,
	}
}

// GetBalance implements the AccountService interface
func (a *MexcAccountServiceAdapter) GetBalance(ctx context.Context) (float64, error) {
	wallet, err := a.client.GetWallet(ctx)
	if err != nil {
		return 0, err
	}

	// Get USDT balance
	var balance float64
	if usdtBalance, ok := wallet.Balances["USDT"]; ok {
		balance = usdtBalance.Free
	}

	return balance, nil
}

// BoughtCoinPositionAdapter adapts the BoughtCoinRepository to the PositionRepository interface
type BoughtCoinPositionAdapter struct {
	repo repository.BoughtCoinRepository
}

// NewBoughtCoinPositionAdapter creates a new BoughtCoinPositionAdapter
func NewBoughtCoinPositionAdapter(repo repository.BoughtCoinRepository) *BoughtCoinPositionAdapter {
	return &BoughtCoinPositionAdapter{
		repo: repo,
	}
}

// GetOpenPositions implements the PositionRepository interface
func (a *BoughtCoinPositionAdapter) GetOpenPositions(ctx context.Context) ([]controls.Position, error) {
	coins, err := a.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	positions := make([]controls.Position, 0, len(coins))
	for _, coin := range coins {
		positions = append(positions, controls.Position{
			Symbol:     coin.Symbol,
			Quantity:   coin.Quantity,
			EntryPrice: coin.BuyPrice,
		})
	}

	return positions, nil
}

// BoughtCoinTradeAdapter adapts the BoughtCoinRepository to the TradeRepository interface
type BoughtCoinTradeAdapter struct {
	repo repository.BoughtCoinRepository
}

// NewBoughtCoinTradeAdapter creates a new BoughtCoinTradeAdapter
func NewBoughtCoinTradeAdapter(repo repository.BoughtCoinRepository) *BoughtCoinTradeAdapter {
	return &BoughtCoinTradeAdapter{
		repo: repo,
	}
}

// GetTradesByDateRange implements the TradeRepository interface
func (a *BoughtCoinTradeAdapter) GetTradesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]controls.Trade, error) {
	// For now, we'll just return an empty slice
	// In a real implementation, we would query the trade history
	return []controls.Trade{}, nil
}

// ZapLoggerAdapter adapts the zap.Logger to the Logger interface
type ZapLoggerAdapter struct {
	logger *zap.Logger
}

// NewZapLoggerAdapter creates a new ZapLoggerAdapter
func NewZapLoggerAdapter(logger *zap.Logger) *ZapLoggerAdapter {
	return &ZapLoggerAdapter{
		logger: logger,
	}
}

// Info implements the Logger interface
func (a *ZapLoggerAdapter) Info(msg string, keysAndValues ...interface{}) {
	fields := convertToZapFields(keysAndValues)
	a.logger.Info(msg, fields...)
}

// Warn implements the Logger interface
func (a *ZapLoggerAdapter) Warn(msg string, keysAndValues ...interface{}) {
	fields := convertToZapFields(keysAndValues)
	a.logger.Warn(msg, fields...)
}

// Error implements the Logger interface
func (a *ZapLoggerAdapter) Error(msg string, keysAndValues ...interface{}) {
	fields := convertToZapFields(keysAndValues)
	a.logger.Error(msg, fields...)
}

// convertToZapFields converts a slice of interface{} to a slice of zap.Field
func convertToZapFields(keysAndValues []interface{}) []zap.Field {
	if len(keysAndValues) == 0 {
		return []zap.Field{}
	}

	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			continue
		}

		if i+1 < len(keysAndValues) {
			value := keysAndValues[i+1]
			switch v := value.(type) {
			case string:
				fields = append(fields, zap.String(key, v))
			case int:
				fields = append(fields, zap.Int(key, v))
			case float64:
				fields = append(fields, zap.Float64(key, v))
			case bool:
				fields = append(fields, zap.Bool(key, v))
			case time.Time:
				fields = append(fields, zap.Time(key, v))
			case error:
				fields = append(fields, zap.Error(v))
			default:
				fields = append(fields, zap.Any(key, v))
			}
		}
	}

	return fields
}

// PriceServiceAdapter adapts the PriceService interface to the controls.PriceService interface
type PriceServiceAdapter struct {
	priceService interfaces.PriceService
}

// NewPriceServiceAdapter creates a new PriceServiceAdapter
func NewPriceServiceAdapter(priceService interfaces.PriceService) *PriceServiceAdapter {
	return &PriceServiceAdapter{
		priceService: priceService,
	}
}

// GetPrice implements the controls.PriceService interface
func (a *PriceServiceAdapter) GetPrice(ctx context.Context, symbol string) (float64, error) {
	return a.priceService.GetPrice(ctx, symbol)
}

// PositionRepositoryAdapter adapts the PositionRepository interface to the controls.PositionRepository interface
type PositionRepositoryAdapter struct {
	positionRepo interfaces.PositionRepository
}

// NewPositionRepositoryAdapter creates a new PositionRepositoryAdapter
func NewPositionRepositoryAdapter(positionRepo interfaces.PositionRepository) *PositionRepositoryAdapter {
	return &PositionRepositoryAdapter{
		positionRepo: positionRepo,
	}
}

// GetOpenPositions implements the controls.PositionRepository interface
func (a *PositionRepositoryAdapter) GetOpenPositions(ctx context.Context) ([]controls.Position, error) {
	positions, err := a.positionRepo.FindAll(ctx, interfaces.PositionFilter{})
	if err != nil {
		return nil, err
	}

	result := make([]controls.Position, 0, len(positions))
	for _, pos := range positions {
		result = append(result, controls.Position{
			Symbol:     pos.Symbol,
			Quantity:   pos.Quantity,
			EntryPrice: pos.EntryPrice,
		})
	}

	return result, nil
}

// AccountServiceAdapter adapts the ExchangeService interface to the controls.AccountService interface
type AccountServiceAdapter struct {
	exchangeSvc interfaces.ExchangeService
}

// NewAccountServiceAdapter creates a new AccountServiceAdapter
func NewAccountServiceAdapter(exchangeSvc interfaces.ExchangeService) *AccountServiceAdapter {
	return &AccountServiceAdapter{
		exchangeSvc: exchangeSvc,
	}
}

// GetBalance implements the controls.AccountService interface
func (a *AccountServiceAdapter) GetBalance(ctx context.Context) (float64, error) {
	wallet, err := a.exchangeSvc.GetWallet(ctx)
	if err != nil {
		return 0, err
	}

	// Get USDT balance
	var balance float64
	if usdtBalance, ok := wallet.Balances["USDT"]; ok {
		balance = usdtBalance.Free
	}

	return balance, nil
}

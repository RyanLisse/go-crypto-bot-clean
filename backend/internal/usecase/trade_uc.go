package usecase

import (
	"context"
	"errors"

	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/rs/zerolog"
)

// Common errors
var (
	ErrInvalidOrderData = errors.New("invalid order data")
	ErrOrderNotFound    = errors.New("order not found")
)

// TradeUseCase defines methods for trade operations
type TradeUseCase interface {
	// Place a new order
	PlaceOrder(ctx context.Context, req model.OrderRequest) (*model.Order, error)
}

// tradeUseCase implements the TradeUseCase interface
type tradeUseCase struct {
	mexcAPI    port.MexcAPI
	orderRepo  port.OrderRepository
	symbolRepo port.SymbolRepository
	logger     zerolog.Logger
}

// NewTradeUseCase creates a new TradeUseCase
func NewTradeUseCase(
	mexcAPI port.MexcAPI,
	orderRepo port.OrderRepository,
	symbolRepo port.SymbolRepository,
	logger zerolog.Logger,
) TradeUseCase {
	return &tradeUseCase{
		mexcAPI:    mexcAPI,
		orderRepo:  orderRepo,
		symbolRepo: symbolRepo,
		logger:     logger.With().Str("component", "trade_usecase").Logger(),
	}
}

// PlaceOrder places a new order
func (uc *tradeUseCase) PlaceOrder(ctx context.Context, req model.OrderRequest) (*model.Order, error) {
	// Validate symbol exists
	symbol, err := uc.symbolRepo.GetBySymbol(ctx, req.Symbol)
	if err != nil {
		uc.logger.Error().Err(err).Str("symbol", req.Symbol).Msg("Failed to validate symbol")
		return nil, err
	}
	if symbol == nil {
		uc.logger.Warn().Str("symbol", req.Symbol).Msg("Symbol not found")
		return nil, ErrSymbolNotFound
	}

	// Place order on exchange
	timeInForce := model.TimeInForceGTC // Default time in force
	order, err := uc.mexcAPI.PlaceOrder(
		ctx,
		req.Symbol,
		req.Side,
		req.Type,
		req.Quantity,
		req.Price,
		timeInForce,
	)
	if err != nil {
		uc.logger.Error().Err(err).
			Str("symbol", req.Symbol).
			Str("side", string(req.Side)).
			Str("type", string(req.Type)).
			Float64("quantity", req.Quantity).
			Float64("price", req.Price).
			Msg("Failed to place order on exchange")
		return nil, err
	}

	// Save order to repository
	err = uc.orderRepo.Create(ctx, order)
	if err != nil {
		uc.logger.Error().Err(err).
			Str("orderId", order.ID).
			Str("symbol", order.Symbol).
			Msg("Failed to save order to repository")
		// Note: We still return the order even if saving to repository fails
		// because the order was successfully placed on the exchange
	}

	uc.logger.Info().
		Str("orderId", order.ID).
		Str("symbol", order.Symbol).
		Str("side", string(order.Side)).
		Str("type", string(order.Type)).
		Float64("quantity", order.Quantity).
		Float64("price", order.Price).
		Msg("Order placed successfully")

	return order, nil
}

package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// Common errors
var (
	ErrInvalidOrderData    = errors.New("invalid order data")
	ErrOrderNotFound       = errors.New("order not found")
	ErrInsufficientBalance = errors.New("insufficient balance for order")
	ErrSymbolNotFound      = errors.New("symbol not found")
)

// TradeUseCase defines methods for trade operations
type TradeUseCase interface {
	// Place a new order
	PlaceOrder(ctx context.Context, req model.OrderRequest) (*model.Order, error)
	// Cancel an existing order
	CancelOrder(ctx context.Context, symbol, orderID string) error
	// Get the current status of an order
	GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error)
	// Get all open orders for a symbol
	GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error)
	// Get order history for a symbol with pagination
	GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error)
	// Calculate the required quantity for an order based on amount in quote currency
	CalculateRequiredQuantity(ctx context.Context, symbol string, side model.OrderSide, amount float64) (float64, error)
}

// tradeUseCase implements the TradeUseCase interface
type tradeUseCase struct {
	mexcClient   port.MEXCClient
	orderRepo    port.OrderRepository
	symbolRepo   port.SymbolRepository
	tradeService port.TradeService
	riskUC       RiskUseCase
	txManager    port.TransactionManager
	logger       zerolog.Logger
}

// NewTradeUseCase creates a new TradeUseCase
func NewTradeUseCase(
	mexcClient port.MEXCClient,
	orderRepo port.OrderRepository,
	symbolRepo port.SymbolRepository,
	tradeService port.TradeService,
	riskUC RiskUseCase,
	txManager port.TransactionManager,
	logger zerolog.Logger,
) TradeUseCase {
	return &tradeUseCase{
		mexcClient:   mexcClient,
		orderRepo:    orderRepo,
		symbolRepo:   symbolRepo,
		tradeService: tradeService,
		riskUC:       riskUC,
		txManager:    txManager,
		logger:       logger.With().Str("component", "trade_usecase").Logger(),
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

	// Perform risk assessment before placing the order
	if uc.riskUC != nil {
		allowed, assessments, err := uc.riskUC.EvaluateOrderRisk(ctx, req.UserID, req)
		if err != nil {
			uc.logger.Error().Err(err).
				Str("symbol", req.Symbol).
				Str("side", string(req.Side)).
				Msg("Failed to evaluate order risk")
			return nil, fmt.Errorf("failed to evaluate risk: %w", err)
		}

		if !allowed {
			// Log risk assessments
			for _, assessment := range assessments {
				if assessment.Level == model.RiskLevelHigh || assessment.Level == model.RiskLevelCritical {
					uc.logger.Warn().
						Str("riskType", string(assessment.Type)).
						Str("riskLevel", string(assessment.Level)).
						Str("message", assessment.Message).
						Str("recommendation", assessment.Recommendation).
						Msg("High risk detected")
				}
			}
			return nil, errors.New("order rejected due to risk assessment: " + getHighestRiskMessage(assessments))
		}
	}

	// Use transaction manager to ensure atomicity of order placement
	var response *model.OrderResponse
	err = uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Delegate to the trade service within transaction
		resp, err := uc.tradeService.PlaceOrder(txCtx, &req)
		if err != nil {
			uc.logger.Error().Err(err).
				Str("symbol", req.Symbol).
				Str("side", string(req.Side)).
				Str("type", string(req.Type)).
				Float64("quantity", req.Quantity).
				Float64("price", req.Price).
				Msg("Failed to place order")
			return err
		}

		// Save order response to use outside transaction
		response = resp
		return nil
	})

	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("order response is nil after transaction")
	}

	uc.logger.Info().
		Str("orderId", response.OrderID).
		Str("symbol", response.Symbol).
		Str("side", string(response.Side)).
		Str("type", string(response.Type)).
		Float64("quantity", response.Quantity).
		Float64("price", response.Price).
		Msg("Order placed successfully")

	return &response.Order, nil
}

// getHighestRiskMessage returns the message from the highest risk assessment
func getHighestRiskMessage(assessments []*model.RiskAssessment) string {
	if len(assessments) == 0 {
		return "unknown risk"
	}

	// Priority: CRITICAL > HIGH > MEDIUM > LOW
	var criticalRisk, highRisk, mediumRisk, lowRisk *model.RiskAssessment

	for _, assessment := range assessments {
		switch assessment.Level {
		case model.RiskLevelCritical:
			criticalRisk = assessment
		case model.RiskLevelHigh:
			if highRisk == nil {
				highRisk = assessment
			}
		case model.RiskLevelMedium:
			if mediumRisk == nil {
				mediumRisk = assessment
			}
		case model.RiskLevelLow:
			if lowRisk == nil {
				lowRisk = assessment
			}
		}
	}

	if criticalRisk != nil {
		return criticalRisk.Message
	}
	if highRisk != nil {
		return highRisk.Message
	}
	if mediumRisk != nil {
		return mediumRisk.Message
	}
	if lowRisk != nil {
		return lowRisk.Message
	}

	return "unknown risk"
}

// CancelOrder cancels an existing order
func (uc *tradeUseCase) CancelOrder(ctx context.Context, symbol, orderID string) error {
	// Delegate to the trade service
	err := uc.tradeService.CancelOrder(ctx, symbol, orderID)
	if err != nil {
		uc.logger.Error().Err(err).
			Str("symbol", symbol).
			Str("orderId", orderID).
			Msg("Failed to cancel order")
		return err
	}

	uc.logger.Info().
		Str("symbol", symbol).
		Str("orderId", orderID).
		Msg("Order canceled successfully")

	return nil
}

// GetOrderStatus retrieves the current status of an order
func (uc *tradeUseCase) GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	// Delegate to the trade service
	order, err := uc.tradeService.GetOrderStatus(ctx, symbol, orderID)
	if err != nil {
		uc.logger.Error().Err(err).
			Str("symbol", symbol).
			Str("orderId", orderID).
			Msg("Failed to get order status")
		return nil, err
	}

	return order, nil
}

// GetOpenOrders retrieves all open orders for a symbol
func (uc *tradeUseCase) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	// Delegate to the trade service
	orders, err := uc.tradeService.GetOpenOrders(ctx, symbol)
	if err != nil {
		uc.logger.Error().Err(err).
			Str("symbol", symbol).
			Msg("Failed to get open orders")
		return nil, err
	}

	return orders, nil
}

// GetOrderHistory retrieves order history for a symbol with pagination
func (uc *tradeUseCase) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	// Delegate to the trade service
	orders, err := uc.tradeService.GetOrderHistory(ctx, symbol, limit, offset)
	if err != nil {
		uc.logger.Error().Err(err).
			Str("symbol", symbol).
			Int("limit", limit).
			Int("offset", offset).
			Msg("Failed to get order history")
		return nil, err
	}

	return orders, nil
}

// CalculateRequiredQuantity calculates the required quantity for an order based on amount
func (uc *tradeUseCase) CalculateRequiredQuantity(ctx context.Context, symbol string, side model.OrderSide, amount float64) (float64, error) {
	// Delegate to the trade service
	quantity, err := uc.tradeService.CalculateRequiredQuantity(ctx, symbol, side, amount)
	if err != nil {
		uc.logger.Error().Err(err).
			Str("symbol", symbol).
			Str("side", string(side)).
			Float64("amount", amount).
			Msg("Failed to calculate required quantity")
		return 0, err
	}

	return quantity, nil
}

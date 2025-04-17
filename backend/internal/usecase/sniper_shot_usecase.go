package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/service"
	"github.com/rs/zerolog"
)

// Common error definitions
var (
	ErrInvalidRequest    = errors.New("invalid sniper shot request")
	ErrInsufficientFunds = errors.New("insufficient funds for sniper shot")
	ErrTradeExecution    = errors.New("error executing sniper shot trade")
	ErrOrderNotFound     = errors.New("sniper shot order not found")
)

// SniperShotServicer defines the interface for sniper shot operations
type SniperShotServicer interface {
	ExecuteSniper(ctx context.Context, req *service.SniperShotRequest) (*service.SniperShotResult, error)
	CancelSniper(ctx context.Context, symbol, orderID string) error
	GetSniperOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error)
}

// SniperShotUseCase orchestrates the sniper shot trading functionality
type SniperShotUseCase struct {
	logger        *zerolog.Logger
	sniperService SniperShotServicer
	walletRepo    port.WalletRepository
	orderRepo     port.OrderRepository
}

// NewSniperShotUseCase creates a new instance of the SniperShot use case
func NewSniperShotUseCase(
	logger *zerolog.Logger,
	sniperService SniperShotServicer,
	walletRepo port.WalletRepository,
	orderRepo port.OrderRepository,
) *SniperShotUseCase {
	return &SniperShotUseCase{
		logger:        logger,
		sniperService: sniperService,
		walletRepo:    walletRepo,
		orderRepo:     orderRepo,
	}
}

// SniperShotParams represents the parameters for executing a sniper shot
type SniperShotParams struct {
	UserID         string                 // ID of the user making the request
	Symbol         string                 // Symbol to trade (e.g., "BTCUSDT")
	Side           model.OrderSide        // BUY or SELL
	Quantity       float64                // Amount to buy or sell
	Price          float64                // Price limit (0 for market orders)
	Type           model.OrderType        // LIMIT or MARKET
	TimeLimit      time.Duration          // Maximum time to try executing the order
	PriceThreshold float64                // Optional price threshold to trigger the trade
	ComparisonType service.ComparisonType // Optional comparison type (Above, Below)
	MaxSlippage    float64                // Maximum allowed slippage percentage
}

// ValidateParams validates the parameters for a sniper shot
func (uc *SniperShotUseCase) ValidateParams(params *SniperShotParams) error {
	// Check required fields
	if params.UserID == "" {
		return fmt.Errorf("%w: missing user ID", ErrInvalidRequest)
	}
	if params.Symbol == "" {
		return fmt.Errorf("%w: missing trading symbol", ErrInvalidRequest)
	}
	if params.Quantity <= 0 {
		return fmt.Errorf("%w: quantity must be greater than zero", ErrInvalidRequest)
	}

	// Validate order type and required fields
	if params.Type == model.OrderTypeLimit && params.Price <= 0 {
		return fmt.Errorf("%w: price must be specified for limit orders", ErrInvalidRequest)
	}

	// Validate side
	if params.Side != model.OrderSideBuy && params.Side != model.OrderSideSell {
		return fmt.Errorf("%w: invalid order side", ErrInvalidRequest)
	}

	// If price threshold is set, validate comparison type
	if params.PriceThreshold > 0 {
		if params.ComparisonType != service.Above && params.ComparisonType != service.Below {
			return fmt.Errorf("%w: invalid comparison type for price threshold", ErrInvalidRequest)
		}
	}

	// Validate time limit
	if params.TimeLimit <= 0 {
		// Set a default time limit if not specified
		params.TimeLimit = 30 * time.Second
	}

	return nil
}

// ExecuteSniper validates parameters, checks balance, and executes a sniper shot trade
func (uc *SniperShotUseCase) ExecuteSniper(ctx context.Context, params *SniperShotParams) (*service.SniperShotResult, error) {
	// Validate parameters
	if err := uc.ValidateParams(params); err != nil {
		return nil, err
	}

	// Check user balance
	wallet, err := uc.walletRepo.GetByUserID(ctx, params.UserID)
	if err != nil {
		uc.logger.Error().Err(err).Str("userID", params.UserID).Msg("Failed to retrieve user wallet")
		return nil, fmt.Errorf("failed to check user balance: %w", err)
	}

	// Extract base and quote assets from the symbol (e.g., BTCUSDT -> BTC, USDT)
	baseAsset, quoteAsset := extractAssetsFromSymbol(params.Symbol)

	// Check if user has sufficient balance
	if !uc.hasEnoughBalance(wallet, params, baseAsset, quoteAsset) {
		return nil, ErrInsufficientFunds
	}

	// Prepare sniper shot request
	var condition *service.TriggerCondition
	if params.PriceThreshold > 0 {
		condition = &service.TriggerCondition{
			PriceThreshold: params.PriceThreshold,
			Comparison:     params.ComparisonType,
			MaxSlippage:    params.MaxSlippage,
		}
	}

	sniperReq := &service.SniperShotRequest{
		UserID:    params.UserID,
		Symbol:    params.Symbol,
		Side:      params.Side,
		Quantity:  params.Quantity,
		Price:     params.Price,
		Type:      params.Type,
		TimeLimit: params.TimeLimit,
		Condition: condition,
	}

	// Execute the sniper shot
	result, err := uc.sniperService.ExecuteSniper(ctx, sniperReq)
	if err != nil {
		uc.logger.Error().Err(err).
			Str("userID", params.UserID).
			Str("symbol", params.Symbol).
			Msg("Failed to execute sniper shot")
		return nil, fmt.Errorf("%w: %v", ErrTradeExecution, err)
	}

	return result, nil
}

// CancelSniper cancels a sniper shot order
func (uc *SniperShotUseCase) CancelSniper(ctx context.Context, userID, orderID string) error {
	// First, get the order details to verify it exists and belongs to the user
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrOrderNotFound, err)
	}

	// Verify the order belongs to the user
	if order.UserID != userID {
		return fmt.Errorf("%w: order does not belong to the user", ErrInvalidRequest)
	}

	// Cancel the order
	if err := uc.sniperService.CancelSniper(ctx, order.Symbol, order.OrderID); err != nil {
		return fmt.Errorf("failed to cancel sniper shot: %w", err)
	}

	return nil
}

// GetOrderStatus retrieves the status of a sniper shot order
func (uc *SniperShotUseCase) GetOrderStatus(ctx context.Context, userID, orderID string) (*model.Order, error) {
	// First, get the order details from the repository to verify ownership
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOrderNotFound, err)
	}

	// Verify the order belongs to the user
	if order.UserID != userID {
		return nil, fmt.Errorf("%w: order does not belong to the user", ErrInvalidRequest)
	}

	// Get the current status from the exchange
	currentOrder, err := uc.sniperService.GetSniperOrderStatus(ctx, order.Symbol, order.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status: %w", err)
	}

	// If the order status has changed, update it in the repository
	if currentOrder.Status != order.Status {
		order.Status = currentOrder.Status
		order.ExecutedQty = currentOrder.ExecutedQty
		order.AvgFillPrice = currentOrder.AvgFillPrice
		order.UpdatedAt = time.Now()

		if err := uc.orderRepo.Update(ctx, order); err != nil {
			uc.logger.Warn().
				Err(err).
				Str("orderID", orderID).
				Msg("Failed to update order status in repository")
		}
	}

	return currentOrder, nil
}

// hasEnoughBalance checks if the user has enough balance for the sniper shot
func (uc *SniperShotUseCase) hasEnoughBalance(wallet *model.Wallet, params *SniperShotParams, baseAsset, quoteAsset model.Asset) bool {
	// For BUY orders, check if user has enough quote asset (e.g., USDT)
	if params.Side == model.OrderSideBuy {
		estimatedCost := params.Quantity * params.Price
		if params.Type == model.OrderTypeMarket {
			// Add a buffer for market orders to account for slippage
			estimatedCost *= 1.05
		}

		// Check if the wallet has the quote asset with sufficient balance
		for _, balance := range wallet.Balances {
			if balance.Asset == quoteAsset && balance.Free >= estimatedCost {
				return true
			}
		}
		return false
	}

	// For SELL orders, check if user has enough base asset (e.g., BTC)
	if params.Side == model.OrderSideSell {
		for _, balance := range wallet.Balances {
			if balance.Asset == baseAsset && balance.Free >= params.Quantity {
				return true
			}
		}
		return false
	}

	return false
}

// extractAssetsFromSymbol extracts base and quote assets from a symbol (e.g., BTCUSDT -> BTC, USDT)
func extractAssetsFromSymbol(symbol string) (model.Asset, model.Asset) {
	// Common quote assets to look for
	quoteAssets := []string{"USDT", "BTC", "ETH", "BNB", "BUSD", "USDC"}

	for _, quote := range quoteAssets {
		if len(symbol) > len(quote) && symbol[len(symbol)-len(quote):] == quote {
			base := symbol[:len(symbol)-len(quote)]
			return model.Asset(base), model.Asset(quote)
		}
	}

	// Default fallback - best guess dividing the symbol
	if len(symbol) >= 6 {
		base := symbol[:len(symbol)-4]
		quote := symbol[len(symbol)-4:]
		return model.Asset(base), model.Asset(quote)
	} else if len(symbol) >= 3 {
		base := symbol[:3]
		quote := symbol[3:]
		return model.Asset(base), model.Asset(quote)
	}

	// If we cannot determine, return empty
	return "", ""
}

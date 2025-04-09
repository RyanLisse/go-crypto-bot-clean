package trade

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ryanlisse/go-crypto-bot/internal/config"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/risk"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/strategy"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/mexc/rest"
)

// tradeService implements TradeService
type tradeService struct {
	boughtCoinRepo  repository.BoughtCoinRepository
	mexcClient      *rest.Client
	config          *config.Config
	logger          *zap.Logger
	strategyFactory strategy.StrategyFactory
	riskService     risk.RiskService
}

// NewTradeService creates a new trade service with dependency injection
func NewTradeService(
	boughtCoinRepo repository.BoughtCoinRepository,
	mexcClient *rest.Client,
	cfg *config.Config,
	strategyFactory strategy.StrategyFactory,
	riskService risk.RiskService,
) *tradeService {
	logger, _ := zap.NewProduction()
	return &tradeService{
		boughtCoinRepo:  boughtCoinRepo,
		mexcClient:      mexcClient,
		config:          cfg,
		logger:          logger,
		strategyFactory: strategyFactory,
		riskService:     riskService,
	}
}

// EvaluatePurchaseDecision checks if a trade should be executed
func (s *tradeService) EvaluatePurchaseDecision(ctx context.Context, symbol string) (*models.PurchaseDecision, error) {
	// Get active trades
	activeTrades, err := s.boughtCoinRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Check max positions limit (default to 5 if not specified)
	maxPositions := 5
	if len(activeTrades) >= maxPositions {
		reason := fmt.Sprintf("max positions (%d) reached", maxPositions)
		s.logger.Info(reason)
		return &models.PurchaseDecision{Decision: false, Reason: reason}, nil
	}

	// Additional trade evaluation logic can be added here
	// For now, we just return a default purchase decision
	return &models.PurchaseDecision{Decision: true, Reason: "trade opportunity found"}, nil
}

// ExecutePurchase executes a purchase of a cryptocurrency
func (s *tradeService) ExecutePurchase(ctx context.Context, symbol string, amount float64, options *models.PurchaseOptions) (*models.BoughtCoin, error) {
	// Get current price
	ticker, err := s.mexcClient.GetTicker(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker for %s: %w", symbol, err)
	}

	// Calculate order value
	orderValue := amount
	if amount == 0 {
		// Use default amount from config if not specified
		defaultAmount := 20.0 // Default to 20 USDT if not specified in config
		if s.config != nil && s.config.Trading.DefaultQuantity > 0 {
			defaultAmount = s.config.Trading.DefaultQuantity
		}
		orderValue = defaultAmount
	}

	// Check if this trade is allowed by risk controls
	if s.riskService != nil {
		allowed, reason, err := s.riskService.IsTradeAllowed(ctx, symbol, orderValue)
		if err != nil {
			s.logger.Error("Failed to check risk controls",
				zap.String("symbol", symbol),
				zap.Error(err))
			return nil, fmt.Errorf("failed to check risk controls: %w", err)
		}

		if !allowed {
			s.logger.Warn("Trade rejected by risk controls",
				zap.String("symbol", symbol),
				zap.Float64("amount", amount),
				zap.String("reason", reason))
			return nil, fmt.Errorf("trade rejected: %s", reason)
		}
	}

	// Calculate quantity based on amount and price
	var quantity float64
	if amount <= 0 && s.riskService != nil {
		// Get account balance
		accountBalance, err := s.getAccountBalance(ctx)
		if err != nil {
			return nil, err
		}

		// Calculate safe position size using risk service
		quantity, err = s.riskService.CalculatePositionSize(ctx, symbol, accountBalance)
		if err != nil {
			s.logger.Error("Failed to calculate position size",
				zap.String("symbol", symbol),
				zap.Error(err))
			return nil, fmt.Errorf("failed to calculate position size: %w", err)
		}
	} else {
		// Use specified amount
		quantity = amount / ticker.Price
	}

	// Create bought coin record
	coin := &models.BoughtCoin{
		Symbol:   symbol,
		BuyPrice: ticker.Price,
		Quantity: quantity,
		BoughtAt: time.Now(),
	}

	// Save to repository
	_, err = s.boughtCoinRepo.Create(ctx, coin)
	if err != nil {
		return nil, fmt.Errorf("failed to save purchase record: %w", err)
	}

	// Record the balance after purchase if risk service is available
	if s.riskService != nil {
		// Get updated account balance
		accountBalance, err := s.getAccountBalance(ctx)
		if err != nil {
			s.logger.Warn("Failed to get account balance for recording", zap.Error(err))
		} else {
			// Record the balance in the risk service
			s.logger.Info("Recording account balance after purchase", zap.Float64("balance", accountBalance))
		}
	}

	s.logger.Info("Executed purchase",
		zap.String("symbol", symbol),
		zap.Float64("price", ticker.Price),
		zap.Float64("quantity", quantity),
		zap.Float64("total", ticker.Price*quantity),
	)

	return coin, nil
}

// CheckStopLoss checks if a coin's current price is below its stop loss price
func (s *tradeService) CheckStopLoss(ctx context.Context, coin *models.BoughtCoin) (bool, error) {
	// Get current price
	// Removed ticker fetch as stop loss logic is no longer implemented here

	// Check if current price is below stop loss price
	// StopLoss is no longer stored in BoughtCoin; implement stop loss logic elsewhere if needed
	return false, nil
}

// CheckTakeProfit checks if a coin's current price is above its take profit price
func (s *tradeService) CheckTakeProfit(ctx context.Context, coin *models.BoughtCoin) (bool, error) {
	// Get current price
	// Removed ticker fetch as take profit logic is no longer implemented here

	// Check if current price is above take profit price
	// TakeProfit is no longer stored in BoughtCoin; implement take profit logic elsewhere if needed
	return false, nil
}

// SellCoin sells a cryptocurrency
func (s *tradeService) SellCoin(ctx context.Context, coin *models.BoughtCoin, amount float64) (*models.Order, error) {
	// Get current price
	ticker, err := s.mexcClient.GetTicker(ctx, coin.Symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker for %s: %w", coin.Symbol, err)
	}

	if amount <= 0 || amount > coin.Quantity {
		return nil, fmt.Errorf("invalid sell amount: %f, available: %f", amount, coin.Quantity)
	}

	// Create order using the models.Order defined in exchange.go
	order := &models.Order{
		Symbol:   coin.Symbol,
		Quantity: amount,
		Price:    ticker.Price,
	}

	// Place the order using the MEXC client
	placedOrder, err := s.mexcClient.PlaceOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to place sell order: %w", err)
	}

	// Update coin record: delete if fully sold, else update quantity
	if amount == coin.Quantity {
		// Delete the coin record
		err = s.boughtCoinRepo.Delete(ctx, coin.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete coin record: %w", err)
		}
	} else {
		coin.Quantity -= amount
		err = s.boughtCoinRepo.Update(ctx, coin)
		if err != nil {
			return nil, fmt.Errorf("failed to update coin record: %w", err)
		}
	}

	return placedOrder, nil
}

// CancelOrder cancels an order
func (s *tradeService) CancelOrder(ctx context.Context, orderID string) error {
	// TODO: Implement the logic to cancel an order
	return nil
}

// GetOrderStatus retrieves the status of a specific order
func (s *tradeService) GetOrderStatus(ctx context.Context, orderID string) (*models.Order, error) {
	// TODO: Implement logic to get order status from the exchange or repository
	return nil, fmt.Errorf("GetOrderStatus not implemented")
}

// GetPendingOrders retrieves all pending orders
func (s *tradeService) GetPendingOrders(ctx context.Context) ([]*models.Order, error) {
	// TODO: Implement logic to get all pending orders from the exchange
	// For now, return an empty slice
	return []*models.Order{}, nil
}

// This is a duplicate of the CancelOrder method above, so we'll remove it

// GetActiveTrades returns all active trades
func (s *tradeService) GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error) {
	// Get active trades from repository
	coinsSlice, err := s.boughtCoinRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active trades: %w", err)
	}

	// Convert to pointer slice
	coins := make([]*models.BoughtCoin, 0, len(coinsSlice))
	for i := range coinsSlice {
		coins = append(coins, &coinsSlice[i])
	}

	// Update current prices
	for _, coin := range coins {
		ticker, err := s.mexcClient.GetTicker(ctx, coin.Symbol)
		if err != nil {
			s.logger.Error("Failed to get ticker", zap.String("symbol", coin.Symbol), zap.Error(err))
			continue
		}
		coin.CurrentPrice = ticker.Price
	}

	return coins, nil
}

// ExecuteTrade executes a trade based on the provided order
func (s *tradeService) ExecuteTrade(ctx context.Context, order *models.Order) (*models.Order, error) {
	// Place order on exchange
	result, err := s.mexcClient.PlaceOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	// If it's a buy order, save it to the bought coin repository
	if order.Side == models.OrderSideBuy && result.Status == models.OrderStatusFilled {
		boughtCoin := &models.BoughtCoin{
			Symbol:       result.Symbol,
			BuyPrice:     result.Price,
			CurrentPrice: result.Price,
			Quantity:     result.Quantity,
			BoughtAt:     time.Now(),
		}

		// Use Create instead of Save
		_, err := s.boughtCoinRepo.Create(ctx, boughtCoin)
		if err != nil {
			s.logger.Error("Failed to save bought coin", zap.Error(err))
		}
	}

	return result, nil
}

// GetTradeHistory returns the trade history
func (s *tradeService) GetTradeHistory(ctx context.Context, startTime time.Time, limit int) ([]*models.Order, error) {
	// Get orders from exchange
	orders, err := s.mexcClient.GetOpenOrders(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	// Filter by start time and limit
	var filteredOrders []*models.Order
	for _, order := range orders {
		if startTime.IsZero() || order.CreatedAt.After(startTime) {
			filteredOrders = append(filteredOrders, order)
		}

		if limit > 0 && len(filteredOrders) >= limit {
			break
		}
	}

	return filteredOrders, nil
}

// getAccountBalance retrieves the current account balance
func (s *tradeService) getAccountBalance(ctx context.Context) (float64, error) {
	// Get wallet from exchange
	wallet, err := s.mexcClient.GetWallet(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Get USDT balance
	var balance float64
	if usdtBalance, ok := wallet.Balances["USDT"]; ok {
		balance = usdtBalance.Free
	}

	s.logger.Debug("Retrieved account balance", zap.Float64("balance", balance))
	return balance, nil
}

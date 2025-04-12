package api

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repository"
	"go-crypto-bot-clean/backend/internal/platform/mexc/rest"

	"go.uber.org/zap"
)

// RealPortfolioService implements the PortfolioServiceInterface using real MEXC API data
type RealPortfolioService struct {
	mexcClient     *rest.Client
	boughtCoinRepo repository.BoughtCoinRepository
	accountService interface {
		GetWallet(ctx context.Context) (*models.Wallet, error)
	}
	logger *zap.Logger
}

// NewRealPortfolioService creates a new real portfolio service
func NewRealPortfolioService(
	mexcClient *rest.Client,
	boughtCoinRepo repository.BoughtCoinRepository,
	accountService interface {
		GetWallet(ctx context.Context) (*models.Wallet, error)
	},
	logger *zap.Logger,
) *RealPortfolioService {
	return &RealPortfolioService{
		mexcClient:     mexcClient,
		boughtCoinRepo: boughtCoinRepo,
		accountService: accountService,
		logger:         logger,
	}
}

// GetPortfolioValue returns the total portfolio value using real data from MEXC
func (s *RealPortfolioService) GetPortfolioValue(ctx context.Context) (float64, error) {
	s.logger.Debug("Getting portfolio value from real MEXC API")

	// Get wallet from account service
	wallet, err := s.accountService.GetWallet(ctx)
	if err != nil {
		s.logger.Error("Failed to get wallet", zap.Error(err))
		return 0, err
	}

	// Calculate total value
	var totalValue float64
	for _, balance := range wallet.Balances {
		// Get current price for the asset
		price := balance.Price
		if price == 0 {
			// If price is not set in the balance, try to get it from the ticker
			ticker, err := s.mexcClient.GetTicker(ctx, balance.Asset+"USDT")
			if err == nil && ticker != nil {
				price = ticker.Price
			}
		}

		// Add to total value
		totalValue += balance.Total * price
	}

	s.logger.Debug("Got portfolio value", zap.Float64("value", totalValue))
	return totalValue, nil
}

// GetActiveTrades returns all active trades
func (s *RealPortfolioService) GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error) {
	s.logger.Debug("Getting active trades from database")

	// Get all bought coins from the repository
	coins, err := s.boughtCoinRepo.FindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get active trades", zap.Error(err))
		return nil, err
	}

	// Update current prices for all coins
	for _, coin := range coins {
		// Get current price from ticker
		ticker, err := s.mexcClient.GetTicker(ctx, coin.Symbol)
		if err == nil && ticker != nil {
			coin.CurrentPrice = ticker.Price
			coin.UpdatedAt = time.Now()

			// Update the coin in the repository
			if err := s.boughtCoinRepo.Update(ctx, &coin); err != nil {
				s.logger.Warn("Failed to update coin price",
					zap.String("symbol", coin.Symbol),
					zap.Error(err))
			}
		}
	}

	s.logger.Debug("Got active trades", zap.Int("count", len(coins)))
	// Convert []models.BoughtCoin to []*models.BoughtCoin
	result := make([]*models.BoughtCoin, len(coins))
	for i := range coins {
		result[i] = &coins[i]
	}
	return result, nil
}

// GetTradePerformance returns performance metrics for trades
func (s *RealPortfolioService) GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error) {
	s.logger.Debug("Getting trade performance", zap.String("timeRange", timeRange))

	// Get all bought coins
	coins, err := s.boughtCoinRepo.FindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get bought coins", zap.Error(err))
		return nil, err
	}

	// Calculate performance metrics
	metrics := &models.PerformanceMetrics{
		TotalTrades: len(coins),
	}

	var totalProfit float64
	var winningTrades, losingTrades int
	var largestProfit, largestLoss float64

	for _, coin := range coins {
		profit := (coin.CurrentPrice - coin.PurchasePrice) * coin.Quantity

		if profit > 0 {
			winningTrades++
			if profit > largestProfit {
				largestProfit = profit
			}
		} else {
			losingTrades++
			if profit < largestLoss {
				largestLoss = profit
			}
		}

		totalProfit += profit
	}

	metrics.WinningTrades = winningTrades
	metrics.LosingTrades = losingTrades
	metrics.TotalProfitLoss = totalProfit

	if metrics.TotalTrades > 0 {
		metrics.WinRate = float64(winningTrades) / float64(metrics.TotalTrades) * 100
		metrics.AverageProfitPerTrade = totalProfit / float64(metrics.TotalTrades)
	}

	metrics.LargestProfit = largestProfit
	metrics.LargestLoss = largestLoss

	s.logger.Debug("Got trade performance",
		zap.Int("totalTrades", metrics.TotalTrades),
		zap.Float64("winRate", metrics.WinRate),
		zap.Float64("totalProfit", metrics.TotalProfitLoss))

	return metrics, nil
}

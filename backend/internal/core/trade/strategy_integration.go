package trade

import (
	"context"
	"fmt"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/strategy"
)

// EvaluateWithStrategy evaluates a trading decision using the strategy framework
func (s *tradeService) EvaluateWithStrategy(ctx context.Context, symbol string) (*models.PurchaseDecision, error) {
	// Get market data for the symbol
	marketData, err := s.prepareMarketData(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare market data: %w", err)
	}

	// Get the appropriate strategy for the current market regime
	strategy, err := s.strategyFactory.GetStrategyForMarketRegime(ctx, marketData)
	if err != nil {
		return nil, fmt.Errorf("failed to get strategy for market regime: %w", err)
	}

	// Get the latest candle
	candles, err := s.getCandles(ctx, symbol, "1h", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get candles: %w", err)
	}

	if len(candles) == 0 {
		return nil, fmt.Errorf("no candles available for %s", symbol)
	}

	// Evaluate the strategy
	signal, err := strategy.OnCandleUpdate(ctx, candles[0])
	if err != nil {
		return nil, fmt.Errorf("strategy evaluation failed: %w", err)
	}

	// Create a purchase decision based on the signal
	decision := &models.PurchaseDecision{
		Symbol:     symbol,
		Decision:   signal.Type == "BUY",
		Reason:     fmt.Sprintf("Strategy %s signal: %s", strategy.GetName(), signal.Type),
		Strategy:   strategy.GetName(),
		Confidence: signal.Confidence,
		Timestamp:  signal.Timestamp,
	}

	return decision, nil
}

// ExecuteStrategySignal executes a trading signal
func (s *tradeService) ExecuteStrategySignal(ctx context.Context, signal *strategy.Signal) (interface{}, error) {
	switch signal.Type {
	case "BUY":
		// Execute a buy order
		options := &models.PurchaseOptions{
			StopLossPercent: 0.05, // Default stop loss
		}

		// Use signal's stop loss if available
		if signal.StopLoss > 0 {
			options.StopLossPercent = (signal.Price - signal.StopLoss) / signal.Price
		}

		// Calculate quantity based on recommended size
		quantity := signal.RecommendedSize
		if quantity <= 0 {
			quantity = 0.1 // Default to 10% of available funds
		}

		// Execute the purchase
		return s.ExecutePurchase(ctx, signal.Symbol, quantity, options)

	case "SELL":
		// Find the coin to sell
		coin, err := s.boughtCoinRepo.FindBySymbol(ctx, signal.Symbol)
		if err != nil {
			return nil, fmt.Errorf("failed to find coin %s: %w", signal.Symbol, err)
		}

		if coin == nil {
			return nil, fmt.Errorf("no position found for %s", signal.Symbol)
		}

		// Sell the coin
		return s.SellCoin(ctx, coin, coin.Quantity)

	default:
		return nil, fmt.Errorf("unsupported signal type: %s", signal.Type)
	}
}

// prepareMarketData prepares market data for strategy evaluation
func (s *tradeService) prepareMarketData(ctx context.Context, symbol string) (*strategy.MarketData, error) {
	// Get candles for regime detection
	candles, err := s.getCandles(ctx, symbol, "1h", 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get candles: %w", err)
	}

	// Detect market regime
	regimeResult, err := strategy.DetectMarketRegime(ctx, candles)
	if err != nil {
		return nil, fmt.Errorf("failed to detect market regime: %w", err)
	}

	// Create market data
	marketData := &strategy.MarketData{
		Symbol: symbol,
		Regime: regimeResult,
	}

	return marketData, nil
}

// getCandles gets candles for a symbol and converts them to the models.Candle format
func (s *tradeService) getCandles(ctx context.Context, symbol, interval string, limit int) ([]*models.Candle, error) {
	// Get klines from the exchange
	klines, err := s.mexcClient.GetKlines(ctx, symbol, interval, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}

	// Convert klines to candles
	candles := make([]*models.Candle, len(klines))
	for i, k := range klines {
		candles[i] = &models.Candle{
			Symbol:     k.Symbol,
			Interval:   k.Interval,
			OpenTime:   k.OpenTime,
			CloseTime:  k.CloseTime,
			OpenPrice:  k.Open,
			HighPrice:  k.High,
			LowPrice:   k.Low,
			ClosePrice: k.Close,
			Volume:     k.Volume,
		}
	}

	return candles, nil
}

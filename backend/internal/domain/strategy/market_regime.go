package strategy

import (
	"context"
	"fmt"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/strategy/advanced"
)

// DetectMarketRegime detects the current market regime based on candles
func DetectMarketRegime(ctx context.Context, candles []*models.Candle) (string, error) {
	if len(candles) < 20 {
		return "UNKNOWN", fmt.Errorf("not enough candles to detect market regime, need at least 20, got %d", len(candles))
	}

	// Use the advanced regime detection
	result, err := advanced.DetectMarketRegime(candles)
	if err != nil {
		return "UNKNOWN", err
	}

	return string(result.Regime), nil
}

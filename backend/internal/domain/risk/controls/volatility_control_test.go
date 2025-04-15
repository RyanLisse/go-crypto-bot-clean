package controls

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
)

// MockMarketDataService implements port.MarketDataService for testing
type MockMarketDataService struct {
	mock.Mock
}

func (m *MockMarketDataService) GetAllSymbols(ctx context.Context) ([]*market.Symbol, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

func (m *MockMarketDataService) GetSymbol(ctx context.Context, symbol string) (*market.Symbol, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(*market.Symbol), args.Error(1)
}

func (m *MockMarketDataService) GetTicker(ctx context.Context, symbol string) (*market.Ticker, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(*market.Ticker), args.Error(1)
}

func (m *MockMarketDataService) GetCandles(ctx context.Context, symbol string, interval string, limit int) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, interval, limit)
	return args.Get(0).([]*market.Candle), args.Error(1)
}

func (m *MockMarketDataService) GetOrderBook(ctx context.Context, symbol string, limit int) (*market.OrderBook, error) {
	args := m.Called(ctx, symbol, limit)
	return args.Get(0).(*market.OrderBook), args.Error(1)
}

func (m *MockMarketDataService) GetHistoricalPrices(ctx context.Context, symbol string, startTime time.Time, endTime time.Time, interval string) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, startTime, endTime, interval)
	return args.Get(0).([]*market.Candle), args.Error(1)
}

func (m *MockMarketDataService) GetSymbolInfo(ctx context.Context, symbol string) (*market.Symbol, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).(*market.Symbol), args.Error(1)
}

func TestVolatilityControl_AssessRisk(t *testing.T) {
	// Setup
	testCases := []struct {
		name                string
		volatilityThreshold float64
		priceChanges        []float64
		expectedRiskLevel   model.RiskLevel
		expectNoRisk        bool
	}{
		{
			name:                "Volatility below threshold",
			volatilityThreshold: 0.05,                                      // 5%
			priceChanges:        []float64{0.01, -0.01, 0.02, -0.02, 0.01}, // Low volatility
			expectNoRisk:        true,
		},
		{
			name:                "Volatility slightly above threshold - medium risk",
			volatilityThreshold: 0.02,                                      // 2%
			priceChanges:        []float64{0.01, -0.01, 0.02, -0.02, 0.01}, // ~2.2% volatility
			expectedRiskLevel:   model.RiskLevelMedium,
			expectNoRisk:        false,
		},
		{
			name:                "Volatility significantly above threshold - high risk",
			volatilityThreshold: 0.02,                                      // 2%
			priceChanges:        []float64{0.03, -0.03, 0.04, -0.04, 0.03}, // ~3.5% volatility
			expectedRiskLevel:   model.RiskLevelHigh,
			expectNoRisk:        false,
		},
		{
			name:                "Volatility extremely high - critical risk",
			volatilityThreshold: 0.02,                                      // 2%
			priceChanges:        []float64{0.05, -0.05, 0.06, -0.06, 0.05}, // >4% volatility
			expectedRiskLevel:   model.RiskLevelCritical,
			expectNoRisk:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock market data with the specified price changes
			marketData := createMockMarketData("BTC-USD", tc.priceChanges)

			// Create mock market data service
			mockService := new(MockMarketDataService)

			// Create a mock risk profile with our test threshold
			profile := &model.RiskProfile{
				UserID:              "test-user",
				VolatilityThreshold: tc.volatilityThreshold,
			}

			// Create the control
			control := NewVolatilityControl(mockService)

			// Call the method under test - using the AssessRisk method which is a simplified version for testing
			assessment, err := control.AssessRisk("test-user", marketData, profile)

			// Assert no error
			require.NoError(t, err)

			if tc.expectNoRisk {
				assert.Nil(t, assessment, "Expected no risk assessment to be created")
			} else {
				// Verify the assessment was created correctly
				require.NotNil(t, assessment, "Expected a risk assessment to be created")
				assert.Equal(t, model.RiskTypeVolatility, assessment.Type)
				assert.Equal(t, tc.expectedRiskLevel, assessment.Level)
				assert.Equal(t, "test-user", assessment.UserID)
				assert.Equal(t, "BTC-USD", assessment.Symbol)
				assert.Contains(t, assessment.Message, "exceeding threshold")
			}
		})
	}
}

// Helper function to create mock market data with specified price changes
func createMockMarketData(symbol string, priceChanges []float64) market.Data {
	// Start with a base price
	basePrice := 10000.0

	// Create price history based on the specified changes
	var klines []market.Kline
	timestamp := time.Now().Add(-time.Hour * 24) // Start 24 hours ago

	for i, change := range priceChanges {
		// Calculate the new price
		price := basePrice * (1 + change)

		// Create a kline for this price
		kline := market.Kline{
			Symbol:    symbol,
			Interval:  market.KlineInterval15m,
			OpenTime:  timestamp.Add(time.Hour * time.Duration(i)),
			CloseTime: timestamp.Add(time.Hour * time.Duration(i+1)),
			Open:      basePrice,
			Close:     price,
			High:      basePrice * 1.01,
			Low:       basePrice * 0.99,
			Volume:    1000.0,
		}

		klines = append(klines, kline)

		// Update the base price for the next iteration
		basePrice = price
	}

	// Create a ticker for current price
	ticker := market.Ticker{
		Symbol: symbol,
		Price:  basePrice,
	}

	// Return the market data
	return market.Data{
		Symbol:      symbol,
		CurrentData: ticker,
		HistoricalData: market.HistoricalData{
			Klines: klines,
		},
	}
}

package service

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test setup helper
func setupPositionServiceTest() (*PositionService, *MockPositionUseCase, *MockMarketDataService) {
	mockPositionUC := new(MockPositionUseCase)
	mockMarketService := new(MockMarketDataService)
	logger := zerolog.Nop()

	service := NewPositionService(mockPositionUC, mockMarketService, &logger)
	return service, mockPositionUC, mockMarketService
}

func TestCreatePosition(t *testing.T) {
	// Setup
	service, mockPositionUC, _ := setupPositionServiceTest()
	ctx := context.Background()

	// Test data
	req := model.PositionCreateRequest{
		Symbol:     "BTC-USDT",
		Side:       model.PositionSideLong,
		Type:       model.PositionTypeManual,
		EntryPrice: 30000.0,
		Quantity:   0.5,
		OrderIDs:   []string{"order123"},
	}

	expectedPosition := &model.Position{
		ID:         "pos123",
		Symbol:     "BTC-USDT",
		Side:       model.PositionSideLong,
		Status:     model.PositionStatusOpen,
		Type:       model.PositionTypeManual,
		EntryPrice: 30000.0,
		Quantity:   0.5,
		OpenedAt:   time.Now(),
	}

	// Expectations
	mockPositionUC.On("CreatePosition", ctx, req).Return(expectedPosition, nil)

	// Execute
	position, err := service.CreatePosition(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPosition, position)
	mockPositionUC.AssertExpectations(t)
}

func TestUpdatePositionWithMarketData(t *testing.T) {
	// Setup
	service, mockPositionUC, mockMarketService := setupPositionServiceTest()
	ctx := context.Background()
	positionID := "pos123"

	// Test data
	existingPosition := &model.Position{
		ID:         positionID,
		Symbol:     "BTC-USDT",
		Side:       model.PositionSideLong,
		Status:     model.PositionStatusOpen,
		EntryPrice: 30000.0,
		Quantity:   0.5,
	}

	currentTicker := &market.Ticker{
		Symbol: "BTC-USDT",
		Price:  32000.0,
	}

	updatedPosition := &model.Position{
		ID:           positionID,
		Symbol:       "BTC-USDT",
		Side:         model.PositionSideLong,
		Status:       model.PositionStatusOpen,
		EntryPrice:   30000.0,
		Quantity:     0.5,
		CurrentPrice: 32000.0,
		PnL:          1000.0, // (32000-30000)*0.5
		PnLPercent:   6.67,   // (32000-30000)/30000*100
	}

	// Expectations
	mockPositionUC.On("GetPositionByID", ctx, positionID).Return(existingPosition, nil)
	mockMarketService.On("RefreshTicker", ctx, existingPosition.Symbol).Return(currentTicker, nil)
	mockPositionUC.On("UpdatePositionPrice", ctx, positionID, currentTicker.Price).Return(updatedPosition, nil)

	// Execute
	position, err := service.UpdatePositionWithMarketData(ctx, positionID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, updatedPosition, position)
	mockPositionUC.AssertExpectations(t)
	mockMarketService.AssertExpectations(t)
}

func TestClosePosition(t *testing.T) {
	// Setup
	service, mockPositionUC, mockMarketService := setupPositionServiceTest()
	ctx := context.Background()
	positionID := "pos123"
	exitOrderIDs := []string{"order456"}

	// Test data
	existingPosition := &model.Position{
		ID:         positionID,
		Symbol:     "BTC-USDT",
		Side:       model.PositionSideLong,
		Status:     model.PositionStatusOpen,
		EntryPrice: 30000.0,
		Quantity:   0.5,
	}

	currentTicker := &market.Ticker{
		Symbol: "BTC-USDT",
		Price:  32000.0,
	}

	now := time.Now()
	closedPosition := &model.Position{
		ID:           positionID,
		Symbol:       "BTC-USDT",
		Side:         model.PositionSideLong,
		Status:       model.PositionStatusClosed,
		EntryPrice:   30000.0,
		Quantity:     0.5,
		CurrentPrice: 32000.0,
		PnL:          1000.0,
		PnLPercent:   6.67,
		ExitOrderIDs: exitOrderIDs,
		ClosedAt:     &now,
	}

	// Expectations
	mockPositionUC.On("GetPositionByID", ctx, positionID).Return(existingPosition, nil)
	mockMarketService.On("RefreshTicker", ctx, existingPosition.Symbol).Return(currentTicker, nil)
	mockPositionUC.On("ClosePosition", ctx, positionID, currentTicker.Price, exitOrderIDs).Return(closedPosition, nil)

	// Execute
	position, err := service.ClosePosition(ctx, positionID, exitOrderIDs)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, closedPosition, position)
	mockPositionUC.AssertExpectations(t)
	mockMarketService.AssertExpectations(t)
}

func TestAnalyzePositionPerformance(t *testing.T) {
	// Setup
	service, mockPositionUC, mockMarketService := setupPositionServiceTest()
	ctx := context.Background()
	positionID := "pos123"

	// Test data
	now := time.Now()
	openedAt := now.Add(-24 * time.Hour)

	position := &model.Position{
		ID:              positionID,
		Symbol:          "BTC-USDT",
		Side:            model.PositionSideLong,
		Status:          model.PositionStatusOpen,
		Type:            model.PositionTypeManual,
		EntryPrice:      30000.0,
		Quantity:        0.5,
		CurrentPrice:    32000.0,
		PnL:             1000.0,
		PnLPercent:      6.67,
		MaxDrawdown:     -200.0,
		MaxProfit:       1200.0,
		RiskRewardRatio: 2.5,
		OpenedAt:        openedAt,
	}

	// Historical ticker data
	tickers := []market.Ticker{
		{Symbol: "BTC-USDT", Price: 30000.0},
		{Symbol: "BTC-USDT", Price: 29000.0}, // Lowest
		{Symbol: "BTC-USDT", Price: 33000.0}, // Highest
		{Symbol: "BTC-USDT", Price: 32000.0},
	}

	// Expectations
	mockPositionUC.On("GetPositionByID", ctx, positionID).Return(position, nil)
	mockMarketService.On("GetHistoricalPrices", ctx, position.Symbol, position.OpenedAt, mock.AnythingOfType("time.Time")).Return(tickers, nil)

	// Execute
	analysis, err := service.AnalyzePositionPerformance(ctx, positionID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, positionID, analysis["positionId"])
	assert.Equal(t, "BTC-USDT", analysis["symbol"])
	assert.Equal(t, model.PositionSideLong, analysis["side"])
	assert.Equal(t, model.PositionTypeManual, analysis["type"])
	assert.Equal(t, 30000.0, analysis["entryPrice"])
	assert.Equal(t, 32000.0, analysis["currentPrice"])
	assert.Equal(t, 0.5, analysis["quantity"])
	assert.Equal(t, 1000.0, analysis["pnl"])
	assert.Equal(t, 6.67, analysis["pnlPercent"])
	assert.Equal(t, -200.0, analysis["maxDrawdown"])
	assert.Equal(t, 1200.0, analysis["maxProfit"])
	assert.InDelta(t, 24.0, analysis["durationHours"], 0.01)
	assert.Equal(t, 2.5, analysis["riskRewardRatio"])

	// Volatility calculations
	assert.Equal(t, 33000.0, analysis["marketHighest"])
	assert.Equal(t, 29000.0, analysis["marketLowest"])
	assert.InDelta(t, 13.79, analysis["marketVolatility"], 0.01) // (33000-29000)/29000*100 = 13.79

	mockPositionUC.AssertExpectations(t)
	mockMarketService.AssertExpectations(t)
}

func TestGetOpenPositionsSummary(t *testing.T) {
	// Setup
	service, mockPositionUC, _ := setupPositionServiceTest()
	ctx := context.Background()

	// Test data
	positions := []*model.Position{
		{
			ID:         "pos1",
			Symbol:     "BTC-USDT",
			Side:       model.PositionSideLong,
			Type:       model.PositionTypeManual,
			EntryPrice: 30000.0,
			Quantity:   0.5,
			PnL:        1000.0,
		},
		{
			ID:         "pos2",
			Symbol:     "ETH-USDT",
			Side:       model.PositionSideLong,
			Type:       model.PositionTypeManual,
			EntryPrice: 2000.0,
			Quantity:   5.0,
			PnL:        500.0,
		},
		{
			ID:         "pos3",
			Symbol:     "SOL-USDT",
			Side:       model.PositionSideShort,
			Type:       model.PositionTypeAutomatic,
			EntryPrice: 100.0,
			Quantity:   20.0,
			PnL:        -200.0,
		},
	}

	// Expectations
	mockPositionUC.On("GetOpenPositions", ctx).Return(positions, nil)

	// Execute
	summary, err := service.GetOpenPositionsSummary(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 3, summary["totalOpenPositions"])
	assert.Equal(t, 27000.0, summary["totalValue"]) // 30000*0.5 + 2000*5 + 100*20 = 15000 + 10000 + 2000 = 27000
	assert.Equal(t, 1300.0, summary["totalPnL"])    // 1000 + 500 - 200
	assert.Equal(t, 2, summary["longPositions"])
	assert.Equal(t, 1, summary["shortPositions"])
	assert.InDelta(t, 433.33, summary["averagePnL"], 0.01) // 1300/3 = 433.33333...

	positionsByType := summary["positionsByType"].(map[string]int)
	assert.Equal(t, 2, positionsByType["MANUAL"])
	assert.Equal(t, 1, positionsByType["AUTOMATIC"])
	assert.Equal(t, 0, positionsByType["NEWCOIN"])

	mockPositionUC.AssertExpectations(t)
}

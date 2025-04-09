package backtest

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup common variables for tests
var (
	testDB         *gorm.DB
	testLogger     *zap.Logger
	testConfig     *RiskManagerConfig
	initialCapital = 10000.0
)

// setupTestDB initializes an in-memory SQLite database for testing
func setupTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	testDB = db
	// Migrate the schema
	err = db.AutoMigrate(&models.TakeProfitLevel{})
	assert.NoError(t, err)
}

// setupTestLogger initializes a Nop logger for testing
func setupTestLogger(t *testing.T) {
	testLogger = zap.NewNop()
}

// setupTestConfig initializes a default RiskManagerConfig for testing
func setupTestConfig(t *testing.T) {
	testConfig = &RiskManagerConfig{
		MaxRiskPerTrade:     1.0,  // 1%
		MaxPositionSize:     10.0, // 10%
		MaxTotalExposure:    50.0, // 50%
		MaxPositions:        5,
		MaxDrawdown:         20.0, // 20%
		MaxDailyLoss:        5.0,  // 5%
		UseTrailingStops:    true,
		TrailingStopPercent: 2.0, // 2%
		UseTakeProfits:      true,
		TakeProfitLevels:    []float64{2.0, 4.0, 6.0},    // 2%, 4%, 6%
		TakeProfitSizes:     []float64{30.0, 30.0, 40.0}, // 30%, 30%, 40%
		Logger:              testLogger,
		DB:                  testDB,
	}
}

func TestNewEventDrivenRiskManager(t *testing.T) {
	setupTestDB(t)
	setupTestLogger(t)
	setupTestConfig(t)

	rm := NewEventDrivenRiskManager(testConfig, initialCapital, testDB)
	assert.NotNil(t, rm)
	assert.Equal(t, testConfig, rm.config)
	assert.Equal(t, initialCapital, rm.initialCapital)
	assert.Equal(t, initialCapital, rm.currentCapital)
	assert.NotNil(t, rm.openPositions)
	assert.NotNil(t, rm.dailyPnL)
	assert.Equal(t, 0.0, rm.maxDrawdown)
	assert.Equal(t, 0.0, rm.currentDrawdown)
	assert.Equal(t, initialCapital, rm.highWaterMark)
	assert.Equal(t, testLogger, rm.logger)
	assert.Equal(t, testDB, rm.db)
	assert.NotNil(t, rm.trailingStops)
	assert.NotNil(t, rm.takeProfitLevels)
	assert.NotNil(t, rm.positionSizeCache)
}

func TestEventDrivenRiskManager_OnPositionOpened(t *testing.T) {
	setupTestDB(t)
	setupTestLogger(t)
	setupTestConfig(t)

	rm := NewEventDrivenRiskManager(testConfig, initialCapital, testDB)
	position := &models.Position{
		ID:         "pos1",
		Symbol:     "BTCUSDT",
		EntryPrice: 50000,
		Quantity:   0.1,
		Side:       models.OrderSideBuy,
		OpenTime:   time.Now(),
	}

	err := rm.OnPositionOpened(context.Background(), position)
	assert.NoError(t, err)
	assert.Contains(t, rm.openPositions, "pos1")
	assert.Equal(t, position, rm.openPositions["pos1"])

	// Check trailing stop
	expectedTrailingStop := 50000 * (1.0 - 0.02)
	assert.Equal(t, expectedTrailingStop, rm.trailingStops["pos1"])

	// Check take profit levels
	assert.Len(t, rm.takeProfitLevels["pos1"], 3)
	expectedTP1Price := 50000 * (1.0 + 0.02)
	expectedTP1Qty := 0.1 * 0.30
	assert.Equal(t, expectedTP1Price, rm.takeProfitLevels["pos1"][0].Price)
	assert.Equal(t, expectedTP1Qty, rm.takeProfitLevels["pos1"][0].Quantity)

	// Check capital update
	expectedCapital := initialCapital - (50000 * 0.1)
	assert.Equal(t, expectedCapital, rm.currentCapital)
}

func TestEventDrivenRiskManager_CalculatePositionSize_ExceedMaxExposure(t *testing.T) {
	setupTestDB(t)
	setupTestLogger(t)
	setupTestConfig(t)

	rm := NewEventDrivenRiskManager(testConfig, initialCapital, testDB)
	// Open a position that uses 45% of capital
	rm.openPositions["pos1"] = &models.Position{EntryPrice: 100, Quantity: 45}
	rm.currentCapital = initialCapital // Reset for calc

	// Try to open another position that would exceed 50% total exposure
	size, err := rm.CalculatePositionSize(context.Background(), "ETHUSDT", 2000, 1900)
	assert.NoError(t, err)

	// Expected size should be limited by max exposure (50% total = 5000)
	// Current exposure = 4500. Remaining exposure = 500.
	// Max size = 500 / 2000 = 0.25
	assert.InDelta(t, 0.25, size, 0.0001)
}

func TestEventDrivenRiskManager_OnPriceUpdate_TakeProfit(t *testing.T) {
	setupTestDB(t)
	setupTestLogger(t)
	setupTestConfig(t)

	rm := NewEventDrivenRiskManager(testConfig, initialCapital, testDB)
	position := &models.Position{
		ID:         "pos1",
		Symbol:     "BTCUSDT",
		EntryPrice: 50000,
		Quantity:   0.1,
		Side:       models.OrderSideBuy,
		OpenTime:   time.Now(),
	}
	rm.OnPositionOpened(context.Background(), position)

	// Price update triggers first take profit level (50000 * 1.02 = 51000)
	signals, err := rm.OnPriceUpdate(context.Background(), "BTCUSDT", 51050, time.Now())
	assert.NoError(t, err)
	assert.Len(t, signals, 1)
	assert.Equal(t, models.OrderSideSell, signals[0].Side)
	assert.Equal(t, 0.1*0.30, signals[0].Quantity) // 30% of position
	assert.True(t, rm.takeProfitLevels["pos1"][0].Triggered)
}

// func TestRiskManagerDrawdown(t *testing.T) {
// 	// ... (existing code) ...
// }

// // TestRiskManager_ShouldBuy tests the ShouldBuy method
// func TestRiskManager_ShouldBuy(t *testing.T) {
// 	setupTestDB(t)
// 	setupTestLogger(t)
// 	setupTestConfig(t)

// 	rm := NewEventDrivenRiskManager(testConfig, initialCapital, testDB)

// 	// Test case 1: Max drawdown exceeded
// 	rm.currentDrawdown = 25.0 // Exceeds 20% limit
// 	shouldBuy, reason := rm.ShouldBuy(context.Background(), "BTCUSDT")
// 	assert.False(t, shouldBuy)
// 	assert.Contains(t, reason, "drawdown")
// 	rm.currentDrawdown = 0 // Reset

// 	// Test case 2: Max daily loss exceeded
// 	dateStr := time.Now().Format("2006-01-02")
// 	rm.dailyPnL[dateStr] = -600 // Exceeds 5% (500) limit
// 	shouldBuy, reason = rm.ShouldBuy(context.Background(), "BTCUSDT")
// 	assert.False(t, shouldBuy)
// 	assert.Contains(t, reason, "daily loss")
// 	delete(rm.dailyPnL, dateStr) // Reset

// 	// Test case 3: Max positions reached
// 	for i := 0; i < testConfig.MaxPositions; i++ {
// 		rm.openPositions[fmt.Sprintf("pos%d", i)] = &models.Position{}
// 	}
// 	shouldBuy, reason = rm.ShouldBuy(context.Background(), "BTCUSDT")
// 	assert.False(t, shouldBuy)
// 	assert.Contains(t, reason, "maximum number of positions")
// 	rm.openPositions = make(map[string]*models.Position) // Reset

// 	// Test case 4: Okay to buy
// 	shouldBuy, reason = rm.ShouldBuy(context.Background(), "BTCUSDT")
// 	assert.True(t, shouldBuy)
// 	assert.Equal(t, "", reason)
// }

// // TestRiskManager_UpdatePosition tests the UpdatePosition method
// func TestRiskManager_UpdatePosition(t *testing.T) {
// 	setupTestDB(t)
// 	setupTestLogger(t)
// 	setupTestConfig(t)

// 	rm := NewEventDrivenRiskManager(testConfig, initialCapital, testDB)
// 	position := &models.Position{
// 		ID:         "pos1",
// 		Symbol:     "BTCUSDT",
// 		EntryPrice: 50000,
// 		Quantity:   0.1,
// 		Side:       models.OrderSideBuy,
// 		OpenTime:   time.Now(),
// 	}
// 	rm.OnPositionOpened(context.Background(), position)

// 	// Simulate partial close
// 	remainingQty := 0.07
// 	err := rm.UpdatePosition(context.Background(), "pos1", remainingQty)
// 	assert.NoError(t, err)
// 	assert.Equal(t, remainingQty, rm.openPositions["pos1"].Quantity)

// 	// Simulate closing position fully
// 	err = rm.UpdatePosition(context.Background(), "pos1", 0)
// 	assert.Error(t, err) // Should error as closing is done via OnPositionClosed
// 	// assert.NotContains(t, rm.openPositions, "pos1") // Position shouldn't be removed here
// }

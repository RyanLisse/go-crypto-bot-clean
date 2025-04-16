package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// PositionMockRepository is a mock for position repository
type PositionMockRepository struct {
	mock.Mock
}

func (m *PositionMockRepository) Create(ctx context.Context, position *model.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

func (m *PositionMockRepository) GetByID(ctx context.Context, id string) (*model.Position, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *PositionMockRepository) Update(ctx context.Context, position *model.Position) error {
	args := m.Called(ctx, position)
	return args.Error(0)
}

func (m *PositionMockRepository) GetOpenPositions(ctx context.Context) ([]*model.Position, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetOpenPositionsBySymbol(ctx context.Context, symbol string) ([]*model.Position, error) {
	args := m.Called(ctx, symbol)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error) {
	args := m.Called(ctx, positionType)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, symbol, limit, offset)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) GetClosedPositions(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, from, to, limit, offset)
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *PositionMockRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *PositionMockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetActiveByUser returns active positions for a user
func (m *PositionMockRepository) GetActiveByUser(ctx context.Context, userID string) ([]*model.Position, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Position), args.Error(1)
}

// GetBySymbolAndUser retrieves positions for a specific symbol and user with pagination
func (m *PositionMockRepository) GetBySymbolAndUser(ctx context.Context, symbol, userID string, page, limit int) ([]*model.Position, error) {
	args := m.Called(ctx, symbol, userID, page, limit)
	return args.Get(0).([]*model.Position), args.Error(1)
}

// GetOpenPositionsByUserID retrieves all open positions for a specific user
func (m *PositionMockRepository) GetOpenPositionsByUserID(ctx context.Context, userID string) ([]*model.Position, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Position), args.Error(1)
}

// PositionMockMarketRepository is a mock for market repository
type PositionMockMarketRepository struct {
	mock.Mock
}

func (m *PositionMockMarketRepository) GetOrderBook(ctx context.Context, symbol, exchange string, depth int) (*market.OrderBook, error) {
	args := m.Called(ctx, symbol, exchange, depth)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.OrderBook), args.Error(1)
}

func (m *PositionMockMarketRepository) GetTicker(ctx context.Context, symbol, exchange string) (*market.Ticker, error) {
	args := m.Called(ctx, symbol, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Ticker), args.Error(1)
}

// Add missing methods for MarketRepository interface
func (m *PositionMockMarketRepository) GetAllTickers(ctx context.Context, exchange string) ([]*market.Ticker, error) {
	args := m.Called(ctx, exchange)
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

func (m *PositionMockMarketRepository) GetTickerHistory(ctx context.Context, symbol, exchange string, start, end time.Time) ([]*market.Ticker, error) {
	args := m.Called(ctx, symbol, exchange, start, end)
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

func (m *PositionMockMarketRepository) SaveTicker(ctx context.Context, ticker *market.Ticker) error {
	args := m.Called(ctx, ticker)
	return args.Error(0)
}

func (m *PositionMockMarketRepository) SaveCandle(ctx context.Context, candle *market.Candle) error {
	args := m.Called(ctx, candle)
	return args.Error(0)
}

func (m *PositionMockMarketRepository) SaveCandles(ctx context.Context, candles []*market.Candle) error {
	args := m.Called(ctx, candles)
	return args.Error(0)
}

func (m *PositionMockMarketRepository) GetCandle(ctx context.Context, symbol, exchange string, interval market.Interval, openTime time.Time) (*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval, openTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Candle), args.Error(1)
}

func (m *PositionMockMarketRepository) GetCandles(ctx context.Context, symbol, exchange string, interval market.Interval, start, end time.Time, limit int) ([]*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval, start, end, limit)
	return args.Get(0).([]*market.Candle), args.Error(1)
}

func (m *PositionMockMarketRepository) GetLatestCandle(ctx context.Context, symbol, exchange string, interval market.Interval) (*market.Candle, error) {
	args := m.Called(ctx, symbol, exchange, interval)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Candle), args.Error(1)
}

func (m *PositionMockMarketRepository) SaveOrderBook(ctx context.Context, orderBook *market.OrderBook) error {
	args := m.Called(ctx, orderBook)
	return args.Error(0)
}

func (m *PositionMockMarketRepository) PurgeOldData(ctx context.Context, olderThan time.Time) error {
	args := m.Called(ctx, olderThan)
	return args.Error(0)
}

func (m *PositionMockMarketRepository) GetLatestTickers(ctx context.Context, limit int) ([]*market.Ticker, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

func (m *PositionMockMarketRepository) GetTickersBySymbol(ctx context.Context, symbol string, limit int) ([]*market.Ticker, error) {
	args := m.Called(ctx, symbol, limit)
	return args.Get(0).([]*market.Ticker), args.Error(1)
}

// PositionMockSymbolRepository is a mock for symbol repository
type PositionMockSymbolRepository struct {
	mock.Mock
}

func (m *PositionMockSymbolRepository) GetBySymbol(ctx context.Context, symbol string) (*market.Symbol, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*market.Symbol), args.Error(1)
}

// Add missing methods for SymbolRepository interface
func (m *PositionMockSymbolRepository) Create(ctx context.Context, symbol *market.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

func (m *PositionMockSymbolRepository) GetByExchange(ctx context.Context, exchange string) ([]*market.Symbol, error) {
	args := m.Called(ctx, exchange)
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

func (m *PositionMockSymbolRepository) GetAll(ctx context.Context) ([]*market.Symbol, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*market.Symbol), args.Error(1)
}

func (m *PositionMockSymbolRepository) Update(ctx context.Context, symbol *market.Symbol) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

func (m *PositionMockSymbolRepository) Delete(ctx context.Context, symbol string) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

// Test setup helper
func setupPositionUseCase(
	positionRepo *PositionMockRepository,
	marketRepo *PositionMockMarketRepository,
	symbolRepo *PositionMockSymbolRepository,
) usecase.PositionUseCase {
	// Create a null logger for testing
	logger := zerolog.Nop()
	return usecase.NewPositionUseCase(positionRepo, marketRepo, symbolRepo, logger)
}

// Tests
func TestCreatePosition(t *testing.T) {
	// Setup
	positionRepo := new(PositionMockRepository)
	marketRepo := new(PositionMockMarketRepository)
	symbolRepo := new(PositionMockSymbolRepository)
	positionUC := setupPositionUseCase(positionRepo, marketRepo, symbolRepo)
	ctx := context.Background()

	// Valid symbol
	symbol := &market.Symbol{
		Symbol:     "BTCUSDT",
		BaseAsset:  "BTC",
		QuoteAsset: "USDT",
	}

	// Input data
	createReq := model.PositionCreateRequest{
		Symbol:     "BTCUSDT",
		Side:       model.PositionSideLong,
		Type:       model.PositionTypeManual,
		EntryPrice: 50000.0,
		Quantity:   0.1,
		OrderIDs:   []string{"order1"},
		Notes:      "Test position",
	}

	t.Run("Success", func(t *testing.T) {
		// Setup expectations
		symbolRepo.On("GetBySymbol", ctx, "BTCUSDT").Return(symbol, nil).Once()
		positionRepo.On("Create", ctx, mock.MatchedBy(func(p *model.Position) bool {
			return p.Symbol == "BTCUSDT" && p.Side == model.PositionSideLong
		})).Return(nil).Once()

		// Execute
		position, err := positionUC.CreatePosition(ctx, createReq)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, position)
		assert.Equal(t, "BTCUSDT", position.Symbol)
		assert.Equal(t, model.PositionSideLong, position.Side)
		assert.Equal(t, model.PositionTypeManual, position.Type)
		assert.Equal(t, 50000.0, position.EntryPrice)
		assert.Equal(t, 0.1, position.Quantity)
		assert.Equal(t, model.PositionStatusOpen, position.Status)
		assert.Equal(t, []string{"order1"}, position.EntryOrderIDs)

		// Verify mocks
		symbolRepo.AssertExpectations(t)
		positionRepo.AssertExpectations(t)
	})

	t.Run("Symbol Not Found", func(t *testing.T) {
		// Setup expectations
		symbolRepo.On("GetBySymbol", ctx, "BTCUSDT").Return(nil, nil).Once()

		// Execute
		position, err := positionUC.CreatePosition(ctx, createReq)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, usecase.ErrSymbolNotFound, err)
		assert.Nil(t, position)

		// Verify mocks
		symbolRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		// Setup expectations
		symbolRepo.On("GetBySymbol", ctx, "BTCUSDT").Return(symbol, nil).Once()
		positionRepo.On("Create", ctx, mock.Anything).Return(errors.New("db error")).Once()

		// Execute
		position, err := positionUC.CreatePosition(ctx, createReq)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, position)

		// Verify mocks
		symbolRepo.AssertExpectations(t)
		positionRepo.AssertExpectations(t)
	})
}

func TestGetPositionByID(t *testing.T) {
	// Setup
	positionRepo := new(PositionMockRepository)
	marketRepo := new(PositionMockMarketRepository)
	symbolRepo := new(PositionMockSymbolRepository)
	positionUC := setupPositionUseCase(positionRepo, marketRepo, symbolRepo)
	ctx := context.Background()

	// Test data
	position := &model.Position{
		ID:           "pos1",
		Symbol:       "BTCUSDT",
		Side:         model.PositionSideLong,
		Status:       model.PositionStatusOpen,
		EntryPrice:   50000.0,
		Quantity:     0.1,
		CurrentPrice: 55000.0,
		PnL:          500.0,
	}

	t.Run("Success", func(t *testing.T) {
		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos1").Return(position, nil).Once()

		// Execute
		result, err := positionUC.GetPositionByID(ctx, "pos1")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, position, result)

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Setup expectations
		positionRepo.On("GetByID", ctx, "nonexistent").Return(nil, errors.New("not found")).Once()

		// Execute
		result, err := positionUC.GetPositionByID(ctx, "nonexistent")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})
}

func TestUpdatePositionPrice(t *testing.T) {
	// Setup
	positionRepo := new(PositionMockRepository)
	marketRepo := new(PositionMockMarketRepository)
	symbolRepo := new(PositionMockSymbolRepository)
	positionUC := setupPositionUseCase(positionRepo, marketRepo, symbolRepo)
	ctx := context.Background()

	// Test position
	now := time.Now()
	position := &model.Position{
		ID:            "pos1",
		Symbol:        "BTCUSDT",
		Side:          model.PositionSideLong,
		Status:        model.PositionStatusOpen,
		EntryPrice:    50000.0,
		Quantity:      0.1,
		CurrentPrice:  50000.0,
		PnL:           0.0,
		PnLPercent:    0.0,
		OpenedAt:      now,
		LastUpdatedAt: now,
	}

	t.Run("Update Long Position Price", func(t *testing.T) {
		// Clone position to avoid test interference
		pos := *position
		testPos := &pos

		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos1").Return(testPos, nil).Once()
		positionRepo.On("Update", ctx, mock.MatchedBy(func(p *model.Position) bool {
			return p.ID == "pos1" && p.CurrentPrice == 55000.0
		})).Return(nil).Once()

		// Execute
		result, err := positionUC.UpdatePositionPrice(ctx, "pos1", 55000.0)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 55000.0, result.CurrentPrice)
		assert.Equal(t, 500.0, result.PnL)       // (55000-50000)*0.1
		assert.Equal(t, 10.0, result.PnLPercent) // (55000-50000)/50000*100

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Update Short Position Price", func(t *testing.T) {
		// Create short position
		shortPos := &model.Position{
			ID:            "pos2",
			Symbol:        "BTCUSDT",
			Side:          model.PositionSideShort,
			Status:        model.PositionStatusOpen,
			EntryPrice:    50000.0,
			Quantity:      0.1,
			CurrentPrice:  50000.0,
			PnL:           0.0,
			PnLPercent:    0.0,
			OpenedAt:      now,
			LastUpdatedAt: now,
		}

		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos2").Return(shortPos, nil).Once()
		positionRepo.On("Update", ctx, mock.MatchedBy(func(p *model.Position) bool {
			return p.ID == "pos2" && p.CurrentPrice == 45000.0
		})).Return(nil).Once()

		// Execute
		result, err := positionUC.UpdatePositionPrice(ctx, "pos2", 45000.0)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 45000.0, result.CurrentPrice)
		assert.Equal(t, 500.0, result.PnL)       // (50000-45000)*0.1
		assert.Equal(t, 10.0, result.PnLPercent) // (50000-45000)/50000*100

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Position Not Found", func(t *testing.T) {
		// Setup expectations
		positionRepo.On("GetByID", ctx, "nonexistent").Return(nil, errors.New("not found")).Once()

		// Execute
		result, err := positionUC.UpdatePositionPrice(ctx, "nonexistent", 55000.0)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})
}

func TestClosePosition(t *testing.T) {
	// Setup
	positionRepo := new(PositionMockRepository)
	marketRepo := new(PositionMockMarketRepository)
	symbolRepo := new(PositionMockSymbolRepository)
	positionUC := setupPositionUseCase(positionRepo, marketRepo, symbolRepo)
	ctx := context.Background()

	// Test position
	now := time.Now()
	position := &model.Position{
		ID:            "pos1",
		Symbol:        "BTCUSDT",
		Side:          model.PositionSideLong,
		Status:        model.PositionStatusOpen,
		EntryPrice:    50000.0,
		Quantity:      0.1,
		CurrentPrice:  50000.0,
		PnL:           0.0,
		PnLPercent:    0.0,
		OpenedAt:      now,
		LastUpdatedAt: now,
	}

	t.Run("Close Position Successfully", func(t *testing.T) {
		// Clone position to avoid test interference
		pos := *position
		testPos := &pos

		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos1").Return(testPos, nil).Once()
		positionRepo.On("Update", ctx, mock.MatchedBy(func(p *model.Position) bool {
			return p.ID == "pos1" &&
				p.Status == model.PositionStatusClosed &&
				p.CurrentPrice == 55000.0 &&
				p.ExitOrderIDs[0] == "exit1"
		})).Return(nil).Once()

		// Execute
		result, err := positionUC.ClosePosition(ctx, "pos1", 55000.0, []string{"exit1"})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, model.PositionStatusClosed, result.Status)
		assert.Equal(t, 55000.0, result.CurrentPrice)
		assert.NotNil(t, result.ClosedAt)
		assert.Equal(t, []string{"exit1"}, result.ExitOrderIDs)
		assert.Equal(t, 500.0, result.PnL) // (55000-50000)*0.1

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Position Already Closed", func(t *testing.T) {
		// Create closed position
		closedTime := now.Add(-1 * time.Hour)
		closedPos := &model.Position{
			ID:            "pos2",
			Symbol:        "BTCUSDT",
			Side:          model.PositionSideLong,
			Status:        model.PositionStatusClosed,
			EntryPrice:    50000.0,
			Quantity:      0.1,
			CurrentPrice:  55000.0,
			PnL:           500.0,
			PnLPercent:    10.0,
			OpenedAt:      now.Add(-2 * time.Hour),
			ClosedAt:      &closedTime,
			LastUpdatedAt: closedTime,
			ExitOrderIDs:  []string{"existing"},
		}

		// Setup expectations
		positionRepo.On("GetByID", ctx, "pos2").Return(closedPos, nil).Once()

		// Execute - trying to close an already closed position
		result, err := positionUC.ClosePosition(ctx, "pos2", 60000.0, []string{"new"})

		// Assert
		assert.NoError(t, err) // This is not an error, just a no-op
		assert.NotNil(t, result)
		assert.Equal(t, model.PositionStatusClosed, result.Status)
		assert.Equal(t, 55000.0, result.CurrentPrice)              // Original price, not updated
		assert.Equal(t, []string{"existing"}, result.ExitOrderIDs) // Original orders not updated

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})

	t.Run("Position Not Found", func(t *testing.T) {
		// Setup expectations
		positionRepo.On("GetByID", ctx, "nonexistent").Return(nil, errors.New("not found")).Once()

		// Execute
		result, err := positionUC.ClosePosition(ctx, "nonexistent", 55000.0, []string{"exit1"})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify mocks
		positionRepo.AssertExpectations(t)
	})
}

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo/crypto-bot/internal/adapter/http/response"
	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPositionUseCase is a mock implementation of the PositionUseCase interface
type MockPositionUseCase struct {
	mock.Mock
}

func (m *MockPositionUseCase) CreatePosition(ctx context.Context, req model.PositionCreateRequest) (*model.Position, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetPositionByID(ctx context.Context, id string) (*model.Position, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetOpenPositions(ctx context.Context) ([]*model.Position, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetOpenPositionsBySymbol(ctx context.Context, symbol string) ([]*model.Position, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetOpenPositionsByType(ctx context.Context, positionType model.PositionType) ([]*model.Position, error) {
	args := m.Called(ctx, positionType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetPositionsBySymbol(ctx context.Context, symbol string, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, symbol, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) GetClosedPositions(ctx context.Context, from, to time.Time, limit, offset int) ([]*model.Position, error) {
	args := m.Called(ctx, from, to, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) UpdatePosition(ctx context.Context, id string, req model.PositionUpdateRequest) (*model.Position, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) UpdatePositionPrice(ctx context.Context, id string, currentPrice float64) (*model.Position, error) {
	args := m.Called(ctx, id, currentPrice)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) ClosePosition(ctx context.Context, id string, exitPrice float64, exitOrderIDs []string) (*model.Position, error) {
	args := m.Called(ctx, id, exitPrice, exitOrderIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) SetStopLoss(ctx context.Context, id string, stopLoss float64) (*model.Position, error) {
	args := m.Called(ctx, id, stopLoss)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) SetTakeProfit(ctx context.Context, id string, takeProfit float64) (*model.Position, error) {
	args := m.Called(ctx, id, takeProfit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Position), args.Error(1)
}

func (m *MockPositionUseCase) DeletePosition(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter(mockUseCase *MockPositionUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create a real logger that writes to a discard writer
	logger := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = io.Discard
		w.NoColor = true
	})).With().Timestamp().Logger()

	handler := NewPositionHandler(mockUseCase, &logger)

	api := router.Group("/api")
	handler.RegisterRoutes(api)

	return router
}

func TestCreatePosition(t *testing.T) {
	// Setup
	mockUseCase := new(MockPositionUseCase)
	router := setupTestRouter(mockUseCase)

	// Test data
	testPosition := &model.Position{
		ID:            "pos123",
		Symbol:        "BTCUSDT",
		Side:          model.PositionSideLong,
		Status:        model.PositionStatusOpen,
		Type:          model.PositionTypeManual,
		EntryPrice:    50000.0,
		Quantity:      1.0,
		CurrentPrice:  50000.0,
		PnL:           0.0,
		PnLPercent:    0.0,
		EntryOrderIDs: []string{"order123"},
		OpenedAt:      time.Now(),
		LastUpdatedAt: time.Now(),
	}

	createReq := model.PositionCreateRequest{
		Symbol:     "BTCUSDT",
		Side:       model.PositionSideLong,
		Type:       model.PositionTypeManual,
		EntryPrice: 50000.0,
		Quantity:   1.0,
		OrderIDs:   []string{"order123"},
	}

	reqBody, _ := json.Marshal(createReq)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("CreatePosition", mock.Anything, createReq).Return(testPosition, nil).Once()

		req, _ := http.NewRequest("POST", "/api/positions", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)

		var response response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Invalid Request", func(t *testing.T) {
		// Invalid request (missing required fields)
		invalidReq := map[string]interface{}{
			"symbol": "BTCUSDT",
			// Missing other required fields
		}
		invalidReqBody, _ := json.Marshal(invalidReq)

		req, _ := http.NewRequest("POST", "/api/positions", bytes.NewBuffer(invalidReqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, response.ErrorCodeBadRequest, resp.Error.Code)
	})

	t.Run("Symbol Not Found", func(t *testing.T) {
		mockUseCase.On("CreatePosition", mock.Anything, createReq).
			Return(nil, usecase.ErrSymbolNotFound).Once()

		req, _ := http.NewRequest("POST", "/api/positions", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, response.ErrorCodeBadRequest, resp.Error.Code)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Internal Error", func(t *testing.T) {
		mockUseCase.On("CreatePosition", mock.Anything, createReq).
			Return(nil, errors.New("internal error")).Once()

		req, _ := http.NewRequest("POST", "/api/positions", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, response.ErrorCodeInternalError, resp.Error.Code)

		mockUseCase.AssertExpectations(t)
	})
}

func TestGetOpenPositions(t *testing.T) {
	// Setup
	mockUseCase := new(MockPositionUseCase)
	router := setupTestRouter(mockUseCase)

	// Test data
	testPositions := []*model.Position{
		{
			ID:            "pos123",
			Symbol:        "BTCUSDT",
			Side:          model.PositionSideLong,
			Status:        model.PositionStatusOpen,
			Type:          model.PositionTypeManual,
			EntryPrice:    50000.0,
			Quantity:      1.0,
			CurrentPrice:  51000.0,
			PnL:           1000.0,
			PnLPercent:    2.0,
			EntryOrderIDs: []string{"order123"},
			OpenedAt:      time.Now(),
			LastUpdatedAt: time.Now(),
		},
		{
			ID:            "pos456",
			Symbol:        "ETHUSDT",
			Side:          model.PositionSideShort,
			Status:        model.PositionStatusOpen,
			Type:          model.PositionTypeAutomatic,
			EntryPrice:    3000.0,
			Quantity:      2.0,
			CurrentPrice:  2950.0,
			PnL:           100.0,
			PnLPercent:    1.67,
			EntryOrderIDs: []string{"order456"},
			OpenedAt:      time.Now(),
			LastUpdatedAt: time.Now(),
		},
	}

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("GetOpenPositions", mock.Anything).Return(testPositions, nil).Once()

		req, _ := http.NewRequest("GET", "/api/positions/open", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)

		// Verify that we got positions in the response
		positionsData, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 2, len(positionsData))

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase.On("GetOpenPositions", mock.Anything).
			Return(nil, errors.New("database error")).Once()

		req, _ := http.NewRequest("GET", "/api/positions/open", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, response.ErrorCodeInternalError, resp.Error.Code)

		mockUseCase.AssertExpectations(t)
	})
}

func TestGetPositionByID(t *testing.T) {
	// Setup
	mockUseCase := new(MockPositionUseCase)
	router := setupTestRouter(mockUseCase)

	// Test data
	testPosition := &model.Position{
		ID:            "pos123",
		Symbol:        "BTCUSDT",
		Side:          model.PositionSideLong,
		Status:        model.PositionStatusOpen,
		Type:          model.PositionTypeManual,
		EntryPrice:    50000.0,
		Quantity:      1.0,
		CurrentPrice:  51000.0,
		PnL:           1000.0,
		PnLPercent:    2.0,
		EntryOrderIDs: []string{"order123"},
		OpenedAt:      time.Now(),
		LastUpdatedAt: time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("GetPositionByID", mock.Anything, "pos123").Return(testPosition, nil).Once()

		req, _ := http.NewRequest("GET", "/api/positions/pos123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUseCase.On("GetPositionByID", mock.Anything, "nonexistent").
			Return(nil, usecase.ErrPositionNotFound).Once()

		req, _ := http.NewRequest("GET", "/api/positions/nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, response.ErrorCodeNotFound, resp.Error.Code)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Internal Error", func(t *testing.T) {
		mockUseCase.On("GetPositionByID", mock.Anything, "pos123").
			Return(nil, errors.New("database error")).Once()

		req, _ := http.NewRequest("GET", "/api/positions/pos123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, response.ErrorCodeInternalError, resp.Error.Code)

		mockUseCase.AssertExpectations(t)
	})
}

func TestUpdatePosition(t *testing.T) {
	// Setup
	mockUseCase := new(MockPositionUseCase)
	router := setupTestRouter(mockUseCase)

	// Test data
	updatedPosition := &model.Position{
		ID:            "pos123",
		Symbol:        "BTCUSDT",
		Side:          model.PositionSideLong,
		Status:        model.PositionStatusOpen,
		Type:          model.PositionTypeManual,
		EntryPrice:    50000.0,
		Quantity:      1.0,
		CurrentPrice:  52000.0,
		PnL:           2000.0,
		PnLPercent:    4.0,
		StopLoss:      new(float64),
		EntryOrderIDs: []string{"order123"},
		OpenedAt:      time.Now(),
		LastUpdatedAt: time.Now(),
	}
	*updatedPosition.StopLoss = 48000.0

	updateReq := model.PositionUpdateRequest{
		CurrentPrice: new(float64),
		StopLoss:     new(float64),
	}
	*updateReq.CurrentPrice = 52000.0
	*updateReq.StopLoss = 48000.0

	reqBody, _ := json.Marshal(updateReq)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("UpdatePosition", mock.Anything, "pos123", mock.MatchedBy(func(req model.PositionUpdateRequest) bool {
			return *req.CurrentPrice == 52000.0 && *req.StopLoss == 48000.0
		})).Return(updatedPosition, nil).Once()

		req, _ := http.NewRequest("PUT", "/api/positions/pos123", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUseCase.On("UpdatePosition", mock.Anything, "nonexistent", mock.Anything).
			Return(nil, usecase.ErrPositionNotFound).Once()

		req, _ := http.NewRequest("PUT", "/api/positions/nonexistent", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, response.ErrorCodeNotFound, resp.Error.Code)

		mockUseCase.AssertExpectations(t)
	})
}

func TestClosePosition(t *testing.T) {
	// Setup
	mockUseCase := new(MockPositionUseCase)
	router := setupTestRouter(mockUseCase)

	// Test data
	now := time.Now()
	closedPosition := &model.Position{
		ID:            "pos123",
		Symbol:        "BTCUSDT",
		Side:          model.PositionSideLong,
		Status:        model.PositionStatusClosed,
		Type:          model.PositionTypeManual,
		EntryPrice:    50000.0,
		Quantity:      1.0,
		CurrentPrice:  52000.0,
		PnL:           2000.0,
		PnLPercent:    4.0,
		EntryOrderIDs: []string{"order123"},
		ExitOrderIDs:  []string{"exit456"},
		OpenedAt:      now.Add(-24 * time.Hour),
		ClosedAt:      &now,
		LastUpdatedAt: now,
	}

	closeReq := struct {
		ExitPrice    float64  `json:"exitPrice"`
		ExitOrderIDs []string `json:"exitOrderIds"`
	}{
		ExitPrice:    52000.0,
		ExitOrderIDs: []string{"exit456"},
	}

	reqBody, _ := json.Marshal(closeReq)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("ClosePosition", mock.Anything, "pos123", 52000.0, []string{"exit456"}).
			Return(closedPosition, nil).Once()

		req, _ := http.NewRequest("PUT", "/api/positions/pos123/close", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUseCase.On("ClosePosition", mock.Anything, "nonexistent", 52000.0, []string{"exit456"}).
			Return(nil, usecase.ErrPositionNotFound).Once()

		req, _ := http.NewRequest("PUT", "/api/positions/nonexistent/close", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, response.ErrorCodeNotFound, resp.Error.Code)

		mockUseCase.AssertExpectations(t)
	})
}

func TestDeletePosition(t *testing.T) {
	// Setup
	mockUseCase := new(MockPositionUseCase)
	router := setupTestRouter(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("DeletePosition", mock.Anything, "pos123").Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/api/positions/pos123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNoContent, w.Code)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUseCase.On("DeletePosition", mock.Anything, "nonexistent").
			Return(usecase.ErrPositionNotFound).Once()

		req, _ := http.NewRequest("DELETE", "/api/positions/nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp response.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Equal(t, response.ErrorCodeNotFound, resp.Error.Code)

		mockUseCase.AssertExpectations(t)
	})
}

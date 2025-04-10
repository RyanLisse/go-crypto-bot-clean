// Package api provides API server initialization.
package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/api/middleware"
	"go-crypto-bot-clean/backend/internal/api/websocket"
	"go-crypto-bot-clean/backend/internal/config"
	"go-crypto-bot-clean/backend/internal/core/analytics"
	"go-crypto-bot-clean/backend/internal/core/newcoin"
	"go-crypto-bot-clean/backend/internal/core/status"
	"go-crypto-bot-clean/backend/internal/core/trade"
	"go-crypto-bot-clean/backend/internal/domain/audit"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/security"

	"go.uber.org/zap"
)

// Helper function to convert time.Time to *time.Time
func toTimePtr(t time.Time) *time.Time {
	return &t
}

// ServerDependencies holds references to services and repositories.
type ServerDependencies struct {
	// Core services
	TradeService     trade.TradeService
	NewCoinService   newcoin.NewCoinService
	AnalyticsService analytics.TradeAnalyticsService

	// Handlers
	HealthHandler    *handlers.HealthHandler
	StatusHandler    *handlers.StatusHandler
	WebSocketHandler *websocket.Handler
	PortfolioHandler *handlers.PortfolioHandler
	TradeHandler     *handlers.TradeHandler
	NewCoinHandler   *handlers.NewCoinsHandler // Handler for new coin detection
	CoinHandler      *handlers.CoinHandler     // Handler for market and tradable coin endpoints
	ConfigHandler    *handlers.ConfigHandler
	AnalyticsHandler *handlers.AnalyticsHandler

	// Middleware dependencies
	Logger       middleware.Logger
	ValidAPIKeys map[string]struct{}
	RateLimit    struct {
		Rate     float64
		Capacity float64
	}

	// Security dependencies
	AuditService      audit.Service
	EncryptionService security.EncryptionService
	ZapLogger         *zap.Logger
}

// NewServerDependencies initializes server dependencies.
func NewServerDependencies() *ServerDependencies {
	// Initialize logger
	logger := &mockLogger{}
	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()

	// Create a mock status provider for now
	statusProvider := &mockStatusProvider{}
	statusService := status.NewStatusService(statusProvider, "1.0.0")
	statusHandler := handlers.NewStatusHandler(statusService)

	// Create WebSocket hub and handler
	hub := websocket.NewHub()
	go hub.Run()

	// Create mock services
	mockTradeService := &mockTradeService{}
	mockNewCoinService := &mockNewCoinService{}
	mockBoughtCoinRepo := &mockBoughtCoinRepository{}
	mockPortfolioService := &mockPortfolioService{}
	mockExchangeService := &mockExchangeService{}
	webSocketHandler := websocket.NewHandler(hub, mockTradeService, mockNewCoinService, mockExchangeService)

	// Create actual handlers
	portfolioHandler := handlers.NewPortfolioHandler(mockPortfolioService)

	tradeHandler := handlers.NewTradeHandler(mockTradeService, mockBoughtCoinRepo)

	newCoinHandler := handlers.NewNewCoinsHandler(mockNewCoinService)
	coinHandler := handlers.NewCoinHandler(mockExchangeService, mockNewCoinService)

	configHandler := handlers.NewConfigHandler(&config.Config{
		Trading: config.TradingConfig{
			DefaultQuantity:  20.0,
			StopLossPercent:  10.0,
			TakeProfitLevels: []float64{5.0, 10.0, 15.0, 20.0},
			SellPercentages:  []float64{0.25, 0.25, 0.25, 0.25},
		},
	})

	return &ServerDependencies{
		TradeService:     mockTradeService,
		NewCoinService:   mockNewCoinService,
		HealthHandler:    healthHandler,
		StatusHandler:    statusHandler,
		WebSocketHandler: webSocketHandler,
		PortfolioHandler: portfolioHandler,
		TradeHandler:     tradeHandler,
		NewCoinHandler:   newCoinHandler,
		CoinHandler:      coinHandler,
		ConfigHandler:    configHandler,
		Logger:           logger,
		ValidAPIKeys:     map[string]struct{}{"test-api-key": {}},
		RateLimit: struct {
			Rate     float64
			Capacity float64
		}{
			Rate:     10.0,
			Capacity: 30.0,
		},
	}
}

// mockStatusProvider is a temporary mock implementation of status.StatusProvider
type mockStatusProvider struct{}

func (m *mockStatusProvider) GetNewCoinWatcher() status.WatcherStatus {
	return &mockWatcher{}
}

func (m *mockStatusProvider) GetPositionMonitor() status.WatcherStatus {
	return &mockWatcher{}
}

// mockWatcher is a temporary mock implementation of status.WatcherStatus
type mockWatcher struct{}

func (m *mockWatcher) IsRunning() bool {
	return false
}

func (m *mockWatcher) Start(ctx context.Context) error {
	return nil
}

func (m *mockWatcher) Stop() {
	// Do nothing
}

// mockExchangeService is a temporary mock implementation of service.ExchangeService
type mockExchangeService struct{}

func (m *mockExchangeService) GetAllTickers(ctx context.Context) (map[string]*models.Ticker, error) {
	return map[string]*models.Ticker{
		"BTCUSDT": {
			Symbol:         "BTCUSDT",
			Price:          50000.0,
			Volume:         1000000.0,
			PriceChange:    1000.0,
			PriceChangePct: 2.0,
			QuoteVolume:    50000000.0,
			High24h:        51000.0,
			Low24h:         49000.0,
			Timestamp:      time.Now(),
		},
		"ETHUSDT": {
			Symbol:         "ETHUSDT",
			Price:          3000.0,
			Volume:         500000.0,
			PriceChange:    100.0,
			PriceChangePct: 3.5,
			QuoteVolume:    1500000.0,
			High24h:        3100.0,
			Low24h:         2900.0,
			Timestamp:      time.Now(),
		},
	}, nil
}

func (m *mockExchangeService) GetTicker(ctx context.Context, symbol string) (*models.Ticker, error) {
	return &models.Ticker{
		Symbol:         symbol,
		Price:          50000.0,
		Volume:         1000000.0,
		PriceChange:    1000.0,
		PriceChangePct: 2.0,
		QuoteVolume:    50000000.0,
		High24h:        51000.0,
		Low24h:         49000.0,
		Timestamp:      time.Now(),
	}, nil
}

func (m *mockExchangeService) Connect(ctx context.Context) error {
	return nil
}

func (m *mockExchangeService) Disconnect() error {
	return nil
}

func (m *mockExchangeService) GetKlines(ctx context.Context, symbol, interval string, limit int) ([]*models.Kline, error) {
	return []*models.Kline{}, nil
}

func (m *mockExchangeService) GetWallet(ctx context.Context) (*models.Wallet, error) {
	return &models.Wallet{}, nil
}

func (m *mockExchangeService) PlaceOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	return order, nil
}

func (m *mockExchangeService) CancelOrder(ctx context.Context, orderID, symbol string) error {
	return nil
}

func (m *mockExchangeService) GetOrder(ctx context.Context, orderID, symbol string) (*models.Order, error) {
	return &models.Order{}, nil
}

func (m *mockExchangeService) GetOpenOrders(ctx context.Context, symbol string) ([]*models.Order, error) {
	return []*models.Order{}, nil
}

func (m *mockExchangeService) SubscribeToTickers(ctx context.Context, symbols []string, updates chan<- *models.Ticker) error {
	return nil
}

func (m *mockExchangeService) UnsubscribeFromTickers(ctx context.Context, symbols []string) error {
	return nil
}

func (m *mockExchangeService) GetNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	return []*models.NewCoin{}, nil
}

// mockTradeService is a temporary mock implementation of trade.TradeService
type mockTradeService struct{}

func (m *mockTradeService) EvaluatePurchaseDecision(ctx context.Context, symbol string) (*models.PurchaseDecision, error) {
	return &models.PurchaseDecision{}, nil
}

func (m *mockTradeService) ExecutePurchase(ctx context.Context, symbol string, amount float64, options *models.PurchaseOptions) (*models.BoughtCoin, error) {
	return &models.BoughtCoin{}, nil
}

func (m *mockTradeService) CheckStopLoss(ctx context.Context, coin *models.BoughtCoin) (bool, error) {
	return false, nil
}

func (m *mockTradeService) CheckTakeProfit(ctx context.Context, coin *models.BoughtCoin) (bool, error) {
	return false, nil
}

func (m *mockTradeService) SellCoin(ctx context.Context, coin *models.BoughtCoin, amount float64) (*models.Order, error) {
	return &models.Order{}, nil
}

func (m *mockTradeService) CancelOrder(ctx context.Context, orderID string) error {
	return nil
}

func (m *mockTradeService) GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error) {
	return []*models.BoughtCoin{}, nil
}

func (m *mockTradeService) ExecuteTrade(ctx context.Context, order *models.Order) (*models.Order, error) {
	return order, nil
}

func (m *mockTradeService) GetTradeHistory(ctx context.Context, startTime time.Time, limit int) ([]*models.Order, error) {
	return []*models.Order{}, nil
}

func (m *mockTradeService) GetOrderStatus(ctx context.Context, orderID string) (*models.Order, error) {
	return &models.Order{}, nil
}

func (m *mockTradeService) GetPendingOrders(ctx context.Context) ([]*models.Order, error) {
	return []*models.Order{}, nil
}

// mockPortfolioService is a temporary mock implementation of handlers.PortfolioServiceInterface
type mockPortfolioService struct{}

func (m *mockPortfolioService) GetPortfolioValue(ctx context.Context) (float64, error) {
	return 1000.0, nil
}

func (m *mockPortfolioService) GetActiveTrades(ctx context.Context) ([]*models.BoughtCoin, error) {
	return []*models.BoughtCoin{}, nil
}

func (m *mockPortfolioService) GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error) {
	return &models.PerformanceMetrics{
		TotalTrades:           10,
		WinningTrades:         7,
		LosingTrades:          3,
		WinRate:               0.7,
		TotalProfitLoss:       1000.0,
		AverageProfitPerTrade: 100.0,
		LargestProfit:         500.0,
		LargestLoss:           -200.0,
	}, nil
}

// mockBoughtCoinRepository is a temporary mock implementation of repository.BoughtCoinRepository
type mockBoughtCoinRepository struct{}

func (m *mockBoughtCoinRepository) FindAll(ctx context.Context) ([]*models.BoughtCoin, error) {
	return []*models.BoughtCoin{}, nil
}

func (m *mockBoughtCoinRepository) FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error) {
	return &models.BoughtCoin{}, nil
}

func (m *mockBoughtCoinRepository) FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error) {
	return &models.BoughtCoin{}, nil
}

func (m *mockBoughtCoinRepository) Save(ctx context.Context, coin *models.BoughtCoin) error {
	return nil
}

func (m *mockBoughtCoinRepository) Delete(ctx context.Context, symbol string) error {
	return nil
}

func (m *mockBoughtCoinRepository) DeleteByID(ctx context.Context, id int64) error {
	return nil
}

func (m *mockBoughtCoinRepository) UpdatePrice(ctx context.Context, symbol string, price float64) error {
	return nil
}

func (m *mockBoughtCoinRepository) FindAllActive(ctx context.Context) ([]*models.BoughtCoin, error) {
	return []*models.BoughtCoin{}, nil
}

func (m *mockBoughtCoinRepository) HardDelete(ctx context.Context, symbol string) error {
	return nil
}

func (m *mockBoughtCoinRepository) Count(ctx context.Context) (int64, error) {
	return 0, nil
}

// mockLogger is a temporary mock implementation of middleware.Logger
type mockLogger struct{}

func (m *mockLogger) Info(args ...interface{}) {
	log.Println(args...)
}

func (m *mockLogger) Error(args ...interface{}) {
	log.Println(args...)
}

// mockNewCoinService is a temporary mock implementation of newcoin.NewCoinService
type mockNewCoinService struct{}

func (m *mockNewCoinService) DetectNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	return []*models.NewCoin{}, nil
}

func (m *mockNewCoinService) SaveNewCoins(ctx context.Context, coins []*models.NewCoin) error {
	return nil
}

func (m *mockNewCoinService) GetAllNewCoins(ctx context.Context) ([]models.NewCoin, error) {
	return []models.NewCoin{}, nil
}

func (m *mockNewCoinService) StartWatching(ctx context.Context, interval time.Duration) error {
	return nil
}

func (m *mockNewCoinService) StopWatching() {}

func (m *mockNewCoinService) MarkAsProcessed(ctx context.Context, id int64) error {
	return nil
}

func (m *mockNewCoinService) GetCoinByID(ctx context.Context, id int64) (*models.NewCoin, error) {
	return &models.NewCoin{}, nil
}

func (m *mockNewCoinService) GetCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	return []models.NewCoin{}, nil
}

func (m *mockNewCoinService) GetCoinsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error) {
	return []models.NewCoin{}, nil
}

func (m *mockNewCoinService) GetUpcomingCoins(ctx context.Context) ([]models.NewCoin, error) {
	// Return mock upcoming coins
	return []models.NewCoin{
		{
			ID:            1,
			Symbol:        "MOCKCOIN1USDT",
			FoundAt:       time.Now(),
			FirstOpenTime: toTimePtr(time.Now().Add(24 * time.Hour)),
			QuoteVolume:   5000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
		{
			ID:            2,
			Symbol:        "MOCKCOIN2USDT",
			FoundAt:       time.Now(),
			FirstOpenTime: toTimePtr(time.Now().Add(48 * time.Hour)),
			QuoteVolume:   3000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
	}, nil
}

func (m *mockNewCoinService) GetUpcomingCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	// Return mock upcoming coins for the specified date
	return []models.NewCoin{
		{
			ID:            1,
			Symbol:        "MOCKCOIN1USDT",
			FoundAt:       time.Now(),
			FirstOpenTime: toTimePtr(date),
			QuoteVolume:   5000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
	}, nil
}

func (m *mockNewCoinService) GetUpcomingCoinsForTodayAndTomorrow(ctx context.Context) ([]models.NewCoin, error) {
	// Get current time
	now := time.Now()

	// Return mock upcoming coins for today and tomorrow
	return []models.NewCoin{
		{
			ID:            1,
			Symbol:        "TODAYCOIN1USDT",
			FoundAt:       now,
			FirstOpenTime: toTimePtr(now.Add(6 * time.Hour)), // Later today
			QuoteVolume:   5000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
		{
			ID:            2,
			Symbol:        "TODAYCOIN2USDT",
			FoundAt:       now,
			FirstOpenTime: toTimePtr(now.Add(12 * time.Hour)), // Later today
			QuoteVolume:   3000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
		{
			ID:            3,
			Symbol:        "TOMORROWCOIN1USDT",
			FoundAt:       now,
			FirstOpenTime: toTimePtr(now.Add(24 * time.Hour)), // Tomorrow
			QuoteVolume:   7000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
		{
			ID:            4,
			Symbol:        "TOMORROWCOIN2USDT",
			FoundAt:       now,
			FirstOpenTime: toTimePtr(now.Add(36 * time.Hour)), // Tomorrow
			QuoteVolume:   9000.0,
			IsProcessed:   false,
			IsUpcoming:    true,
		},
	}, nil
}

func (m *mockNewCoinService) UpdateCoinStatus(ctx context.Context, symbol string, status string) error {
	return nil
}

func (m *mockNewCoinService) GetTradableCoins(ctx context.Context) ([]models.NewCoin, error) {
	// Get current time
	now := time.Now()

	// Return mock tradable coins
	return []models.NewCoin{
		{
			ID:               1,
			Symbol:           "TRADABLECOIN1USDT",
			FoundAt:          now.Add(-24 * time.Hour),
			FirstOpenTime:    toTimePtr(now.Add(-12 * time.Hour)),
			QuoteVolume:      5000.0,
			Status:           "1",
			BecameTradableAt: toTimePtr(now.Add(-6 * time.Hour)),
			IsProcessed:      false,
			IsUpcoming:       false,
		},
		{
			ID:               2,
			Symbol:           "TRADABLECOIN2USDT",
			FoundAt:          now.Add(-48 * time.Hour),
			FirstOpenTime:    toTimePtr(now.Add(-36 * time.Hour)),
			QuoteVolume:      3000.0,
			Status:           "1",
			BecameTradableAt: toTimePtr(now.Add(-24 * time.Hour)),
			IsProcessed:      false,
			IsUpcoming:       false,
		},
	}, nil
}

func (m *mockNewCoinService) GetTradableCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	// Get current time
	now := time.Now()

	// Check if the date is today
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	specifiedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	if today.Equal(specifiedDate) {
		// Return coins that became tradable today
		return []models.NewCoin{
			{
				ID:               3,
				Symbol:           "TODAYTRADABLECOIN1USDT",
				FoundAt:          now.Add(-12 * time.Hour),
				FirstOpenTime:    toTimePtr(now.Add(-6 * time.Hour)),
				QuoteVolume:      7000.0,
				Status:           "1",
				BecameTradableAt: toTimePtr(now.Add(-3 * time.Hour)),
				IsProcessed:      false,
				IsUpcoming:       false,
			},
			{
				ID:               4,
				Symbol:           "TODAYTRADABLECOIN2USDT",
				FoundAt:          now.Add(-10 * time.Hour),
				FirstOpenTime:    toTimePtr(now.Add(-5 * time.Hour)),
				QuoteVolume:      9000.0,
				Status:           "1",
				BecameTradableAt: toTimePtr(now.Add(-1 * time.Hour)),
				IsProcessed:      false,
				IsUpcoming:       false,
			},
		}, nil
	}

	// Return empty list for other dates
	return []models.NewCoin{}, nil
}

func (m *mockNewCoinService) GetTradableCoinsToday(ctx context.Context) ([]models.NewCoin, error) {
	// Get current time
	now := time.Now()

	// Return coins that became tradable today
	return []models.NewCoin{
		{
			ID:               3,
			Symbol:           "TODAYTRADABLECOIN1USDT",
			FoundAt:          now.Add(-12 * time.Hour),
			FirstOpenTime:    toTimePtr(now.Add(-6 * time.Hour)),
			QuoteVolume:      7000.0,
			Status:           "1",
			BecameTradableAt: toTimePtr(now.Add(-3 * time.Hour)),
			IsProcessed:      false,
			IsUpcoming:       false,
		},
		{
			ID:               4,
			Symbol:           "TODAYTRADABLECOIN2USDT",
			FoundAt:          now.Add(-10 * time.Hour),
			FirstOpenTime:    toTimePtr(now.Add(-5 * time.Hour)),
			QuoteVolume:      9000.0,
			Status:           "1",
			BecameTradableAt: toTimePtr(now.Add(-1 * time.Hour)),
			IsProcessed:      false,
			IsUpcoming:       false,
		},
	}, nil
}

// Server wraps the HTTP server and router.
type Server struct {
	httpServer *http.Server
	router     http.Handler
}

// NewServer creates a new API server instance.
func NewServer(deps *Dependencies, addr string) *Server {
	// Use consolidated Chi router
	router := SetupConsolidatedRouter(deps)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &Server{
		httpServer: srv,
		router:     router,
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	log.Printf("Starting API server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping API server...")
	return s.httpServer.Shutdown(ctx)
}

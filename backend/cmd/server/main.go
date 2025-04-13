package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	httpAdapter "github.com/neo/crypto-bot/internal/adapter/http"
	"github.com/neo/crypto-bot/internal/adapter/http/handler"
	gormAdapter "github.com/neo/crypto-bot/internal/adapter/persistence/gorm"
	"github.com/neo/crypto-bot/internal/config"
	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/domain/port"
	"github.com/neo/crypto-bot/internal/factory"
	"github.com/neo/crypto-bot/internal/platform/logger"
)

// mockTradeUseCase is a temporary mock implementation of the TradeUseCase interface
type mockTradeUseCase struct{}

func (m *mockTradeUseCase) PlaceOrder(ctx context.Context, req model.OrderRequest) (*model.Order, error) {
	// Return a dummy order for now
	return &model.Order{
		ID:     "mock-order-id",
		Symbol: req.Symbol,
		Side:   req.Side,
		Type:   req.Type,
	}, nil
}

// Implement additional methods to satisfy the TradeUseCase interface
func (m *mockTradeUseCase) CancelOrder(ctx context.Context, symbol, orderID string) error {
	return nil
}

func (m *mockTradeUseCase) GetOrderStatus(ctx context.Context, symbol, orderID string) (*model.Order, error) {
	return &model.Order{
		ID:        orderID,
		Symbol:    symbol,
		Status:    model.OrderStatusNew,
		CreatedAt: time.Now(),
	}, nil
}

func (m *mockTradeUseCase) GetOpenOrders(ctx context.Context, symbol string) ([]*model.Order, error) {
	return []*model.Order{}, nil
}

func (m *mockTradeUseCase) GetOrderHistory(ctx context.Context, symbol string, limit, offset int) ([]*model.Order, error) {
	return []*model.Order{}, nil
}

func (m *mockTradeUseCase) CalculateRequiredQuantity(ctx context.Context, symbol string, side model.OrderSide, amount float64) (float64, error) {
	return 0.001, nil
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	loggerInstance := logger.New(cfg.LogLevel)
	l := &loggerInstance
	l.Info().Msg("Starting crypto trading bot server...")

	// Setup Database Connection
	dbConn, err := gormAdapter.NewDBConnection(cfg, *l)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to connect to database")
	}
	// Run Migrations
	gormAdapter.AutoMigrateModels(dbConn, l)

	// Set up router
	router, apiV1 := httpAdapter.SetupRouter(*l)

	// Create factories
	aiFactory := factory.NewAIFactory(cfg, *l)
	marketFactory := factory.NewMarketFactory(cfg, l, dbConn)
	positionFactory := factory.NewPositionFactory(cfg, l, dbConn)
	tradeFactory := factory.NewTradeFactory(cfg, l, dbConn)

	// Create AI handler
	aiHandler, err := aiFactory.CreateAIHandler()
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to create AI handler")
	}

	// Create Market Data Use Case and Handler
	marketUseCase, err := marketFactory.CreateMarketDataUseCase()
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to create Market Data Use Case")
	}
	marketHandler := handler.NewMarketDataHandler(marketUseCase, l)

	// Create Market Data Service
	marketRepo, symbolRepo := marketFactory.CreateMarketRepository()
	cache := marketFactory.CreateMarketCache()

	// TODO: Replace with actual MEXC API implementation when available
	var mexcAPI port.MexcAPI = nil

	marketService := marketFactory.CreateMarketDataService(
		marketRepo,
		symbolRepo,
		cache,
		mexcAPI,
	)

	// Register AI routes
	// For now, we'll use a dummy auth middleware
	aiHandler.RegisterRoutes(router, func(c *gin.Context) {
		// In a real implementation, this would validate JWT tokens
		c.Next()
	})

	// Register Market Data routes
	marketHandler.RegisterRoutes(apiV1)

	// Create and register WebSocket handler
	wsHandler := handler.NewWebSocketHandler(marketUseCase, l)
	wsHandler.RegisterRoutes(apiV1)
	wsHandler.Start()

	// Create Position Use Case and Monitor
	positionUC, err := positionFactory.CreatePositionUseCase(marketRepo, symbolRepo)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to create Position Use Case")
	}

	// Create Position Handler
	positionHandler := handler.NewPositionHandler(positionUC, l)
	positionHandler.RegisterRoutes(apiV1)

	// Create Order Repository
	orderRepo := tradeFactory.CreateOrderRepository()

	// Create Trade Service
	tradeService := tradeFactory.CreateTradeService(
		mexcAPI,
		marketService,
		symbolRepo,
		orderRepo,
	)

	// Create Trade Use Case
	tradeUC := tradeFactory.CreateTradeUseCase(
		mexcAPI,
		symbolRepo,
		orderRepo,
		tradeService,
	)

	// Create Trade Handler
	tradeHandler := tradeFactory.CreateTradeHandler(tradeUC)
	tradeHandler.RegisterRoutes(apiV1)

	// For Position Monitor, use the TradeUseCase we just created
	positionMonitor := positionFactory.CreatePositionMonitor(
		positionUC,
		marketService,
		tradeUC,
	)

	// Start Position Monitor
	positionMonitor.Start()

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		l.Info().Msgf("Server is listening on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info().Msg("Shutting down server...")

	// Stop services
	wsHandler.Stop()
	positionMonitor.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	l.Info().Msg("Server exited properly")
}

package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/event"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// Mocks and dummy implementations

type mockConfigLoader struct {
	config *AutoBuyConfig
}

func (m *mockConfigLoader) LoadAutoBuyConfig() (*AutoBuyConfig, error) {
	return m.config, nil
}

type mockNewCoinRepository struct {
	processed bool
}

func (m *mockNewCoinRepository) IsProcessedForAutobuy(symbol string) bool {
	return m.processed
}

func (m *mockNewCoinRepository) MarkAsProcessed(symbol string) error {
	m.processed = true
	return nil
}

type mockMarketDataService struct {
	price  float64
	volume float64
}

func (m *mockMarketDataService) GetMarketData(symbol string) (float64, float64, error) {
	return m.price, m.volume, nil
}

type mockRiskUsecase struct {
	riskOk bool
}

func (m *mockRiskUsecase) CheckRisk(order OrderParameters) error {
	if m.riskOk {
		return nil
	}
	return errors.New("risk check failed")
}

type mockTradeUsecase struct {
	executed bool
}

func (m *mockTradeUsecase) ExecuteMarketBuy(order OrderParameters) error {
	m.executed = true
	return nil
}

type mockNotificationService struct {
	notified bool
}

func (m *mockNotificationService) Notify(message string) {
	m.notified = true
}

// Helper to create a test NewCoin model
func createTestNewCoin(symbol string, quoteAsset string) *model.NewCoin {
	now := time.Now()
	return &model.NewCoin{
		Symbol:           symbol,
		QuoteAsset:       quoteAsset,
		BecameTradableAt: &now,
	}
}

func TestAutobuyService_DisabledConfig(t *testing.T) {
	config := &AutoBuyConfig{
		Enabled: false,
	}
	service := &AutobuyService{
		configLoader:        &mockConfigLoader{config: config},
		newCoinRepository:   &mockNewCoinRepository{},
		marketDataService:   &mockMarketDataService{},
		riskUsecase:         &mockRiskUsecase{riskOk: true},
		tradeUsecase:        &mockTradeUsecase{},
		notificationService: &mockNotificationService{},
	}

	price := 100.0
	volume := 1000.0
	coin := createTestNewCoin("COIN1", "USDT")
	evt := event.NewNewCoinTradable(coin, &price, &volume)

	err := service.HandleNewCoinEvent(*evt)
	if err == nil {
		t.Error("Expected error due to disabled autobuy config, got nil")
	}
}

func TestAutobuyService_DuplicatePrevention(t *testing.T) {
	config := &AutoBuyConfig{
		Enabled:    true,
		QuoteAsset: "USDT",
		MinPrice:   10,
		MaxPrice:   200,
		MinVolume:  500,
	}
	repo := &mockNewCoinRepository{processed: true}
	service := &AutobuyService{
		configLoader:        &mockConfigLoader{config: config},
		newCoinRepository:   repo,
		marketDataService:   &mockMarketDataService{},
		riskUsecase:         &mockRiskUsecase{riskOk: true},
		tradeUsecase:        &mockTradeUsecase{},
		notificationService: &mockNotificationService{},
	}

	price := 50.0
	volume := 600.0
	coin := createTestNewCoin("COIN2", "USDT")
	evt := event.NewNewCoinTradable(coin, &price, &volume)

	err := service.HandleNewCoinEvent(*evt)
	if err == nil {
		t.Error("Expected error due to duplicate processing, got nil")
	}
}

func TestAutobuyService_RiskCheckFailure(t *testing.T) {
	config := &AutoBuyConfig{
		Enabled:    true,
		QuoteAsset: "USDT",
		MinPrice:   10,
		MaxPrice:   200,
		MinVolume:  500,
	}
	repo := &mockNewCoinRepository{}
	marketData := &mockMarketDataService{price: 50, volume: 600}
	risk := &mockRiskUsecase{riskOk: false}
	trade := &mockTradeUsecase{}
	notify := &mockNotificationService{}
	service := &AutobuyService{
		configLoader:        &mockConfigLoader{config: config},
		newCoinRepository:   repo,
		marketDataService:   marketData,
		riskUsecase:         risk,
		tradeUsecase:        trade,
		notificationService: notify,
	}

	price := 50.0
	volume := 600.0
	coin := createTestNewCoin("COIN3", "USDT")
	evt := event.NewNewCoinTradable(coin, &price, &volume)

	err := service.HandleNewCoinEvent(*evt)
	if err == nil || err.Error() != "risk check failed" {
		t.Errorf("Expected 'risk check failed' error, got: %v", err)
	}
	if repo.processed {
		t.Error("Repository should not mark as processed when risk check fails")
	}
}

func TestAutobuyService_Success(t *testing.T) {
	config := &AutoBuyConfig{
		Enabled:      true,
		QuoteAsset:   "USDT",
		MinPrice:     10,
		MaxPrice:     200,
		MinVolume:    500,
		DelaySeconds: 0,
	}
	repo := &mockNewCoinRepository{}
	marketData := &mockMarketDataService{price: 50, volume: 600}
	risk := &mockRiskUsecase{riskOk: true}
	trade := &mockTradeUsecase{}
	notify := &mockNotificationService{}
	service := &AutobuyService{
		configLoader:        &mockConfigLoader{config: config},
		newCoinRepository:   repo,
		marketDataService:   marketData,
		riskUsecase:         risk,
		tradeUsecase:        trade,
		notificationService: notify,
	}

	price := 50.0
	volume := 600.0
	coin := createTestNewCoin("COIN4", "USDT")
	evt := event.NewNewCoinTradable(coin, &price, &volume)

	err := service.HandleNewCoinEvent(*evt)
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	if !repo.processed {
		t.Error("Repository should be marked as processed")
	}
	if !trade.executed {
		t.Error("Trade should be executed")
	}
	if !notify.notified {
		t.Error("Notification should be sent")
	}
	// Allow any optional delay to pass
	time.Sleep(10 * time.Millisecond)
}

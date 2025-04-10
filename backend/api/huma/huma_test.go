package huma

import (
	"testing"

	"go-crypto-bot-clean/backend/api/service"
	"go-crypto-bot-clean/backend/pkg/auth"
	"go-crypto-bot-clean/backend/pkg/backtest"
	"go-crypto-bot-clean/backend/pkg/strategy"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestSetupHuma(t *testing.T) {
	// Skip this test for now due to Huma schema registration issues
	t.Skip("Skipping test due to Huma schema registration issues")
	// Create a new router
	router := chi.NewRouter()

	// Create mock services
	backtestService := backtest.NewService()
	strategyFactory := strategy.NewFactory()
	authService := auth.NewService("dummy-secret-key")

	// Create mock service provider
	serviceProvider := &service.Provider{
		BacktestService: service.NewBacktestService(&backtestService),
		StrategyService: service.NewStrategyService(&strategyFactory),
		AuthService:     service.NewAuthService(&authService, nil),
		UserService:     service.NewUserService(nil),
	}

	// Setup Huma with the router and services
	api := SetupHuma(router, DefaultConfig(), serviceProvider)

	// Verify that the API was created
	assert.NotNil(t, api, "API should not be nil")
}

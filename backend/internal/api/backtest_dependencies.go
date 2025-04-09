package api

import (
	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/backtest"
)

// InitializeBacktestDependencies initializes the backtest dependencies
func (deps *Dependencies) InitializeBacktestDependencies() {
	// Create backtest service
	backtestService := backtest.NewService()
	
	// Create backtest handler
	deps.BacktestHandler = handlers.NewBacktestHandler(backtestService)
}

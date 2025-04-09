package bot

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/core/newcoin"
	"go-crypto-bot-clean/backend/internal/core/portfolio"
	"go-crypto-bot-clean/backend/internal/core/trade"
)

// Config holds configuration for the trading bot
type Config struct {
	NewCoinCheckInterval   time.Duration
	PortfolioCheckInterval time.Duration
	MaxConcurrentRequests  int
	Headless               bool
}

// Bot represents the main trading bot
type Bot struct {
	tradeService     trade.TradeService
	newCoinService   newcoin.NewCoinService
	portfolioService portfolio.PortfolioService
	config           Config
	stopChan         chan struct{}
}

// NewBot creates a new instance of the trading bot
func NewBot(
	tradeService trade.TradeService,
	newCoinService newcoin.NewCoinService,
	portfolioService portfolio.PortfolioService,
	config Config,
) *Bot {
	return &Bot{
		tradeService:     tradeService,
		newCoinService:   newCoinService,
		portfolioService: portfolioService,
		config:           config,
		stopChan:         make(chan struct{}),
	}
}

// Run starts the bot and begins monitoring for trading opportunities
func (b *Bot) Run(ctx context.Context) error {
	fmt.Println("Bot started with configuration:")
	fmt.Printf("- New coin check interval: %v\n", b.config.NewCoinCheckInterval)
	fmt.Printf("- Portfolio check interval: %v\n", b.config.PortfolioCheckInterval)
	fmt.Printf("- Headless mode: %v\n", b.config.Headless)

	// Start new coin detection in a goroutine
	go b.monitorNewCoins(ctx)

	// Start portfolio monitoring in a goroutine
	go b.monitorPortfolio(ctx)

	// Wait for context cancellation
	<-ctx.Done()
	return ctx.Err()
}

// Stop stops the bot
func (b *Bot) Stop() {
	close(b.stopChan)
}

// monitorNewCoins periodically checks for new coins
func (b *Bot) monitorNewCoins(ctx context.Context) {
	ticker := time.NewTicker(b.config.NewCoinCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-b.stopChan:
			return
		case <-ticker.C:
			b.checkNewCoins(ctx)
		}
	}
}

// monitorPortfolio periodically checks portfolio status
func (b *Bot) monitorPortfolio(ctx context.Context) {
	ticker := time.NewTicker(b.config.PortfolioCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-b.stopChan:
			return
		case <-ticker.C:
			b.checkPortfolio(ctx)
		}
	}
}

// checkNewCoins detects and processes new coins
func (b *Bot) checkNewCoins(ctx context.Context) {
	newCoins, err := b.newCoinService.DetectNewCoins(ctx)
	if err != nil {
		fmt.Printf("Error detecting new coins: %v\n", err)
		return
	}

	if len(newCoins) > 0 {
		fmt.Printf("Detected %d new coins\n", len(newCoins))
		
		// Save the new coins
		err = b.newCoinService.SaveNewCoins(ctx, newCoins)
		if err != nil {
			fmt.Printf("Error saving new coins: %v\n", err)
		}
	}
}

// checkPortfolio checks portfolio status
func (b *Bot) checkPortfolio(ctx context.Context) {
	positions, err := b.portfolioService.GetPositions(ctx)
	if err != nil {
		fmt.Printf("Error checking portfolio: %v\n", err)
		return
	}

	if len(positions) > 0 {
		fmt.Printf("Current positions: %d\n", len(positions))
		// Process positions as needed
	}
}

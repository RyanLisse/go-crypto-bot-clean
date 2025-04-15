package standard

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model/market"
)

func ExampleStandardCache_GetTickerWithError() {
	// Create a new cache with a short TTL for demonstration
	cache := NewStandardCache(5*time.Second, 10*time.Second)

	// Create a sample ticker
	ticker := &market.Ticker{
		Symbol:        "BTC-USDT",
		Exchange:      "mexc",
		Price:         50000.0,
		PriceChange:   1.5,
		Volume:        1000,
		High24h:       51000.0,
		Low24h:        49000.0,
		PercentChange: 3.0,
		LastUpdated:   time.Now(),
	}

	// Cache the ticker
	cache.CacheTicker(ticker)

	// Retrieve the ticker with error handling
	ctx := context.Background()
	retrievedTicker, err := cache.GetTickerWithError(ctx, "mexc", "BTC-USDT")
	if err != nil {
		var cacheErr *CacheError
		if apperror.As(err, &cacheErr) {
			switch cacheErr.Code {
			case ErrCacheKeyNotFound:
				// Handle "not found" case
				log.Printf("Ticker not found in cache: %v", err)
			case ErrCacheExpired:
				// Handle "expired" case
				log.Printf("Ticker expired in cache: %v", err)
			case ErrCacheInvalidType:
				// Handle "invalid type" case
				log.Printf("Invalid ticker type in cache: %v", err)
			default:
				log.Printf("Cache error: %v", err)
			}
		} else {
			// Handle general error
			log.Printf("Error retrieving ticker: %v", err)
		}
		return
	}

	fmt.Printf("Retrieved ticker for %s on %s: price %.2f",
		retrievedTicker.Symbol, retrievedTicker.Exchange, retrievedTicker.Price)
}

func Example_errorHandlingWithCache() {
	// Create a new cache
	cache := NewStandardCache(1*time.Minute, 5*time.Minute)

	// Example function showing how to handle cache errors in an application
	fetchTicker := func(ctx context.Context, exchange, symbol string) (*market.Ticker, error) {
		// Try to get from cache first with error handling
		ticker, err := cache.GetTickerWithError(ctx, exchange, symbol)
		if err == nil {
			return ticker, nil
		}

		// Convert cache error to application error
		appErr := ConvertCacheError(err)
		if apperror.Is(appErr, apperror.ErrNotFound) {
			// Cache miss, fetch from API or database
			// For this example, we'll just create a new ticker
			log.Printf("Cache miss, fetching from API")

			// Simulate API call
			newTicker := &market.Ticker{
				Symbol:        symbol,
				Exchange:      exchange,
				Price:         49000.0,
				PriceChange:   -0.5,
				Volume:        2000,
				High24h:       50000.0,
				Low24h:        48000.0,
				PercentChange: -1.0,
				LastUpdated:   time.Now(),
			}

			// Cache the new ticker
			cache.CacheTicker(newTicker)
			return newTicker, nil
		}

		// For other errors, return the application error
		return nil, appErr
	}

	// Use the function
	ctx := context.Background()
	ticker, err := fetchTicker(ctx, "mexc", "ETH-USDT")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Successfully retrieved ticker for %s: %.2f", ticker.Symbol, ticker.Price)
}

func ExampleConvertCacheError() {
	// Example of converting a cache error to an application error
	ctx := context.Background()
	cache := NewStandardCache(1*time.Minute, 5*time.Minute)

	// Try to get a non-existent ticker
	_, err := cache.GetTickerWithError(ctx, "mexc", "NON-EXISTENT")
	if err != nil {
		// Convert to application error
		appErr := ConvertCacheError(err)

		// Now we can use the application error in our HTTP handlers
		switch {
		case apperror.Is(appErr, apperror.ErrNotFound):
			// Would return 404 in a real handler
			fmt.Println("Resource not found, would return 404")
		case apperror.Is(appErr, apperror.ErrInternal):
			// Would return 500 in a real handler
			fmt.Println("Internal server error, would return 500")
		case apperror.Is(appErr, apperror.ErrInvalidInput):
			// Would return 400 in a real handler
			fmt.Println("Invalid input, would return 400")
		default:
			// Would return 500 in a real handler
			fmt.Println("Unknown error, would return 500")
		}
	}
}

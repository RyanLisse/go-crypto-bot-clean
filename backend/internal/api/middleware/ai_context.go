package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/portfolio"
	"go-crypto-bot-clean/backend/internal/domain/trade"
)

// AIContextMiddleware adds trading context to AI requests
func AIContextMiddleware(
	portfolioSvc portfolio.Service,
	tradeSvc trade.Service,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract user ID from context (set by auth middleware)
			userID, ok := r.Context().Value("userID").(int)
			if !ok {
				http.Error(w, "User ID not found in context", http.StatusUnauthorized)
				return
			}

			// Gather portfolio context
			portfolioContext := map[string]interface{}{
				"timestamp": time.Now().UTC(),
			}

			// Try to get portfolio summary
			portfolio, err := portfolioSvc.GetSummary(r.Context(), userID)
			if err != nil {
				log.Printf("Error fetching portfolio context: %v", err)
				// Continue without portfolio context
			} else {
				portfolioContext["portfolio"] = portfolio
			}

			// Try to get active trades
			trades, err := tradeSvc.GetActiveTrades(r.Context(), userID)
			if err != nil {
				log.Printf("Error fetching trades context: %v", err)
				// Continue without trades context
			} else {
				portfolioContext["trades"] = trades
			}

			// Add context to request
			ctx := context.WithValue(r.Context(), "aiContext", portfolioContext)

			// Pass to next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

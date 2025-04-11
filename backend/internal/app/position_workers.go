package app

import (
	"context"
	"time"

	"go-crypto-bot-clean/backend/api/service"
)

// StartMarketDataListener simulates subscribing to market data feed
func StartMarketDataListener(ctx context.Context, svc *service.PositionService) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				// Simulate receiving a market tick
				tick := service.MarketTick{
					Symbol: "BTCUSDT",
					Price:  50000 + float64(t.Second()), // dummy price
				}
				_ = svc.HandleMarketTick(ctx, tick)
			}
		}
	}()
}

// StartOrderEventListener simulates subscribing to order execution events
func StartOrderEventListener(ctx context.Context, svc *service.PositionService) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				// Simulate receiving an order event
				event := service.OrderEvent{
					PositionID: "some-position-id",
					Status:     "closed",
					ClosePrice: 50500 + float64(t.Second()),
					CloseTime:  time.Now(),
				}
				_ = svc.HandleOrderEvent(ctx, event)
			}
		}
	}()
}

// Example: Real-time MEXC New Listing Watcher using Protobuf WebSocket
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc/websocket"
	mexcproto "github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc/websocket/proto"
)

func main() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := websocket.NewProtobufClient(ctx, &log.Logger)

	// Register handler for new listings
	client.RegisterNewListingHandler(func(msg *mexcproto.MexcMessage) error {
		newListingData := msg.GetNewListingData()
		if newListingData == nil {
			return nil
		}
		for _, listing := range newListingData.Listings {
			log.Info().
				Str("symbol", listing.Symbol).
				Str("base_asset", listing.BaseAsset).
				Str("quote_asset", listing.QuoteAsset).
				Str("status", listing.Status).
				Int64("listing_time", listing.ListingTime).
				Int64("trading_time", listing.TradingTime).
				Str("initial_price", listing.InitialPrice).
				Msg("Detected new MEXC listing")
		}
		return nil
	})

	// Connect to WebSocket
	if err := client.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MEXC WebSocket")
	}
	defer client.Disconnect()

	// Subscribe to new listings
	if err := client.SubscribeToNewListings(); err != nil {
		log.Fatal().Err(err).Msg("Failed to subscribe to new listings")
	}

	log.Info().Msg("MEXC new listing watcher started. Waiting for events...")

	// Wait for interrupt signal to gracefully shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Info().Msg("Shutting down new listing watcher")
}

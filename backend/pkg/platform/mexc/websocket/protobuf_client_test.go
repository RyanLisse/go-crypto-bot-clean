package websocket

import (
    "context"
    "testing"

    "github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc/websocket/proto"
    "github.com/rs/zerolog"
)

func TestHandlerRegistration(t *testing.T) {
    logger := zerolog.Nop()
    client := NewProtobufClient(context.Background(), &logger)

    // No handlers or subscriptions initially
    if len(client.handlers) != 0 {
        t.Errorf("expected no handlers initially, got %d", len(client.handlers))
    }
    if len(client.subscriptions) != 0 {
        t.Errorf("expected no subscriptions initially, got %d", len(client.subscriptions))
    }

    // Register handlers
    dummy := func(msg *proto.MexcMessage) error { return nil }
    client.RegisterNewListingHandler(dummy)
    client.RegisterSymbolStatusHandler(dummy)

    if handlers, ok := client.handlers[channelNewListings]; !ok || len(handlers) != 1 {
        t.Errorf("unexpected new listing handlers: %v", handlers)
    }
    if handlers, ok := client.handlers[channelSymbolStatus]; !ok || len(handlers) != 1 {
        t.Errorf("unexpected symbol status handlers: %v", handlers)
    }
}

package middleware

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestSimple(t *testing.T) {
	// Create a logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create a disabled middleware
	middleware := NewDisabledMiddleware(&logger)

	// Check that the middleware is not nil
	if middleware == nil {
		t.Error("Middleware should not be nil")
	}
}

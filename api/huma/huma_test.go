package huma

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestSetupHuma(t *testing.T) {
	// Create a new router
	router := chi.NewRouter()

	// Setup Huma with the router
	api := SetupHuma(router, DefaultConfig())

	// Verify that the API was created
	assert.NotNil(t, api, "API should not be nil")
}

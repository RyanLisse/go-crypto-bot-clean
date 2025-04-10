package huma

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestMinimal(t *testing.T) {
	// Create a new router
	router := chi.NewRouter()

	// Create a new Huma API
	api := humachi.New(router, huma.DefaultConfig("Test API", "1.0.0"))

	// Register a simple endpoint
	huma.Register(api, huma.Operation{
		OperationID: "test",
		Method:      http.MethodGet,
		Path:        "/test",
		Summary:     "Test endpoint",
		Description: "A simple test endpoint",
		Tags:        []string{"Test"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Message string `json:"message"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Message string `json:"message"`
			}
		}{}
		resp.Body.Message = "Hello, world!"
		return resp, nil
	})

	// Verify that the API was created
	assert.NotNil(t, api, "API should not be nil")
}

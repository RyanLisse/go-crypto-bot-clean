package huma

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	// "encoding/json" // No longer needed directly
	// "net/http" // No longer needed directly
)

// AccountDetailsInput defines the input for the account details endpoint (empty for GET).
type AccountDetailsInput struct{}

// AccountDetailsResponse mirrors the structure expected by the frontend's WalletResponse,
// wrapped in a Body field for Huma.
type AccountDetailsResponse struct {
	Body struct {
		Balances  map[string]AssetBalance `json:"balances"`
		UpdatedAt string                  `json:"updatedAt"`
	}
}

// AssetBalance represents the balance details for a single asset.
type AssetBalance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
	Total  float64 `json:"total"`
	Price  float64 `json:"price,omitempty"` // Optional price field
}

// AccountDetailsHandler handles GET requests to /api/v1/account/details using Huma signature.
// Assumes authentication is handled via middleware.
func AccountDetailsHandler(ctx context.Context, input *AccountDetailsInput) (*AccountDetailsResponse, error) {
	// In a real implementation, fetch actual account details based on authenticated user from ctx.
	// Returning mock data.
	resp := &AccountDetailsResponse{}
	resp.Body.Balances = map[string]AssetBalance{
		"BTC":  {Asset: "BTC", Free: 0.1, Locked: 0.01, Total: 0.11, Price: 60000.0},
		"ETH":  {Asset: "ETH", Free: 1.5, Locked: 0.1, Total: 1.6, Price: 3500.0},
		"USDT": {Asset: "USDT", Free: 5000.0, Locked: 100.0, Total: 5100.0, Price: 1.0},
	}
	resp.Body.UpdatedAt = time.Now().Format(time.RFC3339)

	return resp, nil // Return mock response and nil error
}

// --- Additional Account Endpoints ---

type ValidateKeysResponse struct {
	Body struct {
		Valid   bool   `json:"valid"`
		Message string `json:"message,omitempty"`
	} `json:"body"`
}

type ValidateKeysInput struct{}

// ValidateKeysHandler handles POST /api/v1/account/validate-keys
func ValidateKeysHandler(ctx context.Context, input *ValidateKeysInput) (*ValidateKeysResponse, error) {
	resp := &ValidateKeysResponse{}
	resp.Body.Valid = true
	resp.Body.Message = "API keys are valid (mock)"
	return resp, nil
}

// RegisterAccountEndpoints registers all account-related endpoints with Huma
func RegisterAccountEndpoints(api huma.API, basePath string) {
	// GET /account/details
	huma.Register(api, huma.Operation{
		OperationID: "get-account-details",
		Method:      http.MethodGet,
		Path:        basePath + "/account/details",
		Summary:     "Get account details",
		Description: "Returns the account details including wallet balances",
		Tags:        []string{"Account"},
	}, AccountDetailsHandler)

	// GET /account/wallet
	huma.Register(api, huma.Operation{
		OperationID: "get-wallet",
		Method:      http.MethodGet,
		Path:        basePath + "/account/wallet",
		Summary:     "Get wallet details",
		Description: "Returns the wallet details including balances",
		Tags:        []string{"Account"},
	}, AccountDetailsHandler) // Reuse the same handler for now

	// POST /account/validate-keys
	huma.Register(api, huma.Operation{
		OperationID: "validate-account-keys",
		Method:      http.MethodPost,
		Path:        basePath + "/account/validate-keys",
		Summary:     "Validate account API keys",
		Description: "Validates the provided API keys",
		Tags:        []string{"Account"},
	}, ValidateKeysHandler)
}

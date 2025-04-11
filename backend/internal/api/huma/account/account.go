// Package account provides the account management endpoints for the Huma API.
package account

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"go-crypto-bot-clean/backend/internal/core/account"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// AssetBalance represents the balance details for a single asset
type AssetBalance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
	Total  float64 `json:"total"`
	Price  float64 `json:"price,omitempty"`
}

// AccountResponse represents the response for account endpoints
type AccountResponse struct {
	Balances  map[string]AssetBalance `json:"balances"`
	UpdatedAt string                  `json:"updatedAt"`
}

// RegisterEndpoints registers the account management endpoints
func RegisterEndpoints(api huma.API, basePath string, accountService account.AccountService) {
	// GET /account/details
	huma.Register(api, huma.Operation{
		OperationID: "get-account-details",
		Method:      http.MethodGet,
		Path:        basePath + "/account/details",
		Summary:     "Get account details",
		Description: "Returns detailed account information including balances",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body AccountResponse
	}, error) {
		// Get wallet from account service
		wallet, err := accountService.GetWallet(ctx)
		if err != nil {
			return nil, err
		}

		// Convert wallet to response format
		balances := make(map[string]AssetBalance)
		for symbol, balance := range wallet.Balances {
			balances[symbol] = AssetBalance{
				Asset:  balance.Asset,
				Free:   balance.Free,
				Locked: balance.Locked,
				Total:  balance.Total,
				Price:  balance.Price,
			}
		}

		resp := &struct {
			Body AccountResponse
		}{
			Body: AccountResponse{
				Balances:  balances,
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
		}

		return resp, nil
	})

	// GET /account/wallet
	huma.Register(api, huma.Operation{
		OperationID: "get-wallet",
		Method:      http.MethodGet,
		Path:        basePath + "/account/wallet",
		Summary:     "Get wallet information",
		Description: "Returns wallet balances and information",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body models.Wallet
	}, error) {
		wallet, err := accountService.GetWallet(ctx)
		if err != nil {
			return nil, err
		}

		return &struct {
			Body models.Wallet
		}{
			Body: *wallet,
		}, nil
	})

	// POST /account/validate-keys
	huma.Register(api, huma.Operation{
		OperationID: "validate-account-keys",
		Method:      http.MethodPost,
		Path:        basePath + "/account/validate-keys",
		Summary:     "Validate account API keys",
		Description: "Validates the provided API keys",
		Tags:        []string{"Account"},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Valid   bool   `json:"valid"`
			Message string `json:"message,omitempty"`
		}
	}, error) {
		valid, err := accountService.ValidateAPIKeys(ctx)
		if err != nil {
			return nil, err
		}

		resp := &struct {
			Body struct {
				Valid   bool   `json:"valid"`
				Message string `json:"message,omitempty"`
			}
		}{}
		resp.Body.Valid = valid
		if valid {
			resp.Body.Message = "API keys are valid"
		} else {
			resp.Body.Message = "Invalid API keys"
		}

		return resp, nil
	})
}

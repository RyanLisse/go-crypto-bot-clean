//go:build integration

package http_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseURL = "http://localhost:8080" // Change if needed

func getAuthToken(t *testing.T) string {
	token := os.Getenv("CLERK_API_TOKEN")
	require.NotEmpty(t, token, "Set CLERK_API_TOKEN env var for integration tests")
	return token
}

func TestAccountWallet(t *testing.T) {
	req, err := http.NewRequest("GET", baseURL+"/api/v1/account/wallet", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t))
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Wallet: %s", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Response: %s", resp.StatusCode, string(body))
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
		return
	}
	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success=true, got: %v", result["success"])
	}
}

func TestAccountBalanceAsset(t *testing.T) {
	asset := "BTC" // You can loop over assets as needed
	req, err := http.NewRequest("GET", baseURL+"/api/v1/account/balance/"+asset+"?days=7", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t))
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Balance for %s: %s", asset, string(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Response: %s", resp.StatusCode, string(body))
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
		return
	}
	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success=true, got: %v", result["success"])
	}
}

func TestAccountRefresh(t *testing.T) {
	req, err := http.NewRequest("POST", baseURL+"/api/v1/account/refresh", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+getAuthToken(t))
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Refresh: %s", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Response: %s", resp.StatusCode, string(body))
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
		return
	}
	success, ok := result["success"].(bool)
	if !ok || !success {
		t.Errorf("Expected success=true, got: %v", result["success"])
	}
}

func TestMarketTickers(t *testing.T) {
	req, err := http.NewRequest("GET", baseURL+"/api/v1/market/tickers", nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Tickers: %s", string(body))
	var result map[string]interface{}
	_ = json.Unmarshal(body, &result)
	assert.True(t, result["success"].(bool))
}

func TestMarketTickerSymbol(t *testing.T) {
	symbol := "BTCUSDT" // Loop over more symbols if desired
	req, err := http.NewRequest("GET", baseURL+"/api/v1/market/ticker/"+symbol, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Ticker for %s: %s", symbol, string(body))
	var result map[string]interface{}
	_ = json.Unmarshal(body, &result)
	assert.True(t, result["success"].(bool))
}

func TestAccountTest(t *testing.T) {
	req, err := http.NewRequest("GET", baseURL+"/api/v1/account-test", nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Account test: %s", string(body))
	var result map[string]interface{}
	_ = json.Unmarshal(body, &result)
	assert.True(t, result["success"].(bool))
}

func TestAccountWalletTest(t *testing.T) {
	req, err := http.NewRequest("GET", baseURL+"/api/v1/account-wallet-test", nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	t.Logf("Account wallet test: %s", string(body))
	var result map[string]interface{}
	_ = json.Unmarshal(body, &result)
	assert.True(t, result["success"].(bool))
}

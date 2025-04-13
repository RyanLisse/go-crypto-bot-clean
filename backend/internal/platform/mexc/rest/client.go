package rest

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/neo/crypto-bot/pkg/ratelimiter"
)

const (
	// MEXC API endpoints
	BaseURL     = "https://api.mexc.com"
	SpotBaseURL = BaseURL + "/api/v3"

	// Rate limits (based on MEXC documentation)
	// Spot API rate limits
	SpotPublicRequestsPerMinute  = 1200 // 20 requests per second
	SpotPrivateRequestsPerMinute = 600  // 10 requests per second

	// Default values
	DefaultTimeout = 10 * time.Second
)

// Client implements the MEXC REST API client
type Client struct {
	httpClient         *http.Client
	baseURL            string
	apiKey             string
	secretKey          string
	publicRateLimiter  *ratelimiter.TokenBucket
	privateRateLimiter *ratelimiter.TokenBucket
}

// ClientOption defines a functional option for configuring the Client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new MEXC REST client
func NewClient(apiKey, secretKey string, options ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:   SpotBaseURL,
		apiKey:    apiKey,
		secretKey: secretKey,
		// Initialize rate limiters (tokens per second)
		publicRateLimiter:  ratelimiter.NewTokenBucket(SpotPublicRequestsPerMinute/60.0, 50),
		privateRateLimiter: ratelimiter.NewTokenBucket(SpotPrivateRequestsPerMinute/60.0, 25),
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client
}

// callPublicAPI makes a request to a public API endpoint
func (c *Client) callPublicAPI(ctx context.Context, method, path string, params map[string]string) ([]byte, error) {
	// Apply rate limiting
	c.publicRateLimiter.Wait()

	// Construct URL with query parameters
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	if params != nil && len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	return body, nil
}

// callPrivateAPI makes a request to a private API endpoint requiring authentication
func (c *Client) callPrivateAPI(ctx context.Context, method, path string, params map[string]string, body interface{}) ([]byte, error) {
	// Apply rate limiting
	c.privateRateLimiter.Wait()

	// Add timestamp parameter for signature
	if params == nil {
		params = make(map[string]string)
	}
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	params["timestamp"] = timestamp

	// Create request
	url := c.baseURL + path
	var reqBody io.Reader

	// Handle request body for POST/PUT methods
	var jsonBody []byte
	var err error
	if body != nil && (method == http.MethodPost || method == http.MethodPut) {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters and calculate signature
	q := req.URL.Query()
	queryString := ""

	// Add all parameters to query string for signature
	for k, v := range params {
		q.Add(k, v)
	}
	queryString = q.Encode()

	// Calculate HMAC SHA256 signature
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write([]byte(queryString))
	signature := hex.EncodeToString(h.Sum(nil))

	// Add signature to query parameters
	q.Add("signature", signature)
	req.URL.RawQuery = q.Encode()

	// Set headers
	req.Header.Set("X-MBX-APIKEY", c.apiKey)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s, status code: %d", string(respBody), resp.StatusCode)
	}

	return respBody, nil
}

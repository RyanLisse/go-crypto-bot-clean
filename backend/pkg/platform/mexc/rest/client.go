package rest

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc/apikeystore"
	"github.com/cenkalti/backoff/v4"
	"golang.org/x/time/rate"
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
	DefaultKeyID   = "default"

	// Error types
	ErrInvalidResponse    = "invalid_response"
	ErrRateLimit          = "rate_limit"
	ErrAuth               = "authentication"
	ErrNetwork            = "network"
	ErrServer             = "server"
	ErrInvalidRequest     = "invalid_request"
	ErrInsufficientFunds  = "insufficient_funds"
	ErrOrderNotFound      = "order_not_found"
	ErrSymbolNotFound     = "symbol_not_found"
	ErrInvalidOrderStatus = "invalid_order_status"
	ErrUnknown            = "unknown"
)

// APIError represents an error from the MEXC API
type APIError struct {
	Code       int    `json:"code"`
	Message    string `json:"msg"`
	ErrorType  string
	StatusCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("MEXC API error: %d - %s (type: %s, HTTP status: %d)",
		e.Code, e.Message, e.ErrorType, e.StatusCode)
}

// IsRetryable determines if an API error is retryable
func (e *APIError) IsRetryable() bool {
	switch e.ErrorType {
	case ErrRateLimit, ErrNetwork, ErrServer:
		return true
	default:
		return false
	}
}

// Client implements the MEXC REST API client
// Note: MEXC API requires the APIKEY header (not X-MBX-APIKEY) for authentication
type Client struct {
	httpClient         *http.Client
	baseURL            string
	keyID              string
	keyStore           apikeystore.KeyStore
	publicRateLimiter  *rate.Limiter
	privateRateLimiter *rate.Limiter
	backoffStrategy    backoff.BackOff
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

// WithKeyID sets the key ID to use from the key store
func WithKeyID(keyID string) ClientOption {
	return func(c *Client) {
		c.keyID = keyID
	}
}

// WithBackoffStrategy sets a custom backoff strategy
func WithBackoffStrategy(strategy backoff.BackOff) ClientOption {
	return func(c *Client) {
		c.backoffStrategy = strategy
	}
}

// NewClient creates a new MEXC REST client with direct API key
func NewClient(apiKey, secretKey string, options ...ClientOption) *Client {
	// Create a memory key store with the provided credentials
	keyStore := apikeystore.NewMemoryKeyStore()
	keyStore.SetAPIKey(DefaultKeyID, &apikeystore.APIKeyCredentials{
		APIKey:    apiKey,
		SecretKey: secretKey,
	})

	return NewClientWithKeyStore(keyStore, DefaultKeyID, options...)
}

// NewClientWithKeyStore creates a new client with the provided key store
func NewClientWithKeyStore(keyStore apikeystore.KeyStore, keyID string, options ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL:  SpotBaseURL,
		keyID:    keyID,
		keyStore: keyStore,
		// Initialize rate limiters (tokens per second)
		publicRateLimiter:  rate.NewLimiter(rate.Limit(SpotPublicRequestsPerMinute/60.0), 50),
		privateRateLimiter: rate.NewLimiter(rate.Limit(SpotPrivateRequestsPerMinute/60.0), 25),
		// Use exponential backoff as default
		backoffStrategy: backoff.NewExponentialBackOff(),
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client
}

// getCredentials retrieves API credentials from the key store
func (c *Client) getCredentials() (*apikeystore.APIKeyCredentials, error) {
	creds, err := c.keyStore.GetAPIKey(c.keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API credentials: %w", err)
	}
	return creds, nil
}

// callPublicAPI makes a request to a public API endpoint with retries
func (c *Client) callPublicAPI(ctx context.Context, method, path string, params map[string]string) ([]byte, error) {
	var result []byte
	operation := func() error {
		respBody, err := c.doPublicAPICall(ctx, method, path, params)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && apiErr.IsRetryable() {
				return err // backoff will retry
			}
			return backoff.Permanent(err) // do not retry
		}
		result = respBody
		return nil
	}
	err := backoff.Retry(operation, backoff.WithContext(c.backoffStrategy, ctx))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// doPublicAPICall makes a single request to a public API endpoint
func (c *Client) doPublicAPICall(ctx context.Context, method, path string, params map[string]string) ([]byte, error) {
	// Apply rate limiting
	if !c.publicRateLimiter.Allow() {
		return nil, &APIError{
			Message:    "rate limit exceeded",
			ErrorType:  ErrRateLimit,
			StatusCode: 0,
		}
	}

	// Construct URL with query parameters
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to execute request: %v", err),
			ErrorType:  ErrNetwork,
			StatusCode: 0,
		}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to read response body: %v", err),
			ErrorType:  ErrNetwork,
			StatusCode: resp.StatusCode,
		}
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		apiErr := parseAPIError(body, resp.StatusCode)
		return nil, apiErr
	}

	return body, nil
}

// callPrivateAPI makes a request to a private API endpoint requiring authentication with retries
func (c *Client) callPrivateAPI(ctx context.Context, method, path string, params map[string]string, body interface{}) ([]byte, error) {
	var result []byte
	operation := func() error {
		respBody, err := c.doPrivateAPICall(ctx, method, path, params, body)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && apiErr.IsRetryable() {
				return err // backoff will retry
			}
			return backoff.Permanent(err) // do not retry
		}
		result = respBody
		return nil
	}
	err := backoff.Retry(operation, backoff.WithContext(c.backoffStrategy, ctx))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// doPrivateAPICall makes a single request to a private API endpoint requiring authentication
func (c *Client) doPrivateAPICall(ctx context.Context, method, path string, params map[string]string, body interface{}) ([]byte, error) {
	// Apply rate limiting
	_ = c.privateRateLimiter.Wait(ctx)

	// Get API credentials
	creds, err := c.getCredentials()
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to get API credentials: %v", err),
			ErrorType:  ErrAuth,
			StatusCode: 0,
		}
	}

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
	if body != nil && (method == http.MethodPost || method == http.MethodPut) {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, &APIError{
				Message:    fmt.Sprintf("failed to marshal request body: %v", err),
				ErrorType:  ErrInvalidRequest,
				StatusCode: 0,
			}
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to create request: %v", err),
			ErrorType:  ErrInvalidRequest,
			StatusCode: 0,
		}
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
	h := hmac.New(sha256.New, []byte(creds.SecretKey))
	h.Write([]byte(queryString))
	signature := hex.EncodeToString(h.Sum(nil))

	// Add signature to query parameters
	q.Add("signature", signature)
	req.URL.RawQuery = q.Encode()

	// Set headers
	req.Header.Set("APIKEY", creds.APIKey)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to execute request: %v", err),
			ErrorType:  ErrNetwork,
			StatusCode: 0,
		}
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &APIError{
			Message:    fmt.Sprintf("failed to read response body: %v", err),
			ErrorType:  ErrNetwork,
			StatusCode: resp.StatusCode,
		}
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		apiErr := parseAPIError(respBody, resp.StatusCode)
		return nil, apiErr
	}

	return respBody, nil
}

// parseAPIError parses an API error response
func parseAPIError(body []byte, statusCode int) *APIError {
	var errResp struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}

	if err := json.Unmarshal(body, &errResp); err != nil {
		// Couldn't parse the error response
		return &APIError{
			Message:    fmt.Sprintf("HTTP error %d: %s", statusCode, string(body)),
			ErrorType:  ErrUnknown,
			StatusCode: statusCode,
		}
	}

	// Determine error type based on code and status
	errorType := ErrUnknown
	switch {
	case statusCode == 429:
		errorType = ErrRateLimit
	case statusCode == 401 || statusCode == 403:
		errorType = ErrAuth
	case statusCode >= 500:
		errorType = ErrServer
	case statusCode >= 400 && statusCode < 500:
		// Map common error codes
		switch errResp.Code {
		case -1121, -1122:
			errorType = ErrInvalidRequest
		case -2010, -2011:
			errorType = ErrInsufficientFunds
		case -2013:
			errorType = ErrOrderNotFound
		case -1100, -1101, -1102, -1103:
			errorType = ErrInvalidRequest
		case -1104, -1105, -1106:
			errorType = ErrInvalidRequest
		}
	}

	return &APIError{
		Code:       errResp.Code,
		Message:    errResp.Message,
		ErrorType:  errorType,
		StatusCode: statusCode,
	}
}

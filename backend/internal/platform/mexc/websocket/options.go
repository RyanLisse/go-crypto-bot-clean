package websocket

import (
	"github.com/gorilla/websocket"
	"go-crypto-bot-clean/backend/pkg/ratelimiter"
)

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithEndpoint sets the WebSocket endpoint URL
func WithEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

// WithDialer sets a custom WebSocket dialer
func WithDialer(dialer *websocket.Dialer) ClientOption {
	return func(c *Client) {
		c.dialer = dialer
	}
}

// WithCredentials sets the API credentials
func WithCredentials(apiKey, secretKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
		c.secretKey = secretKey
	}
}

// WithConnRateLimiter sets a custom connection rate limiter
func WithConnRateLimiter(limiter *ratelimiter.TokenBucketRateLimiter) ClientOption {
	return func(c *Client) {
		c.connRateLimiter = limiter
	}
}

// WithSubRateLimiter sets a custom subscription rate limiter
func WithSubRateLimiter(limiter *ratelimiter.TokenBucketRateLimiter) ClientOption {
	return func(c *Client) {
		c.subRateLimiter = limiter
	}
}

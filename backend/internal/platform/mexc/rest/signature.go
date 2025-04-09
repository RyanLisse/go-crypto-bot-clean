package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strconv"
	"time"
)

// SignRequest adds timestamp and signature to the request parameters
func SignRequest(params url.Values, secretKey string) url.Values {
	if params == nil {
		params = url.Values{}
	}

	// Add timestamp if not already present
	if params.Get("timestamp") == "" {
		params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	}

	// Create signature
	signature := CreateSignature(params.Encode(), secretKey)
	params.Set("signature", signature)

	return params
}

// CreateSignature creates an HMAC SHA256 signature for the given payload
func CreateSignature(payload, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateQueryWithSignature creates a query string with timestamp and signature
func CreateQueryWithSignature(params url.Values, secretKey string) string {
	signedParams := SignRequest(params, secretKey)
	return signedParams.Encode()
}

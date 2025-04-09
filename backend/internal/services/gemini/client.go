package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// GeminiClient is a client for the Gemini API
type GeminiClient struct {
	APIKey     string
	Endpoint   string
	HTTPClient *http.Client
}

// GeminiRequest represents a request to the Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents the content of a Gemini request
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of a Gemini content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents a response from the Gemini API
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
	Error      *GeminiError      `json:"error,omitempty"`
}

// GeminiCandidate represents a candidate response from the Gemini API
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// GeminiError represents an error from the Gemini API
type GeminiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NewGeminiClient creates a new GeminiClient
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		APIKey:     apiKey,
		Endpoint:   "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent",
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// AnalyzeMetrics analyzes performance metrics using the Gemini API
func (g *GeminiClient) AnalyzeMetrics(ctx context.Context, metrics models.PerformanceReportMetrics) (string, error) {
	// Convert metrics to JSON for better formatting
	metricsJSON, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal metrics: %w", err)
	}

	// Create the prompt
	prompt := fmt.Sprintf(`Analyze these trading bot performance metrics and provide key insights:
%s

Please provide:
1. A summary of the overall performance
2. Key metrics that stand out (positive or negative)
3. Potential areas for improvement
4. Any anomalies or concerning patterns
5. Recommendations for optimization

Format your response in a clear, professional manner suitable for a performance report.`, string(metricsJSON))

	// Create the request
	req := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	// Marshal the request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.Endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", g.APIKey)

	// Send the request
	resp, err := g.HTTPClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", string(body))
	}

	// Unmarshal the response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API error
	if geminiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", geminiResp.Error.Message)
	}

	// Check for empty response
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	// Return the analysis
	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// ExtractInsights extracts insights from an analysis
func (g *GeminiClient) ExtractInsights(ctx context.Context, analysis string) ([]string, error) {
	// Create the prompt
	prompt := fmt.Sprintf(`Extract the key insights from this performance analysis as a list of concise bullet points:
%s

Format your response as a list of insights, one per line, without numbering or bullet points.`, analysis)

	// Create the request
	req := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	// Marshal the request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, g.Endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", g.APIKey)

	// Send the request
	resp, err := g.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	// Unmarshal the response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API error
	if geminiResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", geminiResp.Error.Message)
	}

	// Check for empty response
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	// Split the response into lines
	insightsText := geminiResp.Candidates[0].Content.Parts[0].Text
	insights := []string{}

	// Split by newline and filter empty lines
	for _, line := range bytes.Split([]byte(insightsText), []byte("\n")) {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) > 0 {
			insights = append(insights, string(trimmed))
		}
	}

	return insights, nil
}

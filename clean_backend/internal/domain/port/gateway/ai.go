package gateway

import "context"

// AIGateway defines the port for interacting with an AI service (e.g., Gemini).
type AIGateway interface {
	// GenerateText generates text based on a given prompt.
	GenerateText(ctx context.Context, prompt string) (string, error)

	// AnalyzeSentiment analyzes the sentiment of a given text.
	AnalyzeSentiment(ctx context.Context, text string) (string, error) // Result could be e.g., "Positive", "Negative", "Neutral"

	// TODO: Add other specific AI interactions as needed (e.g., data analysis, specific model calls).
}

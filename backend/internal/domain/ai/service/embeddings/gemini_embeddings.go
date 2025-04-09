package embeddings

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// GeminiEmbeddingsService implements the EmbeddingsService interface using Google Gemini
type GeminiEmbeddingsService struct {
	client *genai.Client
	model  string
	logger *zap.Logger
}

// NewGeminiEmbeddingsService creates a new Gemini embeddings service
func NewGeminiEmbeddingsService(logger *zap.Logger) (*GeminiEmbeddingsService, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Create Gemini client
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Get model from environment variable or use default
	model := os.Getenv("GEMINI_EMBEDDINGS_MODEL")
	if model == "" {
		model = "gemini-embedding-exp-03-07" // Default embedding model
	}

	return &GeminiEmbeddingsService{
		client: client,
		model:  model,
		logger: logger,
	}, nil
}

// GenerateEmbedding generates an embedding for the given text
func (s *GeminiEmbeddingsService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Clean and prepare text
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty text")
	}

	// Create embedding model
	em := s.client.EmbeddingModel(s.model)

	// Generate embedding
	resp, err := em.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		s.logger.Error("Failed to generate Gemini embedding",
			zap.Error(err),
			zap.String("model", s.model),
			zap.Int("text_length", len(text)),
		)
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Check if we got any embeddings
	if resp.Embedding == nil || len(resp.Embedding.Values) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	// Return the embedding
	return resp.Embedding.Values, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (s *GeminiEmbeddingsService) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	// Clean and prepare texts
	var validTexts []string
	for _, text := range texts {
		text = strings.TrimSpace(text)
		if text != "" {
			validTexts = append(validTexts, text)
		}
	}

	if len(validTexts) == 0 {
		return nil, fmt.Errorf("no valid texts")
	}

	// Create embedding model
	em := s.client.EmbeddingModel(s.model)

	// Generate embeddings one by one (Gemini doesn't support batch embeddings yet)
	embeddings := make([][]float32, len(validTexts))
	for i, text := range validTexts {
		resp, err := em.EmbedContent(ctx, genai.Text(text))
		if err != nil {
			s.logger.Error("Failed to generate Gemini embedding in batch",
				zap.Error(err),
				zap.String("model", s.model),
				zap.Int("text_length", len(text)),
				zap.Int("index", i),
			)
			return nil, fmt.Errorf("failed to generate embedding at index %d: %w", i, err)
		}

		// Check if we got any embeddings
		if resp.Embedding == nil || len(resp.Embedding.Values) == 0 {
			return nil, fmt.Errorf("no embeddings returned at index %d", i)
		}

		embeddings[i] = resp.Embedding.Values
	}

	return embeddings, nil
}

// Close closes the Gemini client
func (s *GeminiEmbeddingsService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

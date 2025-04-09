package embeddings

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// OpenAIEmbeddingsService implements the EmbeddingsService interface using OpenAI
type OpenAIEmbeddingsService struct {
	client *openai.Client
	logger *zap.Logger
	model  string
}

// NewOpenAIEmbeddingsService creates a new OpenAI embeddings service
func NewOpenAIEmbeddingsService(logger *zap.Logger) (*OpenAIEmbeddingsService, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)

	// Get model from environment variable or use default
	model := os.Getenv("OPENAI_EMBEDDINGS_MODEL")
	if model == "" {
		model = openai.AdaEmbeddingV2 // Default to text-embedding-ada-002
	}

	return &OpenAIEmbeddingsService{
		client: client,
		logger: logger,
		model:  model,
	}, nil
}

// GenerateEmbedding generates an embedding for the given text
func (s *OpenAIEmbeddingsService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Clean and prepare text
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty text")
	}

	// Create embedding request
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: s.model,
	}

	// Generate embedding
	resp, err := s.client.CreateEmbeddings(ctx, req)
	if err != nil {
		s.logger.Error("Failed to generate embedding",
			zap.Error(err),
			zap.String("model", s.model),
			zap.Int("text_length", len(text)),
		)
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Check if we got any embeddings
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	// Return the embedding
	return resp.Data[0].Embedding, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (s *OpenAIEmbeddingsService) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
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

	// Create embedding request
	req := openai.EmbeddingRequest{
		Input: validTexts,
		Model: s.model,
	}

	// Generate embeddings
	resp, err := s.client.CreateEmbeddings(ctx, req)
	if err != nil {
		s.logger.Error("Failed to generate batch embeddings",
			zap.Error(err),
			zap.String("model", s.model),
			zap.Int("text_count", len(validTexts)),
		)
		return nil, fmt.Errorf("failed to generate batch embeddings: %w", err)
	}

	// Check if we got any embeddings
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	// Sort embeddings by index to maintain original order
	embeddings := make([][]float32, len(resp.Data))
	for _, data := range resp.Data {
		embeddings[data.Index] = data.Embedding
	}

	return embeddings, nil
}

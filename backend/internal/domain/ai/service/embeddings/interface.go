package embeddings

import (
	"context"
)

// EmbeddingsService defines the interface for generating embeddings
type EmbeddingsService interface {
	// GenerateEmbedding generates an embedding for the given text
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)

	// GenerateBatchEmbeddings generates embeddings for multiple texts
	GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
}

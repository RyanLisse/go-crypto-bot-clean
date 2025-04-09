package embeddings

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// FallbackEmbeddingsService implements the EmbeddingsService interface using multiple providers
type FallbackEmbeddingsService struct {
	primaryService   EmbeddingsService
	fallbackService  EmbeddingsService
	logger           *zap.Logger
	useFallbackCount int // For monitoring how often fallback is used
}

// NewFallbackEmbeddingsService creates a new fallback embeddings service
func NewFallbackEmbeddingsService(
	primaryService EmbeddingsService,
	fallbackService EmbeddingsService,
	logger *zap.Logger,
) *FallbackEmbeddingsService {
	return &FallbackEmbeddingsService{
		primaryService:   primaryService,
		fallbackService:  fallbackService,
		logger:           logger,
		useFallbackCount: 0,
	}
}

// GenerateEmbedding generates an embedding for the given text
func (s *FallbackEmbeddingsService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Try primary service first
	embedding, err := s.primaryService.GenerateEmbedding(ctx, text)
	if err == nil {
		return embedding, nil
	}

	// Log the error and try fallback service
	s.logger.Warn("Primary embedding service failed, using fallback",
		zap.Error(err),
		zap.Int("text_length", len(text)),
	)

	// Increment fallback counter
	s.useFallbackCount++

	// Try fallback service
	embedding, err = s.fallbackService.GenerateEmbedding(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("both primary and fallback embedding services failed: %w", err)
	}

	return embedding, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (s *FallbackEmbeddingsService) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	// Try primary service first
	embeddings, err := s.primaryService.GenerateBatchEmbeddings(ctx, texts)
	if err == nil {
		return embeddings, nil
	}

	// Log the error and try fallback service
	s.logger.Warn("Primary embedding service failed for batch, using fallback",
		zap.Error(err),
		zap.Int("text_count", len(texts)),
	)

	// Increment fallback counter
	s.useFallbackCount++

	// Try fallback service
	embeddings, err = s.fallbackService.GenerateBatchEmbeddings(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("both primary and fallback embedding services failed for batch: %w", err)
	}

	return embeddings, nil
}

// GetFallbackCount returns the number of times the fallback service was used
func (s *FallbackEmbeddingsService) GetFallbackCount() int {
	return s.useFallbackCount
}

// ResetFallbackCount resets the fallback counter
func (s *FallbackEmbeddingsService) ResetFallbackCount() {
	s.useFallbackCount = 0
}

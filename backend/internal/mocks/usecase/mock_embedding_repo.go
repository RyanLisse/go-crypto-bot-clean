package mocks

import (
	"context"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

type MockEmbeddingRepository struct{}

func (m *MockEmbeddingRepository) SaveEmbedding(ctx context.Context, embedding *model.AIEmbedding) error {
	return nil
}
func (m *MockEmbeddingRepository) GetEmbedding(ctx context.Context, embeddingID string) (*model.AIEmbedding, error) {
	return &model.AIEmbedding{}, nil
}
func (m *MockEmbeddingRepository) DeleteEmbedding(ctx context.Context, embeddingID string) error {
	return nil
}
func (m *MockEmbeddingRepository) SearchEmbeddings(ctx context.Context, queryVector []float64, limit int) ([]*model.AIEmbedding, error) {
	return []*model.AIEmbedding{}, nil
}
func (m *MockEmbeddingRepository) FindSimilar(ctx context.Context, vector []float64, limit int) ([]*model.AIEmbedding, error) {
	return m.SearchEmbeddings(ctx, vector, limit)
}

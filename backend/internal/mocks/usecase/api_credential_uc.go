package mocks

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

// MockAPICredentialUseCase is a mock implementation of the APICredentialUseCase interface
type MockAPICredentialUseCase struct {
	mock.Mock
}

func (m *MockAPICredentialUseCase) CreateCredential(ctx context.Context, credential *model.APICredential) error {
	args := m.Called(ctx, credential)
	return args.Error(0)
}

func (m *MockAPICredentialUseCase) GetCredential(ctx context.Context, id string) (*model.APICredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialUseCase) UpdateCredential(ctx context.Context, credential *model.APICredential) error {
	args := m.Called(ctx, credential)
	return args.Error(0)
}

func (m *MockAPICredentialUseCase) DeleteCredential(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAPICredentialUseCase) ListCredentials(ctx context.Context, userID string) ([]*model.APICredential, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialUseCase) GetCredentialByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	args := m.Called(ctx, userID, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

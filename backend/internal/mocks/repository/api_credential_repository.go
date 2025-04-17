package mocks

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/stretchr/testify/mock"
)

// MockAPICredentialRepository is a mock implementation of the port.APICredentialRepository interface
type MockAPICredentialRepository struct {
	mock.Mock
}

func (m *MockAPICredentialRepository) ListAll(ctx context.Context) ([]*model.APICredential, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) Save(ctx context.Context, credential *model.APICredential) error {
	args := m.Called(ctx, credential)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) GetByID(ctx context.Context, id string) (*model.APICredential, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) GetByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	args := m.Called(ctx, userID, exchange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) GetByUserIDAndLabel(ctx context.Context, userID, exchange, label string) (*model.APICredential, error) {
	args := m.Called(ctx, userID, exchange, label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) DeleteByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) ListByUserID(ctx context.Context, userID string) ([]*model.APICredential, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.APICredential), args.Error(1)
}

func (m *MockAPICredentialRepository) UpdateStatus(ctx context.Context, id string, status model.APICredentialStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) UpdateLastUsed(ctx context.Context, id string, lastUsed time.Time) error {
	args := m.Called(ctx, id, lastUsed)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) UpdateLastVerified(ctx context.Context, id string, lastVerified time.Time) error {
	args := m.Called(ctx, id, lastVerified)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) IncrementFailureCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAPICredentialRepository) ResetFailureCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

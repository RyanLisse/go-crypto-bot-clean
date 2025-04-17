#!/usr/bin/env bash
set -e

# This script consolidates all mock files and fixes import paths
# Run from the backend directory: bash scripts/fix_mocks_consolidated.sh

echo "Creating directories for consolidated mocks..."
mkdir -p internal/mocks/repository
mkdir -p internal/mocks/service
mkdir -p internal/mocks/usecase

# First, standardize mock implementations in a single location
echo "Consolidating MockAPICredentialRepository..."
cat > internal/mocks/repository/api_credential_repository.go << 'EOF'
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
EOF

echo "Consolidating MockEncryptionService..."
cat > internal/mocks/service/encryption_service.go << 'EOF'
package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockEncryptionService is a mock implementation of the crypto.EncryptionService interface
type MockEncryptionService struct {
	mock.Mock
}

func (m *MockEncryptionService) Encrypt(plaintext string) ([]byte, error) {
	args := m.Called(plaintext)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEncryptionService) Decrypt(ciphertext []byte) (string, error) {
	args := m.Called(ciphertext)
	return args.String(0), args.Error(1)
}
EOF

echo "Consolidating MockAPICredentialUseCase..."
cat > internal/mocks/usecase/api_credential_uc.go << 'EOF'
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
EOF

# Fix import paths in test files - credential_cache_service_test.go
echo "Fixing imports in credential_cache_service_test.go..."
cat > /tmp/fixed_imports.sed << 'EOF'
/import (/a\
	mockRepo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"
EOF

sed -i '' -f /tmp/fixed_imports.sed internal/domain/service/credential_cache_service_test.go

# Replace MockAPICredentialRepository with mockRepo.MockAPICredentialRepository
sed -i '' 's/MockAPICredentialRepository/mockRepo.MockAPICredentialRepository/g' internal/domain/service/credential_cache_service_test.go

# Fix import paths in test files - credential_encryption_service_test.go
echo "Fixing imports in credential_encryption_service_test.go..."
cat > /tmp/fixed_imports2.sed << 'EOF'
/import (/a\
	mockRepo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"\
	mockService "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/service"
EOF

sed -i '' -f /tmp/fixed_imports2.sed internal/domain/service/credential_encryption_service_test.go

# Replace MockAPICredentialRepository with mockRepo.MockAPICredentialRepository
sed -i '' 's/MockAPICredentialRepository/mockRepo.MockAPICredentialRepository/g' internal/domain/service/credential_encryption_service_test.go
sed -i '' 's/MockEncryptionService/mockService.MockEncryptionService/g' internal/domain/service/credential_encryption_service_test.go

# Fix import paths in test files - api_credential_handler_test.go
echo "Fixing imports in api_credential_handler_test.go..."
cat > /tmp/fixed_imports3.sed << 'EOF'
/import (/a\
	mockUsecase "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"
EOF

sed -i '' -f /tmp/fixed_imports3.sed internal/adapter/delivery/http/handler/api_credential_handler_test.go

# Replace MockAPICredentialUseCase with mockUsecase.MockAPICredentialUseCase
sed -i '' 's/mocks\.MockAPICredentialUseCase/mockUsecase.MockAPICredentialUseCase/g' internal/adapter/delivery/http/handler/api_credential_handler_test.go

# Fix import paths in test files - credential_fallback_service_test.go
echo "Fixing imports in credential_fallback_service_test.go..."
cat > /tmp/fixed_imports4.sed << 'EOF'
/import (/a\
	mockRepo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"
EOF

sed -i '' -f /tmp/fixed_imports4.sed internal/domain/service/credential_fallback_service_test.go

# Replace MockAPICredentialRepository with mockRepo.MockAPICredentialRepository
sed -i '' 's/MockAPICredentialRepository/mockRepo.MockAPICredentialRepository/g' internal/domain/service/credential_fallback_service_test.go

# Fix import paths in test files - credential_manager_test.go
echo "Fixing imports in credential_manager_test.go..."
cat > /tmp/fixed_imports5.sed << 'EOF'
/import (/a\
	mockRepo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"\
	mockService "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/service"
EOF

sed -i '' -f /tmp/fixed_imports5.sed internal/domain/service/credential_manager_test.go

# Replace mocks with the correct imports
sed -i '' 's/MockAPICredentialRepository/mockRepo.MockAPICredentialRepository/g' internal/domain/service/credential_manager_test.go
sed -i '' 's/MockEncryptionService/mockService.MockEncryptionService/g' internal/domain/service/credential_manager_test.go

# Fix ai_handler_test.go for duplicate imports
echo "Fixing imports in ai_handler_test.go..."
sed -i '' '/mocks "github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"/d' internal/adapter/delivery/http/handler/ai_handler_test.go

# Fix usecase_factory.go
echo "Fixing imports in usecase_factory.go..."
cat > /tmp/fixed_imports6.sed << 'EOF'
/import (/a\
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"
EOF

sed -i '' -f /tmp/fixed_imports6.sed internal/factory/usecase_factory.go
sed -i '' 's/mocks\./usecase.Mock/g' internal/factory/usecase_factory.go

# Fix all other test files with duplicate mocks imports
echo "Fixing other test files with duplicate imports..."
find internal -name "*_test.go" -exec sed -i '' '/mocks "github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"/d' {} \;

# Final cleanup
echo "Removing temporary files..."
rm -f /tmp/fixed_imports*.sed

echo "Done! Please run 'go test ./...' to verify fixes." 
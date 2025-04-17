#!/bin/bash
set -e

echo "Fixing additional import issues in test files..."

# Fix duplicate imports in usecase_factory.go
echo "Fixing usecase_factory.go"
sed -i '' 's|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase"|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/usecase" // import usecase mocks|' internal/factory/usecase_factory.go
sed -i '' '/"github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase" \/\/ import usecase mocks/!s/"github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"//g' internal/factory/usecase_factory.go
sed -i '' 's|usecase\.MockMockTradeUseCase|mocks.MockTradeUseCase|g' internal/factory/usecase_factory.go
sed -i '' 's|usecase\.MockMockPositionUseCase|mocks.MockPositionUseCase|g' internal/factory/usecase_factory.go
sed -i '' 's|usecase\.MockMockStatusUseCase|mocks.MockStatusUseCase|g' internal/factory/usecase_factory.go

# Fix import issues in api_credential_handler_test.go
echo "Fixing api_credential_handler_test.go"
sed -i '' '/"github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"/!s/"github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"//g' internal/adapter/delivery/http/handler/api_credential_handler_test.go
sed -i '' '/"github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase\/mocks"/d' internal/adapter/delivery/http/handler/api_credential_handler_test.go

# Fix duplicate imports in domain/service test files
echo "Fixing credential_cache_service_test.go"
sed -i '' '2,3d' internal/domain/service/credential_cache_service_test.go
echo "Fixing credential_encryption_service_test.go"
sed -i '' '2,3d' internal/domain/service/credential_encryption_service_test.go
echo "Fixing credential_fallback_service_test.go" 
sed -i '' '2,3d' internal/domain/service/credential_fallback_service_test.go
echo "Fixing credential_manager_test.go"
sed -i '' '2,3d' internal/domain/service/credential_manager_test.go

# Fix references to mockRepo.mockRepo
echo "Fixing mockRepo references in credential_cache_service_test.go"
sed -i '' 's|mockRepo\.mockRepo|mockRepo\.MockAPICredentialRepository|g' internal/domain/service/credential_cache_service_test.go

# Fix import issues in mock_ai_usecase_test.go
echo "Fixing mock_ai_usecase_test.go"
sed -i '' '/"github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"/d' internal/mocks/adapter/delivery/http/handler/mock_ai_usecase_test.go

# Fix internal/usecase/trade_uc_test.go
echo "Fixing trade_uc_test.go"
sed -i '' 's|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks"|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/service"|g' internal/usecase/trade_uc_test.go

echo "All import issues fixed!" 
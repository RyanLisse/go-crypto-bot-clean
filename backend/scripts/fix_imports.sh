#!/bin/bash
set -e

echo "Fixing import issues in test files..."

# Fix duplicate imports in api_credential_handler_test.go
echo "Fixing api_credential_handler_test.go"
sed -i '' '/^import (/,/^)/ {
    /mockUsecase "github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"/d
}' internal/adapter/delivery/http/handler/api_credential_handler_test.go

# Fix issues in credential_fallback_service_test.go
echo "Fixing credential_fallback_service_test.go"
sed -i '' '/^package service/a \
import (\
\t"context"\
\t"os"\
\t"testing"\
\t"time"\
\
\t"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"\
\t"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/service"\
\t"github.com/rs/zerolog"\
\t"github.com/stretchr/testify/assert"\
\t"github.com/stretchr/testify/mock"\
)' internal/domain/service/credential_fallback_service_test.go
sed -i '' '2,3d' internal/domain/service/credential_fallback_service_test.go

# Fix issues in credential_manager_test.go
echo "Fixing credential_manager_test.go"
sed -i '' '/^package service/a \
import (\
\t"context"\
\t"os"\
\t"testing"\
\t"time"\
\
\t"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"\
\t"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/service"\
\t"github.com/rs/zerolog"\
\t"github.com/stretchr/testify/assert"\
\t"github.com/stretchr/testify/mock"\
)' internal/domain/service/credential_manager_test.go
sed -i '' '2,3d' internal/domain/service/credential_manager_test.go

# Fix issues in credential_cache_service_test.go
echo "Fixing credential_cache_service_test.go"
sed -i '' 's|repoMocks "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/service"|g' internal/domain/service/credential_cache_service_test.go
sed -i '' 's|mockRepo\.MockAPICredentialRepository|service.MockAPICredentialRepository|g' internal/domain/service/credential_cache_service_test.go
sed -i '' 's|new(repoMocks.MockAPICredentialRepository)|new(service.MockAPICredentialRepository)|g' internal/domain/service/credential_cache_service_test.go

# Fix issues in market_data_test.go (duplicate mocks imports)
echo "Fixing market_data_test.go"
sed -i '' '/^import (/,/^)/ {
    /mocks "github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"/d
}' internal/usecase/market_data_test.go

# Fix issues in newcoin_uc_test.go (duplicate mocks imports)
echo "Fixing newcoin_uc_test.go"
sed -i '' '/^import (/,/^)/ {
    /mocks "github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"/d
}' internal/usecase/newcoin_uc_test.go

# Fix issues in position_uc_test.go (duplicate mocks imports)
echo "Fixing position_uc_test.go"
sed -i '' '/^import (/,/^)/ {
    /mocks "github.com\/RyanLisse\/go-crypto-bot-clean\/backend\/internal\/mocks\/usecase"/d
}' internal/usecase/position_uc_test.go

# Fix issues in trade_uc_test.go (undefined mocks.MockMEXCClient)
echo "Fixing trade_uc_test.go"
sed -i '' 's|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/service"|"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/domain/mocks"|g' internal/usecase/trade_uc_test.go

echo "All import issues fixed!" 
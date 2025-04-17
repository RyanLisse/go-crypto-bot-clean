#!/bin/bash
set -e

echo "Fixing credential service test files..."

# Define the correct import block to use
IMPORT_BLOCK='import (
\t"context"
\t"errors"
\t"os"
\t"testing"
\t"time"

\t"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
\tmockRepo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"
\t"github.com/rs/zerolog"
\t"github.com/stretchr/testify/assert"
\t"github.com/stretchr/testify/mock"
)'

# Fix credential_fallback_service_test.go
echo "Fixing credential_fallback_service_test.go"
# First, remove everything from line 1 up to first function
sed -i '' '2,25d' internal/domain/service/credential_fallback_service_test.go
# Then add proper imports after package declaration
sed -i '' '/^package service/ a\
'"$IMPORT_BLOCK" internal/domain/service/credential_fallback_service_test.go
# Fix mock references
sed -i '' 's/mockRepo\.mockRepo\.MockAPICredentialRepository/mockRepo\.MockAPICredentialRepository/g' internal/domain/service/credential_fallback_service_test.go

# Fix credential_manager_test.go
echo "Fixing credential_manager_test.go"
# First, remove everything from line 1 up to first function
sed -i '' '2,25d' internal/domain/service/credential_manager_test.go
# Then add proper imports after package declaration
sed -i '' '/^package service/ a\
'"$IMPORT_BLOCK" internal/domain/service/credential_manager_test.go
# Fix mock references
sed -i '' 's/mockRepo\.mockRepo\.MockAPICredentialRepository/mockRepo\.MockAPICredentialRepository/g' internal/domain/service/credential_manager_test.go

# Fix credential_encryption_service_test.go
echo "Fixing credential_encryption_service_test.go"
# First, remove problematic lines
sed -i '' '2,3d' internal/domain/service/credential_encryption_service_test.go
# Fix mock references
sed -i '' 's/mockRepo\.mockRepo\.MockAPICredentialRepository/mockRepo\.MockAPICredentialRepository/g' internal/domain/service/credential_encryption_service_test.go

# Fix credential_cache_service_test.go
echo "Fixing credential_cache_service_test.go"
# Fix mock references
sed -i '' 's/repoMocks/mockRepo/g' internal/domain/service/credential_cache_service_test.go

echo "All credential service test files fixed!" 
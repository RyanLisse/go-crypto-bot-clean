# Mocks Directory

This directory contains all mock implementations used for testing.

## Directory Structure

- `domain/port`: Mocks for domain ports (repositories, clients, etc.)
- `domain/service`: Mocks for domain services
- `usecase`: Mocks for use cases
- `infrastructure`: Mocks for infrastructure components (DB, cache, etc.)
- `adapter`: Mocks for adapters (HTTP, persistence, etc.)

## Best Practices

1. Use a tool like mockery to generate mocks
2. Keep mocks in sync with their interfaces
3. Don't modify generated mocks directly
4. Include a generation command in a comment

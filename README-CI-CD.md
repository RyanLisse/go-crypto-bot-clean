# CI/CD and Testing Setup

This document outlines the CI/CD pipeline and testing setup for the Go Crypto Bot project.

## CI/CD Pipeline

### GitHub Actions

The project uses GitHub Actions for continuous integration and deployment:

- **Backend Workflow** (`.github/workflows/go.yml`):
  - Triggered on push to main and pull requests
  - Jobs:
    - **Build**: Compiles Go code and runs tests
    - **Lint**: Uses golangci-lint for code quality
    - **Security**: Uses gosec for security scanning
    - **Release**: Creates releases when tags are pushed

- **Frontend Deployment**:
  - Configured with Netlify (`frontend/netlify.toml`)

## Pre-commit Hooks

We use Husky to enforce quality checks before commits:

- **Pre-commit Hook** (`.husky/pre-commit`):
  - Runs frontend linting and tests
  - Runs backend tests
  - Prevents commits if tests fail

- **Commit Message Hook** (`.husky/commit-msg`):
  - Enforces conventional commit format
  - Format: `<type>[optional scope]: <description>`
  - Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert

## Testing Setup

### Backend Testing

- Standard Go testing framework
- Run tests with: `cd backend && go test ./...`

### Frontend Testing

- **Unit Tests**:
  - Vitest for unit testing
  - Run with: `cd frontend && bun run test`
  - Watch mode: `bun run test:watch`
  - Coverage: `bun run test:coverage`

- **End-to-End Tests**:
  - Playwright for e2e testing
  - Configuration in `frontend/playwright.config.ts`

## TDD Workflow

1. Write a failing test
2. Run the test to verify it fails
3. Implement the code to make the test pass
4. Run the test to verify it passes
5. Refactor if needed
6. Commit with conventional commit format

## Running All Tests

Use the provided script to run all tests:

```bash
./scripts/run-tests.sh
```

## Setting Up Pre-commit Hooks

If you've just cloned the repository, run:

```bash
npm install
npm run prepare
```

This will install Husky and set up the Git hooks.

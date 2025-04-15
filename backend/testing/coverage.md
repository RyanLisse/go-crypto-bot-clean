# Test Coverage Guide

## Running Tests with Coverage

To check test coverage across all packages, run:

```sh
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

- This will output coverage for each package and function, and a total coverage percentage.
- For an HTML coverage report:

```sh
go tool cover -html=coverage.out
```

## Coverage Goals

- **Unit tests**: Should mock all external dependencies (DB, APIs, etc.) and cover business logic.
- **Integration tests**: Should test real interactions between components (e.g., DB, external API with test keys).
- Maintain a clear separation: unit tests in regular `_test.go` files, integration tests in `*_integration_test.go` or a `/integration/` subdir.

## Improving Coverage

- Use the coverage report to identify untested functions.
- Add tests for uncovered logic, especially critical business paths.
- Use mocks from `internal/mocks/` or package-local `mocks/` for unit tests.
- For new features, require coverage for all new code.

## Example Workflow

1. Write/modify code.
2. Write or update tests.
3. Run coverage commands above.
4. Review uncovered lines/functions.
5. Add tests as needed.
6. Commit only when coverage is satisfactory.

---

_Last updated: 2025-04-15_

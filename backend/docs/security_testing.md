# Security Testing Guide

This guide provides instructions for testing the security features implemented in the cryptocurrency trading bot.

## Overview

Security testing is an essential part of the development process. It helps identify vulnerabilities and ensures that security features are working as expected. This guide covers:

1. **Unit Testing**: Testing individual security components
2. **Integration Testing**: Testing how security components work together
3. **Penetration Testing**: Simulating attacks to identify vulnerabilities
4. **Security Scanning**: Using automated tools to identify vulnerabilities
5. **Manual Testing**: Manually testing security features

## Prerequisites

Before running security tests, ensure you have the following:

- Go 1.16 or higher
- SQLite
- Clerk account
- Test environment

## Unit Testing

### Rate Limiting Tests

The rate limiting middleware can be tested using the following command:

```bash
go test -v ./internal/adapter/http/middleware -run TestAdvancedRateLimiter
```

This test verifies that:
- Basic rate limiting works
- IP-based rate limiting works
- User-based rate limiting works
- Endpoint-specific rate limiting works
- Excluded paths are not rate limited

### CSRF Protection Tests

The CSRF protection middleware can be tested using the following command:

```bash
go test -v ./internal/adapter/http/middleware -run TestCSRFMiddleware
```

This test verifies that:
- CSRF tokens are generated for GET requests
- CSRF tokens are validated for non-GET requests
- Invalid CSRF tokens are rejected
- Excluded paths are not CSRF protected
- Excluded methods are not CSRF protected

### Secure Headers Tests

The secure headers middleware can be tested using the following command:

```bash
go test -v ./internal/adapter/http/middleware -run TestSecureHeadersMiddleware
```

This test verifies that:
- Secure headers are set correctly
- Excluded paths do not have secure headers
- Custom headers are set correctly
- Server and X-Powered-By headers are removed

### Authentication Tests

The authentication middleware can be tested using the following command:

```bash
go test -v ./internal/adapter/http/middleware -run TestClerkMiddleware
```

This test verifies that:
- Valid tokens are accepted
- Invalid tokens are rejected
- User ID is set in the context
- User roles are set in the context

### Encryption Tests

The encryption service can be tested using the following command:

```bash
go test -v ./internal/domain/service -run TestEncryptionService
```

This test verifies that:
- Data can be encrypted and decrypted
- Invalid data cannot be decrypted
- Key rotation works correctly

## Integration Testing

### API Security Tests

The API security features can be tested using the following command:

```bash
go test -v ./internal/adapter/http/controller -run TestAPISecurityIntegration
```

This test verifies that:
- Rate limiting is applied to API endpoints
- CSRF protection is applied to API endpoints
- Secure headers are set in API responses
- Authentication is required for protected endpoints
- Role-based access control works correctly

### End-to-End Security Tests

End-to-end security tests can be run using the following command:

```bash
go test -v ./test/e2e -run TestSecurityE2E
```

This test verifies that:
- The entire security stack works correctly
- Security features are applied in the correct order
- Security features do not interfere with each other

## Penetration Testing

### OWASP ZAP

[OWASP ZAP](https://www.zaproxy.org/) is a free security tool that can be used to find vulnerabilities in web applications.

To run a basic scan:

1. Start the application:
   ```bash
   go run cmd/server/main.go
   ```

2. Run ZAP and configure it to scan your application:
   ```bash
   zap.sh -cmd -quickurl http://localhost:8080
   ```

3. Review the results and fix any vulnerabilities.

### Burp Suite

[Burp Suite](https://portswigger.net/burp) is a popular security testing tool that can be used to find vulnerabilities in web applications.

To run a basic scan:

1. Start the application:
   ```bash
   go run cmd/server/main.go
   ```

2. Configure Burp Suite to proxy your requests.

3. Use the Burp Scanner to scan your application.

4. Review the results and fix any vulnerabilities.

## Security Scanning

### GoSec

[GoSec](https://github.com/securego/gosec) is a security scanner for Go code.

To run a basic scan:

```bash
gosec ./...
```

This will scan your code for common security issues and report any findings.

### Dependency Scanning

Use [Go Dependency Scanner](https://github.com/sonatype-nexus-community/nancy) to scan your dependencies for known vulnerabilities:

```bash
go list -json -m all | nancy sleuth
```

This will scan your dependencies for known vulnerabilities and report any findings.

## Manual Testing

### Rate Limiting

To manually test rate limiting:

1. Start the application:
   ```bash
   go run cmd/server/main.go
   ```

2. Send multiple requests to an endpoint:
   ```bash
   for i in {1..100}; do curl -i http://localhost:8080/api/v1/users; done
   ```

3. Verify that rate limiting is applied after the configured number of requests.

### CSRF Protection

To manually test CSRF protection:

1. Start the application:
   ```bash
   go run cmd/server/main.go
   ```

2. Send a GET request to get a CSRF token:
   ```bash
   curl -i -c cookies.txt http://localhost:8080/api/v1/users
   ```

3. Send a POST request without a CSRF token:
   ```bash
   curl -i -X POST -b cookies.txt http://localhost:8080/api/v1/users
   ```

4. Verify that the request is rejected with a 403 Forbidden response.

5. Send a POST request with a CSRF token:
   ```bash
   TOKEN=$(grep csrf_token cookies.txt | awk '{print $7}')
   curl -i -X POST -b cookies.txt -H "X-CSRF-Token: $TOKEN" http://localhost:8080/api/v1/users
   ```

6. Verify that the request is accepted.

### Secure Headers

To manually test secure headers:

1. Start the application:
   ```bash
   go run cmd/server/main.go
   ```

2. Send a request to an endpoint:
   ```bash
   curl -i http://localhost:8080/api/v1/users
   ```

3. Verify that the response includes the configured secure headers.

### Authentication

To manually test authentication:

1. Start the application:
   ```bash
   go run cmd/server/main.go
   ```

2. Send a request without authentication:
   ```bash
   curl -i http://localhost:8080/api/v1/users/me
   ```

3. Verify that the request is rejected with a 401 Unauthorized response.

4. Send a request with authentication:
   ```bash
   curl -i -H "Authorization: Bearer <clerk_session_token>" http://localhost:8080/api/v1/users/me
   ```

5. Verify that the request is accepted.

### Role-Based Access Control

To manually test role-based access control:

1. Start the application:
   ```bash
   go run cmd/server/main.go
   ```

2. Send a request with a user that does not have the required role:
   ```bash
   curl -i -H "Authorization: Bearer <clerk_session_token>" http://localhost:8080/api/v1/admin
   ```

3. Verify that the request is rejected with a 403 Forbidden response.

4. Send a request with a user that has the required role:
   ```bash
   curl -i -H "Authorization: Bearer <admin_clerk_session_token>" http://localhost:8080/api/v1/admin
   ```

5. Verify that the request is accepted.

## Continuous Integration

Security tests should be integrated into your CI/CD pipeline to ensure that security issues are caught early.

### GitHub Actions

Here's an example GitHub Actions workflow that runs security tests:

```yaml
name: Security Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  security:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.16
    - name: Install GoSec
      run: go install github.com/securego/gosec/v2/cmd/gosec@latest
    - name: Run GoSec
      run: gosec ./...
    - name: Run Unit Tests
      run: go test -v ./internal/adapter/http/middleware -run "TestAdvancedRateLimiter|TestCSRFMiddleware|TestSecureHeadersMiddleware"
    - name: Run Integration Tests
      run: go test -v ./internal/adapter/http/controller -run TestAPISecurityIntegration
    - name: Scan Dependencies
      run: go list -json -m all | nancy sleuth
```

## Reporting Security Issues

If you discover a security issue, please report it by sending an email to security@example.com. Do not disclose security issues publicly until they have been handled by the security team.

## References

- [OWASP Top Ten](https://owasp.org/www-project-top-ten/)
- [OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)
- [OWASP ZAP](https://www.zaproxy.org/)
- [Burp Suite](https://portswigger.net/burp)
- [GoSec](https://github.com/securego/gosec)
- [Nancy](https://github.com/sonatype-nexus-community/nancy)

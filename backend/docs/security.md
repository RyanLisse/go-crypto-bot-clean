# Security Documentation

This document provides comprehensive documentation for the security features implemented in the cryptocurrency trading bot.

## Overview

The application implements multiple layers of security to protect against common web vulnerabilities and attacks. These include:

1. **Rate Limiting**: Prevents abuse and DoS attacks by limiting the number of requests.
2. **CSRF Protection**: Prevents Cross-Site Request Forgery attacks.
3. **Secure HTTP Headers**: Protects against various attacks by setting appropriate HTTP headers.
4. **Authentication**: Secures user access using Clerk authentication.
5. **Encryption**: Protects sensitive data using strong encryption.

## Rate Limiting

### Overview

Rate limiting is implemented to prevent abuse and DoS attacks by limiting the number of requests that can be made in a given time period. The rate limiter supports:

- IP-based rate limiting
- User-based rate limiting
- Endpoint-specific rate limiting
- Configurable limits
- Redis support for distributed rate limiting

### Configuration

Rate limiting can be configured in the `config.yaml` file:

```yaml
rate_limit:
  enabled: true
  default_limit: 60   # 1 request per second
  default_burst: 10   # Allow bursts of 10 requests
  ip_limit: 300       # 5 requests per second per IP
  ip_burst: 20        # Allow bursts of 20 requests per IP
  user_limit: 600     # 10 requests per second per user
  user_burst: 30      # Allow bursts of 30 requests per user
  auth_user_limit: 1200  # 20 requests per second for authenticated users
  auth_user_burst: 60    # Allow bursts of 60 requests for authenticated users
  cleanup_interval: 5m
  block_duration: 15m
  trusted_proxies:
    - "127.0.0.1"
    - "::1"
  excluded_paths:
    - "/health"
    - "/metrics"
    - "/favicon.ico"
  redis_enabled: false
  redis_key_prefix: "ratelimit:"
```

### Implementation

The rate limiter is implemented in the `AdvancedRateLimiter` middleware, which uses the token bucket algorithm to limit requests. When a request exceeds the rate limit, a 429 Too Many Requests response is returned.

### Headers

The rate limiter sets the following headers in the response:

- `X-RateLimit-Limit`: The maximum number of requests allowed in the current time window.
- `X-RateLimit-Remaining`: The number of requests remaining in the current time window.
- `X-RateLimit-Reset`: The time when the current rate limit window resets.
- `Retry-After`: The number of seconds to wait before making another request.

## CSRF Protection

### Overview

Cross-Site Request Forgery (CSRF) protection is implemented to prevent attackers from tricking users into performing unwanted actions. The CSRF protection middleware:

- Generates a unique token for each user session
- Validates the token for non-GET requests
- Supports both cookie and header-based tokens
- Allows excluding specific paths and methods

### Configuration

CSRF protection can be configured in the `config.yaml` file:

```yaml
csrf:
  enabled: true
  secret: "${CSRF_SECRET}"
  token_length: 32
  cookie_name: "csrf_token"
  cookie_path: "/"
  cookie_max_age: 24h
  cookie_secure: true
  cookie_http_only: true
  cookie_same_site: "Lax"
  header_name: "X-CSRF-Token"
  form_field_name: "csrf_token"
  excluded_paths:
    - "/health"
    - "/metrics"
    - "/favicon.ico"
    - "/api/v1/auth/verify"
  excluded_methods:
    - "GET"
    - "HEAD"
    - "OPTIONS"
    - "TRACE"
  failure_status_code: 403
```

### Implementation

The CSRF protection is implemented in the `CSRFMiddleware` middleware, which:

1. Generates a CSRF token for GET requests
2. Sets the token in a cookie
3. Validates the token for non-GET requests
4. Returns a 403 Forbidden response if the token is invalid

### Usage

To use CSRF protection in your frontend:

1. For GET requests, the CSRF token is automatically set in a cookie.
2. For non-GET requests, include the CSRF token in the request:
   - As a header: `X-CSRF-Token: <token>`
   - As a form field: `csrf_token=<token>`
   - As a query parameter: `?csrf_token=<token>`

## Secure HTTP Headers

### Overview

Secure HTTP headers are implemented to protect against various attacks, including:

- Cross-Site Scripting (XSS)
- Clickjacking
- MIME type sniffing
- Information disclosure
- Cross-Origin Resource Sharing (CORS) attacks

### Configuration

Secure headers can be configured in the `config.yaml` file:

```yaml
secure_headers:
  enabled: true
  content_security_policy: "default-src 'self'; script-src 'self'; object-src 'none'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; font-src 'self'; frame-src 'none'; connect-src 'self'"
  x_content_type_options: "nosniff"
  x_frame_options: "DENY"
  x_xss_protection: "1; mode=block"
  referrer_policy: "strict-origin-when-cross-origin"
  strict_transport_security: "max-age=31536000; includeSubDomains"
  permissions_policy: "camera=(), microphone=(), geolocation=(), interest-cohort=()"
  cross_origin_embedder_policy: "require-corp"
  cross_origin_opener_policy: "same-origin"
  cross_origin_resource_policy: "same-origin"
  cache_control: "no-store, max-age=0"
  excluded_paths:
    - "/health"
    - "/metrics"
    - "/favicon.ico"
  remove_server_header: true
  remove_powered_by_header: true
  content_security_policy_report_only: false
  content_security_policy_report_uri: ""
```

### Implementation

The secure headers are implemented in the `SecureHeadersMiddleware` middleware, which sets the appropriate headers in the response.

### Headers

The following headers are set:

- **Content-Security-Policy**: Restricts the sources from which resources can be loaded.
- **X-Content-Type-Options**: Prevents MIME type sniffing.
- **X-Frame-Options**: Prevents clickjacking by restricting framing.
- **X-XSS-Protection**: Enables the browser's XSS filter.
- **Referrer-Policy**: Controls the information sent in the Referer header.
- **Strict-Transport-Security**: Enforces HTTPS.
- **Permissions-Policy**: Controls access to browser features.
- **Cross-Origin-Embedder-Policy**: Controls which resources can be embedded.
- **Cross-Origin-Opener-Policy**: Controls window.opener behavior.
- **Cross-Origin-Resource-Policy**: Controls which sites can load resources.
- **Cache-Control**: Controls caching behavior.

## Authentication

### Overview

Authentication is implemented using Clerk, a secure authentication service. The authentication system:

- Validates user tokens
- Manages user sessions
- Supports role-based access control
- Integrates with the database

### Configuration

Authentication can be configured in the `config.yaml` file:

```yaml
auth:
  clerk_secret_key: "${CLERK_SECRET_KEY}"
  use_enhanced_auth: true
```

### Implementation

Authentication is implemented in the `ClerkMiddleware` and `EnhancedClerkMiddleware` middlewares, which:

1. Validate the user's token
2. Set the user ID in the request context
3. Retrieve user roles from Clerk
4. Support role-based access control

### Usage

To authenticate a request:

1. Include the Clerk session token in the Authorization header:
   ```
   Authorization: Bearer <clerk_session_token>
   ```

2. The middleware will validate the token and set the user ID in the request context.

3. To require authentication for a route:
   ```go
   router.Group(func(r chi.Router) {
       r.Use(authMiddleware.RequireAuthentication)
       // Protected routes
   })
   ```

4. To require a specific role for a route:
   ```go
   router.Group(func(r chi.Router) {
       r.Use(authMiddleware.RequireRole("admin"))
       // Admin-only routes
   })
   ```

## Encryption

### Overview

Encryption is implemented to protect sensitive data, such as API keys and secrets. The encryption system:

- Uses AES-256-GCM for encryption
- Supports key rotation
- Securely stores encryption keys
- Provides a simple API for encrypting and decrypting data

### Configuration

Encryption can be configured in the `config.yaml` file:

```yaml
encryption:
  enabled: true
  key_file: "${ENCRYPTION_KEY_FILE}"
  key_rotation_interval: 720h  # 30 days
```

### Implementation

Encryption is implemented in the `EncryptionService`, which provides methods for encrypting and decrypting data.

### Usage

To encrypt and decrypt data:

```go
// Encrypt data
encryptedData, err := encryptionService.Encrypt(data)
if err != nil {
    // Handle error
}

// Decrypt data
decryptedData, err := encryptionService.Decrypt(encryptedData)
if err != nil {
    // Handle error
}
```

## Best Practices

### General Security Best Practices

1. **Keep Dependencies Updated**: Regularly update dependencies to patch security vulnerabilities.
2. **Use HTTPS**: Always use HTTPS in production to encrypt data in transit.
3. **Implement Proper Logging**: Log security events but avoid logging sensitive information.
4. **Use Secure Passwords**: Enforce strong password policies.
5. **Implement MFA**: Use multi-factor authentication where possible.
6. **Regular Security Audits**: Conduct regular security audits and penetration testing.
7. **Follow the Principle of Least Privilege**: Only grant the minimum necessary permissions.

### API Security Best Practices

1. **Validate All Input**: Validate and sanitize all user input to prevent injection attacks.
2. **Use Parameterized Queries**: Use parameterized queries to prevent SQL injection.
3. **Implement API Rate Limiting**: Limit the number of requests to prevent abuse.
4. **Use API Keys**: Use API keys to authenticate API requests.
5. **Implement Proper Error Handling**: Return generic error messages to users to avoid information disclosure.
6. **Use Content Security Policy**: Implement a strict Content Security Policy.
7. **Implement CORS**: Configure CORS to restrict which domains can access your API.

## Security Headers Reference

### Content-Security-Policy

The Content-Security-Policy header restricts the sources from which resources can be loaded. It helps prevent Cross-Site Scripting (XSS) attacks.

Example:
```
Content-Security-Policy: default-src 'self'; script-src 'self'; object-src 'none'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; font-src 'self'; frame-src 'none'; connect-src 'self'
```

### X-Content-Type-Options

The X-Content-Type-Options header prevents MIME type sniffing, which can lead to security vulnerabilities.

Example:
```
X-Content-Type-Options: nosniff
```

### X-Frame-Options

The X-Frame-Options header prevents clickjacking attacks by restricting how the page can be framed.

Example:
```
X-Frame-Options: DENY
```

### X-XSS-Protection

The X-XSS-Protection header enables the browser's built-in XSS filter.

Example:
```
X-XSS-Protection: 1; mode=block
```

### Referrer-Policy

The Referrer-Policy header controls how much information is sent in the Referer header.

Example:
```
Referrer-Policy: strict-origin-when-cross-origin
```

### Strict-Transport-Security

The Strict-Transport-Security header enforces HTTPS by telling browsers to always use HTTPS for the domain.

Example:
```
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

### Permissions-Policy

The Permissions-Policy header controls which browser features and APIs can be used.

Example:
```
Permissions-Policy: camera=(), microphone=(), geolocation=(), interest-cohort=()
```

### Cross-Origin-Embedder-Policy

The Cross-Origin-Embedder-Policy header controls which resources can be embedded in the page.

Example:
```
Cross-Origin-Embedder-Policy: require-corp
```

### Cross-Origin-Opener-Policy

The Cross-Origin-Opener-Policy header controls the behavior of window.opener.

Example:
```
Cross-Origin-Opener-Policy: same-origin
```

### Cross-Origin-Resource-Policy

The Cross-Origin-Resource-Policy header controls which sites can load the resource.

Example:
```
Cross-Origin-Resource-Policy: same-origin
```

### Cache-Control

The Cache-Control header controls how the response is cached.

Example:
```
Cache-Control: no-store, max-age=0
```

## References

- [OWASP Top Ten](https://owasp.org/www-project-top-ten/)
- [OWASP Secure Headers Project](https://owasp.org/www-project-secure-headers/)
- [Content Security Policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)
- [HTTP Security Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#security)
- [Clerk Documentation](https://clerk.dev/docs)

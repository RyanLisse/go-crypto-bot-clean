# Middleware Standardization Strategy

This document outlines our approach to standardizing and consolidating middleware across the application, focusing on eliminating redundancy and ensuring consistent behavior.

## Core Middleware Pattern

The application will use a consolidated middleware approach using the standard http.Handler interface:

```go
func (m *SomeMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Pre-processing logic
            
            // Call the next handler
            next.ServeHTTP(w, r)
            
            // Post-processing logic (optional)
        })
    }
}
```

## Standardized Middleware Components

### 1. Authentication Middleware

We will use a single `ConsolidatedAuthMiddleware` for all authentication needs, replacing multiple implementations like `SimpleAuthMiddleware`, `ClerkMiddleware`, etc.

```go
// internal/adapter/http/middleware/consolidated_auth_middleware.go

type ConsolidatedAuthMiddleware struct {
    // Dependencies
    authService      port.AuthService
    userService      port.UserService
    logger           *zap.Logger
    enableDummyAuth  bool
    dummyUserID      string
}

// NewConsolidatedAuthMiddleware creates a new authentication middleware
func NewConsolidatedAuthMiddleware(
    authService port.AuthService,
    userService port.UserService,
    logger *zap.Logger,
    enableDummyAuth bool,
    dummyUserID string,
) *ConsolidatedAuthMiddleware {
    return &ConsolidatedAuthMiddleware{
        authService:     authService,
        userService:     userService,
        logger:          logger,
        enableDummyAuth: enableDummyAuth,
        dummyUserID:     dummyUserID,
    }
}

// Middleware creates a middleware function
func (m *ConsolidatedAuthMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from Authorization header
            authHeader := r.Header.Get("Authorization")
            
            // Test mode short circuit if enabled
            if m.enableDummyAuth {
                // Only enable in non-production environments
                if config.GetEnvironment() != "production" {
                    user, err := m.userService.GetUserByID(r.Context(), m.dummyUserID)
                    if err == nil {
                        // Set user in context
                        ctx := context.WithValue(r.Context(), UserContextKey, user)
                        next.ServeHTTP(w, r.WithContext(ctx))
                        return
                    }
                    // Fall through to normal auth if dummy user not found
                }
            }
            
            // Normal authentication flow
            if authHeader == "" {
                // No auth, continue without user in context
                next.ServeHTTP(w, r)
                return
            }
            
            // Extract token
            token := strings.TrimPrefix(authHeader, "Bearer ")
            
            // Validate token
            userID, err := m.authService.ValidateToken(r.Context(), token)
            if err != nil {
                // Token invalid, but still continue without user context
                m.logger.Debug("Invalid token", zap.Error(err))
                next.ServeHTTP(w, r)
                return
            }
            
            // Get user from database
            user, err := m.userService.GetUserByID(r.Context(), userID)
            if err != nil {
                m.logger.Debug("User not found for valid token", zap.String("userID", userID), zap.Error(err))
                next.ServeHTTP(w, r)
                return
            }
            
            // Set user in context
            ctx := context.WithValue(r.Context(), UserContextKey, user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// RequireAuthentication is a middleware that requires authentication
func (m *ConsolidatedAuthMiddleware) RequireAuthentication(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := GetUserFromContext(r.Context())
        if user == nil {
            errorResponse := domain.NewErrorResponse("Unauthorized", "Authentication required", http.StatusUnauthorized)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(errorResponse)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// RequireRole is a middleware that requires a specific role
func (m *ConsolidatedAuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := GetUserFromContext(r.Context())
            if user == nil {
                errorResponse := domain.NewErrorResponse("Unauthorized", "Authentication required", http.StatusUnauthorized)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusUnauthorized)
                json.NewEncoder(w).Encode(errorResponse)
                return
            }
            
            // Check if user has the required role
            hasRole := false
            for _, userRole := range user.Roles {
                if userRole == role {
                    hasRole = true
                    break
                }
            }
            
            if !hasRole {
                errorResponse := domain.NewErrorResponse("Forbidden", "Insufficient permissions", http.StatusForbidden)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusForbidden)
                json.NewEncoder(w).Encode(errorResponse)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// GetUserFromContext extracts the user from the context
func GetUserFromContext(ctx context.Context) *model.User {
    user, ok := ctx.Value(UserContextKey).(*model.User)
    if !ok {
        return nil
    }
    return user
}
```

### 2. Logging Middleware

A single logging middleware will be used for consistent request/response logging:

```go
// internal/adapter/http/middleware/logging_middleware.go

type LoggingMiddleware struct {
    logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
    return &LoggingMiddleware{logger: logger}
}

func (m *LoggingMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Create a response writer that captures status code
            ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
            
            // Get or generate request ID
            requestID := r.Header.Get("X-Request-ID")
            if requestID == "" {
                requestID = uuid.New().String()
                r.Header.Set("X-Request-ID", requestID)
            }
            
            // Add request ID to context
            ctx := context.WithValue(r.Context(), RequestIDContextKey, requestID)
            r = r.WithContext(ctx)
            
            // Pre-request logging
            m.logger.Debug("Request started",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("request_id", requestID),
                zap.String("remote_addr", r.RemoteAddr),
                zap.String("user_agent", r.UserAgent()),
            )
            
            // Process request
            next.ServeHTTP(ww, r)
            
            // Post-request logging
            duration := time.Since(start)
            
            // Determine log level based on status code
            logFunc := m.logger.Info
            if ww.Status() >= 500 {
                logFunc = m.logger.Error
            } else if ww.Status() >= 400 {
                logFunc = m.logger.Warn
            }
            
            logFunc("Request completed",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("request_id", requestID),
                zap.Int("status", ww.Status()),
                zap.Duration("duration", duration),
                zap.Int("bytes", ww.BytesWritten()),
            )
        })
    }
}
```

### 3. Error Handling Middleware

A standardized approach to error handling with a single middleware:

```go
// internal/adapter/http/middleware/error_middleware.go

type ErrorMiddleware struct {
    logger *zap.Logger
}

func NewErrorMiddleware(logger *zap.Logger) *ErrorMiddleware {
    return &ErrorMiddleware{logger: logger}
}

func (m *ErrorMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Create a panic recovery wrapper
            defer func() {
                if err := recover(); err != nil {
                    // Log the stack trace
                    stackTrace := debug.Stack()
                    m.logger.Error("Panic recovered in HTTP handler",
                        zap.Any("error", err),
                        zap.String("stack", string(stackTrace)),
                        zap.String("path", r.URL.Path),
                        zap.String("method", r.Method),
                    )
                    
                    // Return 500 error
                    errorResponse := domain.NewErrorResponse(
                        "InternalServerError",
                        "An unexpected error occurred",
                        http.StatusInternalServerError,
                    )
                    w.Header().Set("Content-Type", "application/json")
                    w.WriteHeader(http.StatusInternalServerError)
                    json.NewEncoder(w).Encode(errorResponse)
                }
            }()
            
            // Continue to the next middleware/handler
            next.ServeHTTP(w, r)
        })
    }
}
```

### 4. CORS Middleware

A single middleware for handling CORS:

```go
// internal/adapter/http/middleware/cors_middleware.go

type CORSMiddleware struct {
    allowedOrigins []string
    allowedMethods []string
    allowedHeaders []string
    maxAge         int
}

func NewCORSMiddleware(
    allowedOrigins []string,
    allowedMethods []string,
    allowedHeaders []string,
    maxAge int,
) *CORSMiddleware {
    if len(allowedMethods) == 0 {
        allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    }
    
    if len(allowedHeaders) == 0 {
        allowedHeaders = []string{
            "Accept", "Content-Type", "Content-Length", "Authorization",
            "X-CSRF-Token", "X-Requested-With", "X-Request-ID",
        }
    }
    
    return &CORSMiddleware{
        allowedOrigins: allowedOrigins,
        allowedMethods: allowedMethods,
        allowedHeaders: allowedHeaders,
        maxAge:         maxAge,
    }
}

func (m *CORSMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Check if the origin is allowed
            allowed := false
            for _, allowedOrigin := range m.allowedOrigins {
                if allowedOrigin == "*" || allowedOrigin == origin {
                    allowed = true
                    break
                }
            }
            
            if allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.allowedMethods, ", "))
                w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.allowedHeaders, ", "))
                w.Header().Set("Access-Control-Max-Age", strconv.Itoa(m.maxAge))
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
            
            // Handle preflight requests
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

### 5. Rate Limiting Middleware

A consolidated rate limiter that supports both simple and advanced configurations:

```go
// internal/adapter/http/middleware/rate_limiter.go

type RateLimiterMiddleware struct {
    store     port.RateLimitStore
    logger    *zap.Logger
    ipLimits  []*IPLimit
    userLimits []*UserLimit
    routeLimits map[string][]*RouteLimit
}

func NewRateLimiterMiddleware(
    store port.RateLimitStore,
    logger *zap.Logger,
) *RateLimiterMiddleware {
    return &RateLimiterMiddleware{
        store:       store,
        logger:      logger,
        ipLimits:    []*IPLimit{},
        userLimits:  []*UserLimit{},
        routeLimits: make(map[string][]*RouteLimit),
    }
}

// AddIPLimit adds an IP-based rate limit
func (m *RateLimiterMiddleware) AddIPLimit(requests int, duration time.Duration) *RateLimiterMiddleware {
    m.ipLimits = append(m.ipLimits, &IPLimit{
        Requests: requests,
        Duration: duration,
    })
    return m
}

// AddUserLimit adds a user-based rate limit
func (m *RateLimiterMiddleware) AddUserLimit(requests int, duration time.Duration) *RateLimiterMiddleware {
    m.userLimits = append(m.userLimits, &UserLimit{
        Requests: requests,
        Duration: duration,
    })
    return m
}

// AddRouteLimit adds a route-specific rate limit
func (m *RateLimiterMiddleware) AddRouteLimit(
    method string,
    pathPattern string,
    requests int,
    duration time.Duration,
) *RateLimiterMiddleware {
    key := method + ":" + pathPattern
    m.routeLimits[key] = append(m.routeLimits[key], &RouteLimit{
        Method:   method,
        Path:     pathPattern,
        Requests: requests,
        Duration: duration,
    })
    return m
}

// Middleware creates the rate limiting middleware
func (m *RateLimiterMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get client IP
            clientIP := GetClientIP(r)
            
            // Check IP limits
            for _, limit := range m.ipLimits {
                key := fmt.Sprintf("ip:%s:%d:%d", clientIP, limit.Requests, limit.Duration.Seconds())
                
                // Check if IP is rate limited
                remaining, reset, err := m.store.Check(r.Context(), key, limit.Requests, limit.Duration)
                if err != nil {
                    m.logger.Error("Failed to check rate limit",
                        zap.Error(err),
                        zap.String("ip", clientIP),
                    )
                    // Continue despite error to avoid blocking legitimate traffic
                } else if remaining < 0 {
                    m.logger.Warn("IP rate limit exceeded",
                        zap.String("ip", clientIP),
                        zap.Int("limit", limit.Requests),
                        zap.Duration("duration", limit.Duration),
                    )
                    
                    // Set rate limit headers
                    w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit.Requests))
                    w.Header().Set("X-RateLimit-Remaining", "0")
                    w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))
                    w.Header().Set("Retry-After", strconv.FormatInt(int64(reset.Sub(time.Now()).Seconds()), 10))
                    
                    // Return rate limit error
                    errorResponse := domain.NewErrorResponse(
                        "RateLimitExceeded",
                        "Rate limit exceeded. Try again later.",
                        http.StatusTooManyRequests,
                    )
                    w.Header().Set("Content-Type", "application/json")
                    w.WriteHeader(http.StatusTooManyRequests)
                    json.NewEncoder(w).Encode(errorResponse)
                    return
                } else {
                    // Set rate limit headers
                    w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit.Requests))
                    w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
                    w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))
                }
            }
            
            // Check user limits if authenticated
            user := GetUserFromContext(r.Context())
            if user != nil {
                for _, limit := range m.userLimits {
                    key := fmt.Sprintf("user:%s:%d:%d", user.ID, limit.Requests, limit.Duration.Seconds())
                    
                    // Check if user is rate limited
                    remaining, reset, err := m.store.Check(r.Context(), key, limit.Requests, limit.Duration)
                    if err != nil {
                        m.logger.Error("Failed to check user rate limit",
                            zap.Error(err),
                            zap.String("user_id", user.ID),
                        )
                        // Continue despite error to avoid blocking legitimate traffic
                    } else if remaining < 0 {
                        m.logger.Warn("User rate limit exceeded",
                            zap.String("user_id", user.ID),
                            zap.Int("limit", limit.Requests),
                            zap.Duration("duration", limit.Duration),
                        )
                        
                        // Set rate limit headers
                        w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit.Requests))
                        w.Header().Set("X-RateLimit-Remaining", "0")
                        w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))
                        w.Header().Set("Retry-After", strconv.FormatInt(int64(reset.Sub(time.Now()).Seconds()), 10))
                        
                        // Return rate limit error
                        errorResponse := domain.NewErrorResponse(
                            "RateLimitExceeded", 
                            "Rate limit exceeded. Try again later.",
                            http.StatusTooManyRequests,
                        )
                        w.Header().Set("Content-Type", "application/json")
                        w.WriteHeader(http.StatusTooManyRequests)
                        json.NewEncoder(w).Encode(errorResponse)
                        return
                    } else {
                        // Set rate limit headers
                        w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit.Requests))
                        w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
                        w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))
                    }
                }
            }
            
            // Check route limits
            method := r.Method
            path := r.URL.Path
            
            for pattern, limits := range m.routeLimits {
                patternParts := strings.Split(pattern, ":")
                if len(patternParts) != 2 {
                    continue
                }
                
                methodPattern := patternParts[0]
                pathPattern := patternParts[1]
                
                // Check if method matches
                if methodPattern != "*" && methodPattern != method {
                    continue
                }
                
                // Check if path matches
                matched, _ := path.Match(pathPattern, path)
                if !matched {
                    continue
                }
                
                // Apply route-specific limits
                for _, limit := range limits {
                    // Build key based on whether user is authenticated
                    var key string
                    if user != nil {
                        key = fmt.Sprintf("route:user:%s:%s:%s:%d:%d", 
                            user.ID, method, pathPattern, limit.Requests, limit.Duration.Seconds())
                    } else {
                        key = fmt.Sprintf("route:ip:%s:%s:%s:%d:%d", 
                            clientIP, method, pathPattern, limit.Requests, limit.Duration.Seconds())
                    }
                    
                    // Check rate limit
                    remaining, reset, err := m.store.Check(r.Context(), key, limit.Requests, limit.Duration)
                    if err != nil {
                        m.logger.Error("Failed to check route rate limit",
                            zap.Error(err),
                            zap.String("method", method),
                            zap.String("path", path),
                        )
                        // Continue despite error to avoid blocking legitimate traffic
                    } else if remaining < 0 {
                        m.logger.Warn("Route rate limit exceeded",
                            zap.String("method", method),
                            zap.String("path", path),
                            zap.Int("limit", limit.Requests),
                            zap.Duration("duration", limit.Duration),
                        )
                        
                        // Set rate limit headers
                        w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit.Requests))
                        w.Header().Set("X-RateLimit-Remaining", "0")
                        w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))
                        w.Header().Set("Retry-After", strconv.FormatInt(int64(reset.Sub(time.Now()).Seconds()), 10))
                        
                        // Return rate limit error
                        errorResponse := domain.NewErrorResponse(
                            "RateLimitExceeded", 
                            "Rate limit exceeded for this endpoint. Try again later.",
                            http.StatusTooManyRequests,
                        )
                        w.Header().Set("Content-Type", "application/json")
                        w.WriteHeader(http.StatusTooManyRequests)
                        json.NewEncoder(w).Encode(errorResponse)
                        return
                    } else {
                        // Set rate limit headers
                        w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit.Requests))
                        w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
                        w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))
                    }
                }
            }
            
            // Continue to the next middleware/handler
            next.ServeHTTP(w, r)
        })
    }
}

// Helper function to get client IP
func GetClientIP(r *http.Request) string {
    // Check for X-Forwarded-For header
    xForwardedFor := r.Header.Get("X-Forwarded-For")
    if xForwardedFor != "" {
        // X-Forwarded-For can contain multiple IPs, take the first one
        ips := strings.Split(xForwardedFor, ",")
        if len(ips) > 0 {
            return strings.TrimSpace(ips[0])
        }
    }
    
    // Check for X-Real-IP header
    xRealIP := r.Header.Get("X-Real-IP")
    if xRealIP != "" {
        return xRealIP
    }
    
    // Fall back to RemoteAddr
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return r.RemoteAddr
    }
    return ip
}
```

## Middleware Factory

To streamline middleware creation and ensure consistent configuration, we'll use a factory pattern:

```go
// internal/factory/middleware_factory.go

type MiddlewareFactory struct {
    logger           *zap.Logger
    authService      port.AuthService
    userService      port.UserService
    rateLimitStore   port.RateLimitStore
    config           *config.Config
}

func NewMiddlewareFactory(
    logger *zap.Logger,
    authService port.AuthService,
    userService port.UserService,
    rateLimitStore port.RateLimitStore,
    config *config.Config,
) *MiddlewareFactory {
    return &MiddlewareFactory{
        logger:         logger,
        authService:    authService,
        userService:    userService,
        rateLimitStore: rateLimitStore,
        config:         config,
    }
}

// CreateAuthMiddleware creates a new authentication middleware
func (f *MiddlewareFactory) CreateAuthMiddleware() *middleware.ConsolidatedAuthMiddleware {
    // Test mode is only allowed in non-production environments
    enableDummyAuth := f.config.Auth.EnableDummyAuth && f.config.Environment != "production"
    
    return middleware.NewConsolidatedAuthMiddleware(
        f.authService,
        f.userService,
        f.logger,
        enableDummyAuth,
        f.config.Auth.DummyUserID,
    )
}

// CreateLoggingMiddleware creates a new logging middleware
func (f *MiddlewareFactory) CreateLoggingMiddleware() *middleware.LoggingMiddleware {
    return middleware.NewLoggingMiddleware(f.logger)
}

// CreateErrorMiddleware creates a new error handling middleware
func (f *MiddlewareFactory) CreateErrorMiddleware() *middleware.ErrorMiddleware {
    return middleware.NewErrorMiddleware(f.logger)
}

// CreateCORSMiddleware creates a new CORS middleware
func (f *MiddlewareFactory) CreateCORSMiddleware() *middleware.CORSMiddleware {
    return middleware.NewCORSMiddleware(
        f.config.CORS.AllowedOrigins,
        f.config.CORS.AllowedMethods,
        f.config.CORS.AllowedHeaders,
        f.config.CORS.MaxAge,
    )
}

// CreateRateLimiterMiddleware creates a new rate limiter middleware
func (f *MiddlewareFactory) CreateRateLimiterMiddleware() *middleware.RateLimiterMiddleware {
    rateLimiter := middleware.NewRateLimiterMiddleware(
        f.rateLimitStore,
        f.logger,
    )
    
    // Add default IP limits
    for _, limit := range f.config.RateLimit.IPLimits {
        rateLimiter.AddIPLimit(limit.Requests, limit.Duration)
    }
    
    // Add default user limits
    for _, limit := range f.config.RateLimit.UserLimits {
        rateLimiter.AddUserLimit(limit.Requests, limit.Duration)
    }
    
    // Add route-specific limits
    for _, limit := range f.config.RateLimit.RouteLimits {
        rateLimiter.AddRouteLimit(
            limit.Method,
            limit.Path,
            limit.Requests,
            limit.Duration,
        )
    }
    
    return rateLimiter
}
```

## Mock Middleware for Testing

For testing purposes, we'll implement a mock version of the authentication middleware:

```go
// internal/adapter/http/middleware/test_auth_middleware.go

type TestAuthMiddleware struct {
    userService port.UserService
    logger      *zap.Logger
    testUserID  string
}

func NewTestAuthMiddleware(
    userService port.UserService,
    logger *zap.Logger,
    testUserID string,
) *TestAuthMiddleware {
    return &TestAuthMiddleware{
        userService: userService,
        logger:      logger,
        testUserID:  testUserID,
    }
}

func (m *TestAuthMiddleware) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Only allow test middleware in non-production environments
            if config.GetEnvironment() == "production" {
                m.logger.Error("TestAuthMiddleware used in production, this is a security risk")
                errorResponse := domain.NewErrorResponse(
                    "InternalServerError",
                    "An unexpected error occurred", 
                    http.StatusInternalServerError,
                )
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusInternalServerError)
                json.NewEncoder(w).Encode(errorResponse)
                return
            }
            
            // Get test user
            user, err := m.userService.GetUserByID(r.Context(), m.testUserID)
            if err != nil {
                m.logger.Error("Failed to get test user",
                    zap.String("user_id", m.testUserID),
                    zap.Error(err),
                )
                next.ServeHTTP(w, r)
                return
            }
            
            // Set user in context
            ctx := context.WithValue(r.Context(), middleware.UserContextKey, user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## Standard Middleware Chain

To ensure consistent middleware application, we'll define a standard middleware chain:

```go
// internal/server/middleware_chain.go

// SetupStandardMiddleware configures and applies standard middleware to a router
func SetupStandardMiddleware(
    router chi.Router,
    factory *factory.MiddlewareFactory,
) {
    // Apply middleware in order of execution
    router.Use(middleware.RequestID)
    router.Use(factory.CreateErrorMiddleware().Middleware())
    router.Use(factory.CreateLoggingMiddleware().Middleware())
    router.Use(factory.CreateCORSMiddleware().Middleware())
    router.Use(factory.CreateRateLimiterMiddleware().Middleware())
    router.Use(factory.CreateAuthMiddleware().Middleware())
    router.Use(middleware.Recoverer)
}
```

## Middleware Configuration

Configuration for middleware should be defined in a structured way:

```go
// internal/config/middleware_config.go

type MiddlewareConfig struct {
    Auth      AuthConfig      `yaml:"auth"`
    CORS      CORSConfig      `yaml:"cors"`
    RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type AuthConfig struct {
    EnableDummyAuth bool   `yaml:"enable_dummy_auth"`
    DummyUserID     string `yaml:"dummy_user_id"`
}

type CORSConfig struct {
    AllowedOrigins []string `yaml:"allowed_origins"`
    AllowedMethods []string `yaml:"allowed_methods"`
    AllowedHeaders []string `yaml:"allowed_headers"`
    MaxAge         int      `yaml:"max_age"`
}

type RateLimitConfig struct {
    IPLimits    []IPLimitConfig    `yaml:"ip_limits"`
    UserLimits  []UserLimitConfig  `yaml:"user_limits"`
    RouteLimits []RouteLimitConfig `yaml:"route_limits"`
}

type IPLimitConfig struct {
    Requests int           `yaml:"requests"`
    Duration time.Duration `yaml:"duration"`
}

type UserLimitConfig struct {
    Requests int           `yaml:"requests"`
    Duration time.Duration `yaml:"duration"`
}

type RouteLimitConfig struct {
    Method   string        `yaml:"method"`
    Path     string        `yaml:"path"`
    Requests int           `yaml:"requests"`
    Duration time.Duration `yaml:"duration"`
}
```

## Safety Mechanisms

To prevent test/mock middleware from being used in production:

1. All mock/test middleware must check the environment before applying:

```go
if config.GetEnvironment() == "production" {
    log.Error("Test middleware used in production environment")
    // Return error or fall back to real implementation
}
```

2. Factory methods should enforce this check:

```go
func (f *MiddlewareFactory) CreateTestAuthMiddleware() *middleware.TestAuthMiddleware {
    if f.config.Environment == "production" {
        f.logger.Fatal("Attempted to create TestAuthMiddleware in production",
            zap.String("environment", f.config.Environment),
        )
        return nil // This will never be reached due to Fatal, but satisfies the compiler
    }
    
    return middleware.NewTestAuthMiddleware(
        f.userService,
        f.logger,
        f.config.Auth.TestUserID,
    )
}
```

## Migration Plan

To migrate from the current multiple middleware implementations to the standardized approach:

1. Create the new consolidated middleware components
2. Create the middleware factory
3. Update route registration to use the factory
4. Replace individual middleware usages with consolidated versions
5. Remove deprecated middleware implementations

## Conclusion

This standardized middleware approach provides several benefits:

1. **Consistency**: All HTTP handlers use the same middleware chain
2. **Security**: Test/mock middleware cannot be used in production
3. **Maintainability**: Single implementation of each middleware type
4. **Flexibility**: Factory pattern allows for environment-specific configuration
5. **Testability**: Separate test middleware implementations make testing easier 
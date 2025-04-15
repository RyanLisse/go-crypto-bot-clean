# HTTP Layer Unification (Gin to Chi) - Implementation Plan

## Current State Analysis

1. **Router Definition**: 
   - The main entry point in `cmd/server/main.go` already uses Chi
   - A Gin router also exists in `internal/adapter/delivery/http/router.go`

2. **Handler Types**:
   - Some handlers use Chi (e.g., `position_handler_chi.go`, `market_data.go`)
   - Others use Gin (e.g., `position_handler.go`, `risk_handler.go`, `autobuy_handler.go`)

3. **Middleware**:
   - Chi middleware exists in `internal/adapter/http/middleware/middleware.go`
   - Gin middleware exists in various handler files and in the router.go

## Implementation Steps

### 1. Remove or Rename Gin-based Router

```bash
# Delete Gin router
rm internal/adapter/delivery/http/router.go
```

### 2. Convert Each Gin Handler to Chi

For each handler using Gin, follow this pattern:

#### Example Migration Pattern (converting risk_handler.go)

Original Gin-based handler:
```go
// RegisterRoutes registers risk-related routes with the Gin engine
func (h *RiskHandler) RegisterRoutes(router *gin.RouterGroup) {
    riskGroup := router.Group("/risk")
    {
        riskGroup.GET("/profile", h.GetRiskProfile)
        // ...other routes
    }
}

// GetRiskProfile handles profile requests
func (h *RiskHandler) GetRiskProfile(c *gin.Context) {
    userID := getUserIDFromContext(c.Request)
    // ...handler logic
    c.JSON(http.StatusOK, response.Success(profile))
}
```

New Chi-based handler:
```go
// RegisterRoutes registers risk-related routes with the Chi router
func (h *RiskHandler) RegisterRoutes(r chi.Router) {
    r.Route("/risk", func(r chi.Router) {
        r.Get("/profile", h.GetRiskProfile)
        // ...other routes
    })
}

// GetRiskProfile handles profile requests
func (h *RiskHandler) GetRiskProfile(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r)
    // ...handler logic
    response.WriteJSON(w, http.StatusOK, response.Success(profile))
}
```

### 3. Specific Handler Files to Convert

1. **Position Handler** (priority since there are two versions):
   - Remove `position_handler.go` (Gin version)
   - Rename `position_handler_chi.go` to `position_handler.go`

2. **Risk Handler**:
   - Convert `risk_handler.go` from Gin to Chi

3. **Autobuy Handler**:
   - Convert `autobuy_handler.go` from Gin to Chi

4. **Websocket Handler**:
   - Convert `websocket.go` from Gin to Chi

### 4. Update Handler Tests

Test files like `position_handler_test.go` need to be updated to use Chi instead of Gin.

Example update:
```go
// Before: Gin setup
router := gin.Default()
router.GET("/positions/:id", handler.GetPositionByID)
req, _ := http.NewRequest("GET", "/positions/123", nil)
recorder := httptest.NewRecorder()
router.ServeHTTP(recorder, req)

// After: Chi setup
router := chi.NewRouter()
router.Get("/positions/{id}", handler.GetPositionByID)
req, _ := http.NewRequest("GET", "/positions/123", nil)
recorder := httptest.NewRecorder()
router.ServeHTTP(recorder, req)
```

### 5. Update Dependencies

1. Remove Gin dependency from go.mod:
```bash
go mod edit -droprequire=github.com/gin-gonic/gin
go mod tidy
```

### 6. Update Main Server File

Update `cmd/server/main.go` to register all converted handlers.

### 7. Execute Incremental Testing

1. Convert one handler at a time
2. Run tests after each conversion
3. Verify API endpoints work correctly

## Migration Checklist

- [x] Remove Gin router
- [x] Rename `position_handler_chi.go` to `position_handler.go`
- [x] Convert `risk_handler.go` to use Chi
- [x] Convert `autobuy_handler.go` to use Chi
- [x] Convert `websocket.go` to use Chi
- [x] Update handler tests to use Chi
- [ ] Remove Gin dependency
- [ ] Verify all API endpoints work correctly 
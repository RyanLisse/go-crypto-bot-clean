# Date-Based Coin Filtering Implementation Guide

## Overview

This document outlines the implementation of date-based filtering for new coin listings in the Go Crypto Bot project. This functionality allows users to retrieve new coin listings for a specific date or date range, similar to the existing feature in the Python project.

## Current Implementation Analysis

### Python Implementation

The Python project implements date-based filtering in two scripts:

1. **`check_new_coins_today.py`**: Retrieves new coins set to go online today by calling the MEXC API and filtering results based on the current date.
2. **`check_new_coins_database.py`**: Retrieves new coins from the database and filters them based on the current date.

Both scripts use the `NewCoin` component to fetch new coin data from the MEXC API and filter coins based on their `firstOpenTime` attribute.

Key code snippet from Python implementation:
```python
for coin in coin_statuses:
    if coin.firstOpenTime:
        open_time = datetime.datetime.fromtimestamp(int(coin.firstOpenTime)/1000)
        if open_time.date() == today_date:
            today_coins.append(coin)
```

### Go Implementation Status

The Go project currently has functionality to fetch new coin listings but lacks the specific date-based filtering feature. The current implementation includes:

1. **NewCoin Service**: Handles fetching and storing new coin listings
2. **REST API Endpoints**: For retrieving all new coins and processing them
3. **MEXC API Integration**: For fetching new coin data from the exchange

## Missing Components

The following components need to be implemented to add date-based filtering functionality:

1. **Request DTOs**: For date and date range parameters
2. **Service Methods**: To filter coins by date and date range
3. **Repository Methods**: To query the database for coins within a date range
4. **API Endpoints**: To expose the date filtering functionality
5. **Route Registration**: To make the new endpoints accessible

## Implementation Plan

### 1. Update Request DTOs

Add new DTOs to handle date-based filtering requests:

```go
// NewCoinsByDateRequest represents a request to get new coins by date
type NewCoinsByDateRequest struct {
    Date string `json:"date" form:"date" binding:"required"` // Format: YYYY-MM-DD
}

// NewCoinDateRangeRequest represents a request to get new coins within a date range
type NewCoinDateRangeRequest struct {
    StartDate string `json:"start_date" form:"start_date" binding:"required"` // Format: YYYY-MM-DD
    EndDate   string `json:"end_date" form:"end_date"`                        // Format: YYYY-MM-DD, optional (defaults to today)
}
```

### 2. Update NewCoinService Interface

Extend the `NewCoinService` interface with methods for date-based filtering:

```go
type NewCoinService interface {
    // Existing methods...
    
    // GetCoinsByDate returns new coins found on a specific date
    GetCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error)
    
    // GetCoinsByDateRange returns new coins found within a date range
    GetCoinsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error)
}
```

### 3. Implement Service Methods

Add implementations for both service structs:

```go
// For newCoinService
func (s *newCoinService) GetCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
    // Set time boundaries for the day
    startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
    endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)
    
    return s.newCoinRepo.FindByDateRange(ctx, startOfDay, endOfDay)
}

func (s *newCoinService) GetCoinsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error) {
    // Set time boundaries
    startOfDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
    endOfDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
    
    return s.newCoinRepo.FindByDateRange(ctx, startOfDay, endOfDay)
}

// Similar implementations for mockCompatibleCoinService
```

### 4. Update Repository Interface and Implementation

Add a new method to the repository interface:

```go
type NewCoinRepository interface {
    // Existing methods...
    
    // FindByDateRange finds coins within a date range
    FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error)
}
```

Implement the method in the repository:

```go
func (r *newCoinRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error) {
    var coins []models.NewCoin
    
    result := r.db.WithContext(ctx).
        Where("found_at BETWEEN ? AND ?", startDate, endDate).
        Find(&coins)
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return coins, nil
}
```

### 5. Add Handler Methods

Implement new handler methods in the `NewCoinsHandler`:

```go
// GetCoinsByDate godoc
// @Summary Get coins by date
// @Description Returns a list of coins found on a specific date
// @Tags newcoins
// @Accept json
// @Produce json
// @Param date query string true "Date in YYYY-MM-DD format"
// @Success 200 {object} response.NewCoinsListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins/by-date [get]
func (h *NewCoinsHandler) GetCoinsByDate(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Parse request
    var req request.NewCoinsByDateRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.ErrorResponse{
            Code:    "invalid_request",
            Message: "Invalid date format",
            Details: "Date must be in YYYY-MM-DD format",
        })
        return
    }
    
    // Parse date
    date, err := time.Parse("2006-01-02", req.Date)
    if err != nil {
        c.JSON(http.StatusBadRequest, response.ErrorResponse{
            Code:    "invalid_date",
            Message: "Invalid date format",
            Details: err.Error(),
        })
        return
    }
    
    // Get coins by date
    coins, err := h.NewCoinService.GetCoinsByDate(ctx, date)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.ErrorResponse{
            Code:    "internal_error",
            Message: "Failed to get coins by date",
            Details: err.Error(),
        })
        return
    }
    
    // Map to response DTOs and build response
    // ...
}

// GetCoinsByDateRange godoc
// @Summary Get coins by date range
// @Description Returns a list of coins found within a date range
// @Tags newcoins
// @Accept json
// @Produce json
// @Param start_date query string true "Start date in YYYY-MM-DD format"
// @Param end_date query string false "End date in YYYY-MM-DD format (defaults to today)"
// @Success 200 {object} response.NewCoinsListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/newcoins/by-date-range [get]
func (h *NewCoinsHandler) GetCoinsByDateRange(c *gin.Context) {
    // Implementation similar to GetCoinsByDate but handling date range
    // ...
}
```

### 6. Register Routes

Update the route registration:

```go
newCoinsGroup := v1.Group("/newcoins")
{
    // Existing routes...
    
    // Add new routes
    newCoinsGroup.GET("/by-date", newCoinsHandler.GetCoinsByDate)
    newCoinsGroup.GET("/by-date-range", newCoinsHandler.GetCoinsByDateRange)
}
```

## Testing

The new endpoints can be tested with curl:

```bash
# Get coins for a specific date
curl -v "http://localhost:8080/api/v1/newcoins/by-date?date=2025-04-08"

# Get coins for a date range
curl -v "http://localhost:8080/api/v1/newcoins/by-date-range?start_date=2025-04-01&end_date=2025-04-08"
```

## Comparison with Python Implementation

The Go implementation follows a similar approach to the Python version but with these improvements:

1. **Structured Architecture**: Clear separation of concerns (handler, service, repository)
2. **Error Handling**: Comprehensive error handling at each layer
3. **Context Propagation**: Proper context propagation for cancellation and timeouts
4. **API Documentation**: Swagger annotations for API documentation
5. **Date Range Support**: Added support for filtering by date range, not just a single date
6. **Query Parameters**: Using query parameters for better RESTful API design

## Future Enhancements

Potential future enhancements to consider:

1. **Caching**: Add caching for frequently accessed date ranges
2. **Pagination**: Implement pagination for large result sets
3. **Sorting Options**: Allow sorting by different fields (volume, symbol, etc.)
4. **Filtering Options**: Add more filtering options (by status, volume threshold, etc.)
5. **Webhooks**: Implement webhooks for real-time notifications of new coins

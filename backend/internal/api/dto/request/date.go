package request

// DateRequest represents a request with a specific date
type DateRequest struct {
	Date string `json:"date" binding:"required"`
}

// DateRangeRequest represents a request with a date range
type DateRangeRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

package request

// ProcessNewCoinsRequest represents a request to process new coins
type ProcessNewCoinsRequest struct {
	CoinIDs []uint `json:"coin_ids" binding:"required"`
}

// DateFilterRequest represents a request to filter coins by a specific date
type DateFilterRequest struct {
	Date string `json:"date" binding:"required" example:"2023-01-01"`
}

// DateRangeFilterRequest represents a request to filter coins by a date range
type DateRangeFilterRequest struct {
	StartDate string `json:"start_date" binding:"required" example:"2023-01-01"`
	EndDate   string `json:"end_date" binding:"required" example:"2023-01-31"`
}

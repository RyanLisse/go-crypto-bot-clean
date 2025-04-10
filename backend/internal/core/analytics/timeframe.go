package analytics

// TimeFrame represents the granularity of analytics data
type TimeFrame string

const (
	TimeFrameAll     TimeFrame = "all"
	TimeFrameDay     TimeFrame = "day"
	TimeFrameWeek    TimeFrame = "week"
	TimeFrameMonth   TimeFrame = "month"
	TimeFrameQuarter TimeFrame = "quarter"
	TimeFrameYear    TimeFrame = "year"
)

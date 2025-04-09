package huma

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"go-crypto-bot-clean/backend/internal/domain/models"
)

// BacktestRequest represents a request to run a backtest
type BacktestRequest struct {
	Body struct {
		Strategy       string    `json:"strategy" doc:"Strategy to use for backtesting" example:"breakout" binding:"required"`
		Symbol         string    `json:"symbol" doc:"Trading pair symbol" example:"BTC/USDT" binding:"required"`
		Timeframe      string    `json:"timeframe" doc:"Timeframe for backtesting" example:"1h" enum:"1m,5m,15m,30m,1h,4h,1d" binding:"required"`
		StartDate      time.Time `json:"startDate" doc:"Start date for backtesting" example:"2023-01-01T00:00:00Z" binding:"required"`
		EndDate        time.Time `json:"endDate" doc:"End date for backtesting" example:"2023-02-01T00:00:00Z" binding:"required"`
		InitialCapital float64   `json:"initialCapital" doc:"Initial capital for backtesting" example:"1000.0" minimum:"0" binding:"required"`
		RiskPerTrade   float64   `json:"riskPerTrade" doc:"Risk per trade as a percentage" example:"0.02" minimum:"0" maximum:"1" binding:"required"`
	}
}

// BacktestResponse represents the response from a backtest
type BacktestResponse struct {
	Body struct {
		ID                 string                 `json:"id" doc:"Unique identifier for the backtest" example:"bt-123456"`
		Strategy           string                 `json:"strategy" doc:"Strategy used for backtesting" example:"breakout"`
		Symbol             string                 `json:"symbol" doc:"Trading pair symbol" example:"BTC/USDT"`
		Timeframe          string                 `json:"timeframe" doc:"Timeframe used for backtesting" example:"1h"`
		StartDate          time.Time              `json:"startDate" doc:"Start date for backtesting" example:"2023-01-01T00:00:00Z"`
		EndDate            time.Time              `json:"endDate" doc:"End date for backtesting" example:"2023-02-01T00:00:00Z"`
		InitialCapital     float64                `json:"initialCapital" doc:"Initial capital for backtesting" example:"1000.0"`
		FinalCapital       float64                `json:"finalCapital" doc:"Final capital after backtesting" example:"1250.0"`
		TotalReturn        float64                `json:"totalReturn" doc:"Total return percentage" example:"25.0"`
		AnnualizedReturn   float64                `json:"annualizedReturn" doc:"Annualized return percentage" example:"300.0"`
		MaxDrawdown        float64                `json:"maxDrawdown" doc:"Maximum drawdown percentage" example:"15.0"`
		SharpeRatio        float64                `json:"sharpeRatio" doc:"Sharpe ratio" example:"1.5"`
		WinRate            float64                `json:"winRate" doc:"Win rate percentage" example:"65.0"`
		ProfitFactor       float64                `json:"profitFactor" doc:"Profit factor" example:"2.1"`
		TotalTrades        int                    `json:"totalTrades" doc:"Total number of trades" example:"50"`
		WinningTrades      int                    `json:"winningTrades" doc:"Number of winning trades" example:"32"`
		LosingTrades       int                    `json:"losingTrades" doc:"Number of losing trades" example:"18"`
		AverageProfitTrade float64                `json:"averageProfitTrade" doc:"Average profit per winning trade" example:"12.5"`
		AverageLossTrade   float64                `json:"averageLossTrade" doc:"Average loss per losing trade" example:"-8.2"`
		MaxConsecutiveWins int                    `json:"maxConsecutiveWins" doc:"Maximum consecutive winning trades" example:"8"`
		MaxConsecutiveLoss int                    `json:"maxConsecutiveLoss" doc:"Maximum consecutive losing trades" example:"3"`
		EquityCurve        []EquityPoint          `json:"equityCurve" doc:"Equity curve data points"`
		DrawdownCurve      []DrawdownPoint        `json:"drawdownCurve" doc:"Drawdown curve data points"`
		Trades             []models.Order         `json:"trades" doc:"List of trades executed during backtesting"`
		CreatedAt          time.Time              `json:"createdAt" doc:"Time when the backtest was created" example:"2023-02-02T10:00:00Z"`
	}
}

// EquityPoint represents a point on the equity curve
type EquityPoint struct {
	Timestamp time.Time `json:"timestamp" doc:"Timestamp of the equity point" example:"2023-01-01T12:00:00Z"`
	Equity    float64   `json:"equity" doc:"Equity value at this point" example:"1050.0"`
}

// DrawdownPoint represents a point on the drawdown curve
type DrawdownPoint struct {
	Timestamp time.Time `json:"timestamp" doc:"Timestamp of the drawdown point" example:"2023-01-01T12:00:00Z"`
	Drawdown  float64   `json:"drawdown" doc:"Drawdown percentage at this point" example:"5.0"`
}

// BacktestListResponse represents a list of backtests
type BacktestListResponse struct {
	Body struct {
		Backtests []struct {
			ID               string    `json:"id" doc:"Unique identifier for the backtest" example:"bt-123456"`
			Strategy         string    `json:"strategy" doc:"Strategy used for backtesting" example:"breakout"`
			Symbol           string    `json:"symbol" doc:"Trading pair symbol" example:"BTC/USDT"`
			Timeframe        string    `json:"timeframe" doc:"Timeframe used for backtesting" example:"1h"`
			StartDate        time.Time `json:"startDate" doc:"Start date for backtesting" example:"2023-01-01T00:00:00Z"`
			EndDate          time.Time `json:"endDate" doc:"End date for backtesting" example:"2023-02-01T00:00:00Z"`
			TotalReturn      float64   `json:"totalReturn" doc:"Total return percentage" example:"25.0"`
			MaxDrawdown      float64   `json:"maxDrawdown" doc:"Maximum drawdown percentage" example:"15.0"`
			SharpeRatio      float64   `json:"sharpeRatio" doc:"Sharpe ratio" example:"1.5"`
			WinRate          float64   `json:"winRate" doc:"Win rate percentage" example:"65.0"`
			TotalTrades      int       `json:"totalTrades" doc:"Total number of trades" example:"50"`
			CreatedAt        time.Time `json:"createdAt" doc:"Time when the backtest was created" example:"2023-02-02T10:00:00Z"`
		} `json:"backtests" doc:"List of backtests"`
		Count     int       `json:"count" doc:"Number of backtests" example:"10"`
		Timestamp time.Time `json:"timestamp" doc:"Timestamp of the response" example:"2023-02-02T10:00:00Z"`
	}
}

// BacktestCompareRequest represents a request to compare backtests
type BacktestCompareRequest struct {
	Body struct {
		BacktestIDs []string `json:"backtest_ids" doc:"IDs of backtests to compare" example:"[\"bt-123456\", \"bt-789012\"]" binding:"required"`
	}
}

// BacktestCompareResponse represents the response from a backtest comparison
type BacktestCompareResponse struct {
	Body struct {
		Backtests []struct {
			ID               string  `json:"id" doc:"Unique identifier for the backtest" example:"bt-123456"`
			Strategy         string  `json:"strategy" doc:"Strategy used for backtesting" example:"breakout"`
			Symbol           string  `json:"symbol" doc:"Trading pair symbol" example:"BTC/USDT"`
			Timeframe        string  `json:"timeframe" doc:"Timeframe used for backtesting" example:"1h"`
			TotalReturn      float64 `json:"totalReturn" doc:"Total return percentage" example:"25.0"`
			AnnualizedReturn float64 `json:"annualizedReturn" doc:"Annualized return percentage" example:"300.0"`
			MaxDrawdown      float64 `json:"maxDrawdown" doc:"Maximum drawdown percentage" example:"15.0"`
			SharpeRatio      float64 `json:"sharpeRatio" doc:"Sharpe ratio" example:"1.5"`
			WinRate          float64 `json:"winRate" doc:"Win rate percentage" example:"65.0"`
			ProfitFactor     float64 `json:"profitFactor" doc:"Profit factor" example:"2.1"`
			TotalTrades      int     `json:"totalTrades" doc:"Total number of trades" example:"50"`
		} `json:"backtests" doc:"List of backtests being compared"`
		Comparison struct {
			BestTotalReturn      string  `json:"bestTotalReturn" doc:"ID of backtest with best total return" example:"bt-123456"`
			BestSharpeRatio      string  `json:"bestSharpeRatio" doc:"ID of backtest with best Sharpe ratio" example:"bt-789012"`
			BestDrawdown         string  `json:"bestDrawdown" doc:"ID of backtest with best (lowest) drawdown" example:"bt-123456"`
			BestWinRate          string  `json:"bestWinRate" doc:"ID of backtest with best win rate" example:"bt-789012"`
			BestProfitFactor     string  `json:"bestProfitFactor" doc:"ID of backtest with best profit factor" example:"bt-123456"`
			ReturnDifference     float64 `json:"returnDifference" doc:"Difference in total return between best and worst" example:"10.5"`
			DrawdownDifference   float64 `json:"drawdownDifference" doc:"Difference in max drawdown between best and worst" example:"5.2"`
			SharpeRatioDifference float64 `json:"sharpeRatioDifference" doc:"Difference in Sharpe ratio between best and worst" example:"0.8"`
		} `json:"comparison" doc:"Comparison metrics between backtests"`
		Timestamp time.Time `json:"timestamp" doc:"Timestamp of the response" example:"2023-02-02T10:00:00Z"`
	}
}

// registerBacktestEndpoints registers the backtest endpoints.
func registerBacktestEndpoints(api huma.API, basePath string) {
	// POST /backtest
	huma.Register(api, huma.Operation{
		OperationID: "run-backtest",
		Method:      http.MethodPost,
		Path:        basePath + "/backtest",
		Summary:     "Run a backtest",
		Description: "Runs a backtest with the specified parameters",
		Tags:        []string{"Backtesting"},
	}, func(ctx context.Context, input *BacktestRequest) (*struct {
		Body struct {
			ID        string    `json:"id" doc:"Unique identifier for the backtest" example:"bt-123456"`
			Status    string    `json:"status" doc:"Status of the backtest" example:"completed" enum:"running,completed,failed"`
			Message   string    `json:"message,omitempty" doc:"Status message" example:"Backtest completed successfully"`
			CreatedAt time.Time `json:"createdAt" doc:"Time when the backtest was created" example:"2023-02-02T10:00:00Z"`
		}
	}, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// GET /backtest/{id}
	huma.Register(api, huma.Operation{
		OperationID: "get-backtest-result",
		Method:      http.MethodGet,
		Path:        basePath + "/backtest/{id}",
		Summary:     "Get backtest result",
		Description: "Returns the result of a specific backtest",
		Tags:        []string{"Backtesting"},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"Backtest ID" example:"bt-123456"`
	}) (*BacktestResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// GET /backtest/list
	huma.Register(api, huma.Operation{
		OperationID: "list-backtests",
		Method:      http.MethodGet,
		Path:        basePath + "/backtest/list",
		Summary:     "List backtests",
		Description: "Returns a list of backtests",
		Tags:        []string{"Backtesting"},
	}, func(ctx context.Context, input *struct {
		Limit    int    `query:"limit" doc:"Maximum number of backtests to return" example:"10" default:"10" minimum:"1" maximum:"100"`
		Offset   int    `query:"offset" doc:"Offset for pagination" example:"0" default:"0" minimum:"0"`
		Strategy string `query:"strategy" doc:"Filter by strategy" example:"breakout"`
		Symbol   string `query:"symbol" doc:"Filter by symbol" example:"BTC/USDT"`
	}) (*BacktestListResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})

	// POST /backtest/compare
	huma.Register(api, huma.Operation{
		OperationID: "compare-backtests",
		Method:      http.MethodPost,
		Path:        basePath + "/backtest/compare",
		Summary:     "Compare backtests",
		Description: "Compares multiple backtests",
		Tags:        []string{"Backtesting"},
	}, func(ctx context.Context, input *BacktestCompareRequest) (*BacktestCompareResponse, error) {
		// This is just a placeholder for documentation purposes
		return nil, nil
	})
}

package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo/crypto-bot/internal/adapter/http/response"
	"github.com/neo/crypto-bot/internal/domain/model/market"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
)

// MarketDataHandler handles HTTP requests for market data
type MarketDataHandler struct {
	useCase *usecase.MarketDataUseCase
	logger  *zerolog.Logger
}

// NewMarketDataHandler creates a new MarketDataHandler
func NewMarketDataHandler(uc *usecase.MarketDataUseCase, logger *zerolog.Logger) *MarketDataHandler {
	return &MarketDataHandler{
		useCase: uc,
		logger:  logger,
	}
}

// RegisterRoutes registers market data routes with the Gin engine
func (h *MarketDataHandler) RegisterRoutes(router *gin.RouterGroup) {
	marketGroup := router.Group("/market")
	{
		marketGroup.GET("/tickers", h.GetLatestTickers)
		marketGroup.GET("/tickers/:exchange/:symbol", h.GetTicker)
		marketGroup.GET("/candles", h.GetCandles)
		marketGroup.GET("/symbols", h.GetAllSymbols)
		marketGroup.GET("/symbols/:symbol", h.GetSymbolInfo)
	}
}

// GetLatestTickers godoc
// @Summary Get latest tickers
// @Description Get the latest tickers for all available symbols
// @Tags market
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=[]market.Ticker}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /market/tickers [get]
func (h *MarketDataHandler) GetLatestTickers(c *gin.Context) {
	tickers, err := h.useCase.GetLatestTickers(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get latest tickers")
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to retrieve latest tickers"))
		return
	}
	c.JSON(http.StatusOK, response.Success(tickers))
}

// GetTicker godoc
// @Summary Get specific ticker
// @Description Get the latest ticker for a specific exchange and symbol
// @Tags market
// @Accept json
// @Produce json
// @Param exchange path string true "Exchange name (e.g., binance)"
// @Param symbol path string true "Trading symbol (e.g., BTCUSDT)"
// @Success 200 {object} response.APIResponse{data=market.Ticker}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /market/tickers/{exchange}/{symbol} [get]
func (h *MarketDataHandler) GetTicker(c *gin.Context) {
	exchange := c.Param("exchange")
	symbol := c.Param("symbol")

	if exchange == "" || symbol == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Exchange and symbol parameters are required"))
		return
	}

	ticker, err := h.useCase.GetTicker(c.Request.Context(), exchange, symbol)
	if err != nil {
		h.logger.Error().Err(err).Str("exchange", exchange).Str("symbol", symbol).Msg("Failed to get ticker")
		// TODO: Differentiate between Not Found and other errors
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to retrieve ticker"))
		return
	}

	if ticker == nil {
		c.JSON(http.StatusNotFound, response.Error(response.ErrorCodeNotFound, "Ticker not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(ticker))
}

// GetCandles godoc
// @Summary Get candles (k-lines)
// @Description Get historical candles for a specific exchange, symbol, and interval
// @Tags market
// @Accept json
// @Produce json
// @Param exchange query string true "Exchange name (e.g., binance)"
// @Param symbol query string true "Trading symbol (e.g., BTCUSDT)"
// @Param interval query string true "Candle interval (e.g., 1m, 5m, 1h, 1d)"
// @Param start query integer false "Start timestamp (Unix milliseconds)"
// @Param end query integer false "End timestamp (Unix milliseconds)"
// @Param limit query integer false "Limit number of candles (default: 500, max: 1000)"
// @Success 200 {object} response.APIResponse{data=[]market.Candle}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /market/candles [get]
func (h *MarketDataHandler) GetCandles(c *gin.Context) {
	exchange := c.Query("exchange")
	symbol := c.Query("symbol")
	intervalStr := c.Query("interval")

	if exchange == "" || symbol == "" || intervalStr == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "exchange, symbol, and interval parameters are required"))
		return
	}

	// Validate interval string
	valid := false
	validIntervals := []market.Interval{market.Interval1m, market.Interval5m, market.Interval15m, market.Interval30m,
		market.Interval1h, market.Interval4h, market.Interval1d, market.Interval1w, market.Interval1M}
	for _, validInterval := range validIntervals {
		if market.Interval(intervalStr) == validInterval {
			valid = true
			break
		}
	}
	if !valid {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid interval. Supported intervals: 1m, 5m, 15m, 30m, 1h, 4h, 1d, 1w, 1M"))
		return
	}

	interval := market.Interval(intervalStr)

	// Default time range: last 24 hours
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	if startStr := c.Query("start"); startStr != "" {
		startMillis, err := strconv.ParseInt(startStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid start timestamp"))
			return
		}
		startTime = time.UnixMilli(startMillis)
	}

	if endStr := c.Query("end"); endStr != "" {
		endMillis, err := strconv.ParseInt(endStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid end timestamp"))
			return
		}
		endTime = time.UnixMilli(endMillis)
	}

	limit := 500 // Default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Invalid limit parameter"))
			return
		}
		if parsedLimit > 1000 { // Max limit
			parsedLimit = 1000
		}
		limit = parsedLimit
	}

	candles, err := h.useCase.GetCandles(c.Request.Context(), exchange, symbol, interval, startTime, endTime, limit)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get candles")
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to retrieve candles"))
		return
	}

	c.JSON(http.StatusOK, response.Success(candles))
}

// GetAllSymbols godoc
// @Summary Get all symbols
// @Description Get a list of all available trading symbols
// @Tags market
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse{data=[]market.Symbol}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /market/symbols [get]
func (h *MarketDataHandler) GetAllSymbols(c *gin.Context) {
	symbols, err := h.useCase.GetAllSymbols(c.Request.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get all symbols")
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to retrieve symbols"))
		return
	}
	c.JSON(http.StatusOK, response.Success(symbols))
}

// GetSymbolInfo godoc
// @Summary Get symbol info
// @Description Get detailed information for a specific trading symbol
// @Tags market
// @Accept json
// @Produce json
// @Param symbol path string true "Trading symbol (e.g., BTCUSDT)"
// @Success 200 {object} response.APIResponse{data=market.Symbol}
// @Failure 400 {object} response.APIResponse{error=response.APIError}
// @Failure 404 {object} response.APIResponse{error=response.APIError}
// @Failure 500 {object} response.APIResponse{error=response.APIError}
// @Router /market/symbols/{symbol} [get]
func (h *MarketDataHandler) GetSymbolInfo(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, response.Error(response.ErrorCodeBadRequest, "Symbol parameter is required"))
		return
	}

	symbolInfo, err := h.useCase.GetSymbolInfo(c.Request.Context(), symbol)
	if err != nil {
		// TODO: Differentiate between Not Found and other errors
		h.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get symbol info")
		c.JSON(http.StatusInternalServerError, response.Error(response.ErrorCodeInternalError, "Failed to retrieve symbol information"))
		return
	}

	if symbolInfo == nil { // Should be handled by use case/repo, but double check
		c.JSON(http.StatusNotFound, response.Error(response.ErrorCodeNotFound, "Symbol not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(symbolInfo))
}

// For backward compatibility
type ErrorResponse struct {
	Message string `json:"message"`
}

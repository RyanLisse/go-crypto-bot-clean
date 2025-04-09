// Package handlers contains HTTP request handlers.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryanlisse/go-crypto-bot/internal/core/newcoin"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/service"
)

// CoinHandler handles market and coin-related endpoints.
type CoinHandler struct {
	exchangeService service.ExchangeService
	newCoinService  newcoin.NewCoinService
}

// NewCoinHandler creates a new CoinHandler with dependencies.
func NewCoinHandler(exchangeService service.ExchangeService, newCoinService newcoin.NewCoinService) *CoinHandler {
	return &CoinHandler{
		exchangeService: exchangeService,
		newCoinService:  newCoinService,
	}
}

// MarketResponse represents a market summary.
type MarketResponse struct {
	Symbol      string  `json:"symbol"`
	Price       float64 `json:"price"`
	Volume      float64 `json:"volume"`
	QuoteVolume float64 `json:"quote_volume"`
}

// NewCoinResponse represents a new coin notification.
type NewCoinResponse struct {
	Symbol           string  `json:"symbol"`
	FoundAt          int64   `json:"found_at"`
	BaseVolume       float64 `json:"base_volume"`
	QuoteVolume      float64 `json:"quote_volume"`
	Status           string  `json:"status,omitempty"`
	BecameTradableAt int64   `json:"became_tradable_at,omitempty"`
}

// @Summary List all supported markets
// @Description Returns a list of all supported markets
// @Tags markets
// @Produce json
// @Success 200 {array} MarketResponse
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/v1/markets [get]
func (h *CoinHandler) ListMarkets(c *gin.Context) {
	tickers, err := h.exchangeService.GetAllTickers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch markets"})
		return
	}

	var resp []MarketResponse
	for _, m := range tickers {
		resp = append(resp, MarketResponse{
			Symbol:      m.Symbol,
			Price:       m.Price,
			Volume:      m.Volume,
			QuoteVolume: m.QuoteVolume,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Get market data for a symbol
// @Description Returns market data for a specific symbol
// @Tags markets
// @Produce json
// @Param symbol path string true "Market symbol"
// @Success 200 {object} MarketResponse
// @Failure 404 {object} gin.H{"error": "market not found"}
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/v1/market/{symbol} [get]
func (h *CoinHandler) GetMarket(c *gin.Context) {
	symbol := c.Param("symbol")
	market, err := h.exchangeService.GetTicker(c.Request.Context(), symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch market"})
		return
	}
	if market == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "market not found"})
		return
	}
	resp := MarketResponse{
		Symbol:      market.Symbol,
		Price:       market.Price,
		Volume:      market.Volume,
		QuoteVolume: market.QuoteVolume,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary List new coin notifications
// @Description Returns a list of new coin notifications
// @Tags newcoins
// @Produce json
// @Success 200 {array} NewCoinResponse
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/v1/newcoins [get]
func (h *CoinHandler) ListNewCoins(c *gin.Context) {
	coins, err := h.newCoinService.DetectNewCoins(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch new coins"})
		return
	}

	// Initialize with empty slice instead of nil to ensure [] instead of null in JSON
	resp := make([]NewCoinResponse, 0)
	for _, coin := range coins {
		resp = append(resp, NewCoinResponse{
			Symbol:      coin.Symbol,
			FoundAt:     coin.FoundAt.Unix(),
			BaseVolume:  coin.BaseVolume,
			QuoteVolume: coin.QuoteVolume,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary List coins that became tradable today
// @Description Returns a list of coins that became tradable today
// @Tags newcoins
// @Produce json
// @Success 200 {array} NewCoinResponse
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/v1/newcoins/tradable/today [get]
func (h *CoinHandler) ListTradableCoinsToday(c *gin.Context) {
	coins, err := h.newCoinService.GetTradableCoinsToday(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tradable coins"})
		return
	}

	// Initialize with empty slice instead of nil to ensure [] instead of null in JSON
	resp := make([]NewCoinResponse, 0)
	for _, coin := range coins {
		var becameTradableAt int64
		if !coin.BecameTradableAt.IsZero() {
			becameTradableAt = coin.BecameTradableAt.Unix()
		}

		resp = append(resp, NewCoinResponse{
			Symbol:           coin.Symbol,
			FoundAt:          coin.FoundAt.Unix(),
			BaseVolume:       coin.BaseVolume,
			QuoteVolume:      coin.QuoteVolume,
			Status:           coin.Status,
			BecameTradableAt: becameTradableAt,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary List all tradable coins
// @Description Returns a list of all coins that have become tradable
// @Tags newcoins
// @Produce json
// @Success 200 {array} NewCoinResponse
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/v1/newcoins/tradable [get]
func (h *CoinHandler) ListTradableCoins(c *gin.Context) {
	coins, err := h.newCoinService.GetTradableCoins(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tradable coins"})
		return
	}

	// Initialize with empty slice instead of nil to ensure [] instead of null in JSON
	resp := make([]NewCoinResponse, 0)
	for _, coin := range coins {
		var becameTradableAt int64
		if !coin.BecameTradableAt.IsZero() {
			becameTradableAt = coin.BecameTradableAt.Unix()
		}

		resp = append(resp, NewCoinResponse{
			Symbol:           coin.Symbol,
			FoundAt:          coin.FoundAt.Unix(),
			BaseVolume:       coin.BaseVolume,
			QuoteVolume:      coin.QuoteVolume,
			Status:           coin.Status,
			BecameTradableAt: becameTradableAt,
		})
	}
	c.JSON(http.StatusOK, resp)
}

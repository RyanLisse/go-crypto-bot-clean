package backtest

import (
	"sort"
	"time"
)

// FeeModel defines the interface for calculating trading fees
type FeeModel interface {
	// CalculateFee calculates the fee for a trade
	CalculateFee(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64
}

// NoFee implements the FeeModel interface with no fees
type NoFee struct{}

// CalculateFee calculates the fee for a trade (always 0)
func (f *NoFee) CalculateFee(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64 {
	return 0.0
}

// FixedFeeModel implements the FeeModel interface with a fixed fee rate
type FixedFeeModel struct {
	FeeRate float64
}

// NewFixedFeeModel creates a new FixedFeeModel
func NewFixedFeeModel(feeRate float64) *FixedFeeModel {
	return &FixedFeeModel{
		FeeRate: feeRate,
	}
}

// CalculateFee calculates the fee for a trade
func (f *FixedFeeModel) CalculateFee(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64 {
	return price * quantity * f.FeeRate
}

// FeeTier represents a fee tier based on trading volume
type FeeTier struct {
	MinVolume float64
	MaxVolume float64
	MakerFee  float64
	TakerFee  float64
}

// TieredFeeModel implements the FeeModel interface with tiered fee rates based on volume
type TieredFeeModel struct {
	Tiers        []FeeTier
	Volume30Day  map[string]float64 // 30-day volume by symbol
	LastUpdated  time.Time
	DefaultTier  FeeTier
	VolumeByDate map[string]map[time.Time]float64 // Volume by symbol and date
}

// NewTieredFeeModel creates a new TieredFeeModel
func NewTieredFeeModel(tiers []FeeTier, defaultTier FeeTier) *TieredFeeModel {
	return &TieredFeeModel{
		Tiers:        tiers,
		Volume30Day:  make(map[string]float64),
		LastUpdated:  time.Now(),
		DefaultTier:  defaultTier,
		VolumeByDate: make(map[string]map[time.Time]float64),
	}
}

// CalculateFee calculates the fee for a trade
func (f *TieredFeeModel) CalculateFee(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64 {
	// Update 30-day volume
	f.updateVolume(symbol, quantity*price, timestamp)

	// Get 30-day volume for the symbol
	volume := f.Volume30Day[symbol]

	// Find the appropriate tier
	tier := f.findTier(volume)

	// Calculate fee based on side (maker or taker)
	var feeRate float64
	if side == "BUY" {
		feeRate = tier.TakerFee
	} else {
		feeRate = tier.MakerFee
	}

	return price * quantity * feeRate
}

// updateVolume updates the 30-day volume for a symbol
func (f *TieredFeeModel) updateVolume(symbol string, volume float64, timestamp time.Time) {
	// Initialize volume map for the symbol if it doesn't exist
	if _, ok := f.VolumeByDate[symbol]; !ok {
		f.VolumeByDate[symbol] = make(map[time.Time]float64)
	}

	// Add volume for the current date
	date := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
	f.VolumeByDate[symbol][date] += volume

	// Calculate 30-day volume
	thirtyDaysAgo := timestamp.AddDate(0, 0, -30)
	totalVolume := 0.0

	for date, vol := range f.VolumeByDate[symbol] {
		if date.After(thirtyDaysAgo) || date.Equal(thirtyDaysAgo) {
			totalVolume += vol
		}
	}

	f.Volume30Day[symbol] = totalVolume
	f.LastUpdated = timestamp
}

// findTier finds the appropriate fee tier for a given volume
func (f *TieredFeeModel) findTier(volume float64) FeeTier {
	// Sort tiers by min volume
	sort.Slice(f.Tiers, func(i, j int) bool {
		return f.Tiers[i].MinVolume < f.Tiers[j].MinVolume
	})

	// Find the appropriate tier
	for _, tier := range f.Tiers {
		if volume >= tier.MinVolume && (tier.MaxVolume == 0 || volume < tier.MaxVolume) {
			return tier
		}
	}

	// Return default tier if no matching tier is found
	return f.DefaultTier
}

// ExchangeFeeModel implements the FeeModel interface with exchange-specific fee rates
type ExchangeFeeModel struct {
	ExchangeName string
	FeeRates     map[string]float64 // Fee rates by symbol
	DefaultRate  float64
}

// NewExchangeFeeModel creates a new ExchangeFeeModel
func NewExchangeFeeModel(exchangeName string, defaultRate float64) *ExchangeFeeModel {
	return &ExchangeFeeModel{
		ExchangeName: exchangeName,
		FeeRates:     make(map[string]float64),
		DefaultRate:  defaultRate,
	}
}

// SetFeeRate sets the fee rate for a symbol
func (f *ExchangeFeeModel) SetFeeRate(symbol string, feeRate float64) {
	f.FeeRates[symbol] = feeRate
}

// CalculateFee calculates the fee for a trade
func (f *ExchangeFeeModel) CalculateFee(symbol string, side string, quantity float64, price float64, timestamp time.Time) float64 {
	// Get fee rate for the symbol
	feeRate, ok := f.FeeRates[symbol]
	if !ok {
		feeRate = f.DefaultRate
	}

	return price * quantity * feeRate
}

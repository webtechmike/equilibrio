package services

import (
	"context"
	"equilibrio-backend/internal/models"
	"time"
)

// MarketDataProvider defines the interface for market data sources
type MarketDataProvider interface {
	// GetQuote fetches real-time quote for a symbol
	GetQuote(ctx context.Context, symbol string) (*models.Quote, error)

	// GetHistoricalPrices fetches historical price data
	GetHistoricalPrices(ctx context.Context, symbol string, days int) ([]models.CandlestickData, error)

	// SearchSymbols searches for symbols matching criteria
	SearchSymbols(ctx context.Context, query string) ([]string, error)

	// GetMarketSnapshot gets current market overview
	GetMarketSnapshot(ctx context.Context, symbols []string) (map[string]*models.Quote, error)
}

// Quote represents a real-time stock quote
type Quote struct {
	Symbol        string
	Name          string
	Price         float64
	Change        float64
	ChangePercent float64
	Volume        int64
	MarketCap     int64
	PERatio       float64
	DividendYield float64
	Week52High    float64
	Week52Low     float64
	Open          float64
	High          float64
	Low           float64
	PreviousClose float64
	Sector        string
	Industry      string
	Timestamp     time.Time
}

// EquilibriumCalculator calculates equilibrium zones and signals
type EquilibriumCalculator struct {
	lookbackPeriod int
}

// NewEquilibriumCalculator creates a new equilibrium calculator
func NewEquilibriumCalculator(lookbackPeriod int) *EquilibriumCalculator {
	return &EquilibriumCalculator{
		lookbackPeriod: lookbackPeriod,
	}
}

// CalculateEquilibrium calculates equilibrium zones from price history
func (ec *EquilibriumCalculator) CalculateEquilibrium(prices []models.CandlestickData, currentPrice float64) models.EquilibriumData {
	if len(prices) == 0 {
		return models.EquilibriumData{
			Zone:       "neutral",
			Strength:   0.5,
			Support:    currentPrice * 0.95,
			Resistance: currentPrice * 1.05,
		}
	}

	// Find significant support and resistance levels
	support, resistance := ec.findKeyLevels(prices)

	// Determine zone based on current price position
	zone := "neutral"
	strength := 0.5

	priceRange := resistance - support
	if priceRange > 0 {
		position := (currentPrice - support) / priceRange

		if position < 0.3 {
			zone = "support"
			strength = 0.3 - position // Stronger when closer to support
		} else if position > 0.7 {
			zone = "resistance"
			strength = position - 0.7 // Stronger when closer to resistance
		} else {
			zone = "neutral"
			strength = 0.5
		}
	}

	return models.EquilibriumData{
		Zone:       zone,
		Strength:   strength,
		Support:    support,
		Resistance: resistance,
	}
}

// findKeyLevels identifies support and resistance levels from price history
func (ec *EquilibriumCalculator) findKeyLevels(prices []models.CandlestickData) (float64, float64) {
	if len(prices) == 0 {
		return 0, 0
	}

	// Simple approach: use recent lows for support, recent highs for resistance
	var lows, highs []float64

	lookback := ec.lookbackPeriod
	if lookback > len(prices) {
		lookback = len(prices)
	}

	for i := len(prices) - lookback; i < len(prices); i++ {
		lows = append(lows, prices[i].Low)
		highs = append(highs, prices[i].High)
	}

	// Find average of lowest lows and highest highs
	support := findAvgOfLowest(lows, 3)
	resistance := findAvgOfHighest(highs, 3)

	return support, resistance
}

// Helper functions
func findAvgOfLowest(values []float64, n int) float64 {
	if len(values) == 0 {
		return 0
	}
	if n > len(values) {
		n = len(values)
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)

	// Simple bubble sort for small n
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	sum := 0.0
	for i := 0; i < n; i++ {
		sum += sorted[i]
	}
	return sum / float64(n)
}

func findAvgOfHighest(values []float64, n int) float64 {
	if len(values) == 0 {
		return 0
	}
	if n > len(values) {
		n = len(values)
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)

	// Simple bubble sort for small n
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] < sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	sum := 0.0
	for i := 0; i < n; i++ {
		sum += sorted[i]
	}
	return sum / float64(n)
}

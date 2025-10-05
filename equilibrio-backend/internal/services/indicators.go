package services

import (
	"math"

	"equilibrio-backend/internal/models"
)

type IndicatorService struct{}

func NewIndicatorService() *IndicatorService {
	return &IndicatorService{}
}

// CalculateIndicators calculates technical indicators for a given symbol
func (s *IndicatorService) CalculateIndicators(symbol string, period int) (*models.TechnicalIndicators, error) {
	// In a real implementation, this would:
	// 1. Fetch historical price data for the symbol
	// 2. Calculate technical indicators
	// 3. Return the results

	// For now, return mock calculations
	indicators := &models.TechnicalIndicators{
		RSI:            s.calculateRSI(period),
		StochRSI:       s.calculateStochRSI(period),
		HistoricRSIAvg: s.calculateHistoricRSIAvg(period),
		SMA50:          s.calculateSMA(period, 50),
		SMA200:         s.calculateSMA(period, 200),
		EMA20:          s.calculateEMA(period, 20),
		MACD:           s.calculateMACD(period),
		MACDSignal:     s.calculateMACDSignal(period),
		MACDHistogram:  s.calculateMACDHistogram(period),
	}

	return indicators, nil
}

// calculateRSI calculates the Relative Strength Index
func (s *IndicatorService) calculateRSI(period int) float64 {
	// Mock RSI calculation - in real implementation, use actual price data
	return 30 + math.Mod(float64(period), 40)
}

// calculateStochRSI calculates the Stochastic RSI
func (s *IndicatorService) calculateStochRSI(period int) float64 {
	// Mock Stochastic RSI calculation
	return 20 + math.Mod(float64(period), 60)
}

// calculateHistoricRSIAvg calculates the historic RSI average
func (s *IndicatorService) calculateHistoricRSIAvg(period int) float64 {
	// Mock historic RSI average
	return 45 + math.Mod(float64(period), 10)
}

// calculateSMA calculates Simple Moving Average
func (s *IndicatorService) calculateSMA(period, window int) float64 {
	// Mock SMA calculation
	basePrice := 100.0
	return basePrice * (0.9 + math.Mod(float64(period), 0.2))
}

// calculateEMA calculates Exponential Moving Average
func (s *IndicatorService) calculateEMA(period, window int) float64 {
	// Mock EMA calculation
	basePrice := 100.0
	return basePrice * (0.95 + math.Mod(float64(period), 0.1))
}

// calculateMACD calculates MACD
func (s *IndicatorService) calculateMACD(period int) float64 {
	// Mock MACD calculation
	return (math.Mod(float64(period), 10) - 5) * 0.5
}

// calculateMACDSignal calculates MACD Signal line
func (s *IndicatorService) calculateMACDSignal(period int) float64 {
	// Mock MACD Signal calculation
	macd := s.calculateMACD(period)
	return macd + (math.Mod(float64(period), 4)-2)*0.2
}

// calculateMACDHistogram calculates MACD Histogram
func (s *IndicatorService) calculateMACDHistogram(period int) float64 {
	// Mock MACD Histogram calculation
	macd := s.calculateMACD(period)
	signal := s.calculateMACDSignal(period)
	return macd - signal
}

// CalculateEquilibriumLevel calculates the equilibrium level (50% retracement)
func (s *IndicatorService) CalculateEquilibriumLevel(high52Week, low52Week float64) float64 {
	return (high52Week + low52Week) / 2
}

// CalculatePriceToEquilibrium calculates the percentage distance from equilibrium
func (s *IndicatorService) CalculatePriceToEquilibrium(currentPrice, equilibriumLevel float64) float64 {
	return ((currentPrice - equilibriumLevel) / equilibriumLevel) * 100
}

// DetermineTrend determines the trend based on moving averages
func (s *IndicatorService) DetermineTrend(currentPrice, sma50, sma200 float64) string {
	if currentPrice > sma50 && sma50 > sma200 {
		return "bullish"
	} else if currentPrice < sma50 && sma50 < sma200 {
		return "bearish"
	}
	return "neutral"
}

// DetermineSignal determines the trading signal based on RSI and equilibrium
func (s *IndicatorService) DetermineSignal(rsi, priceToEquilibrium float64) string {
	if rsi < 40 && priceToEquilibrium < -10 {
		return "buy"
	} else if rsi > 70 && priceToEquilibrium > 10 {
		return "sell"
	}
	return "hold"
}

// DetermineVolumeProfile determines the volume profile based on volume
func (s *IndicatorService) DetermineVolumeProfile(volume int64) string {
	if volume > 50000000 {
		return "high"
	} else if volume < 10000000 {
		return "low"
	}
	return "medium"
}

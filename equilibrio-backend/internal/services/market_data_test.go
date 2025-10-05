package services

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"equilibrio-backend/internal/models"
)

// generateMockStock generates a single mock stock for testing
func generateMockStock(symbol string, index int) models.StockData {
	rand.Seed(time.Now().UnixNano() + int64(index))

	sectors := []string{
		"Technology", "Healthcare", "Financial", "Consumer", "Energy",
		"Utilities", "Real Estate", "Materials", "Industrials",
	}

	industries := []string{
		"Software", "Semiconductors", "Biotechnology", "Banks", "Insurance",
		"Retail", "Oil & Gas", "Electric Utilities", "REITs", "Chemicals",
	}

	// Generate varied RSI distribution
	var rsi float64
	rsiRand := rand.Float64()
	if rsiRand < 0.15 {
		rsi = 20 + rand.Float64()*10 // 20-30: oversold
	} else if rsiRand < 0.30 {
		rsi = 30 + rand.Float64()*10 // 30-40: approaching oversold
	} else if rsiRand < 0.70 {
		rsi = 40 + rand.Float64()*30 // 40-70: neutral
	} else if rsiRand < 0.85 {
		rsi = 70 + rand.Float64()*10 // 70-80: approaching overbought
	} else {
		rsi = 80 + rand.Float64()*10 // 80-90: overbought
	}

	price := 10 + rand.Float64()*490
	change := (rand.Float64() - 0.5) * 10
	changePercent := (change / price) * 100

	// Determine signal based on multiple factors
	var signal string
	if rsi < 35 && changePercent < 0 {
		signal = "buy"
	} else if rsi > 65 && changePercent > 0 {
		signal = "sell"
	} else {
		signal = "hold"
	}

	volumeProfiles := []string{"low", "normal", "high", "extreme"}
	trends := []string{"bullish", "bearish", "sideways"}
	equilibriumZones := []string{"support", "resistance", "neutral"}

	return models.StockData{
		Symbol:              symbol,
		Name:                "Company " + symbol,
		Sector:              sectors[rand.Intn(len(sectors))],
		Industry:            industries[rand.Intn(len(industries))],
		Price:               math.Round(price*100) / 100,
		Change:              math.Round(change*100) / 100,
		ChangePercent:       math.Round(changePercent*100) / 100,
		Volume:              int64(rand.Intn(10000000) + 1000000),
		MarketCap:           int64(rand.Intn(1000000000000) + 1000000000),
		PERatio:             math.Round((5+rand.Float64()*45)*100) / 100,
		DividendYield:       math.Round(rand.Float64()*5*100) / 100,
		Week52High:          math.Round((price*(1+rand.Float64()*0.3))*100) / 100,
		Week52Low:           math.Round((price*(1-rand.Float64()*0.3))*100) / 100,
		RSI:                 math.Round(rsi*100) / 100,
		MACD:                math.Round((rand.Float64()-0.5)*10*100) / 100,
		MovingAvg50:         math.Round((price*(1+(rand.Float64()-0.5)*0.1))*100) / 100,
		MovingAvg200:        math.Round((price*(1+(rand.Float64()-0.5)*0.2))*100) / 100,
		BollingerUpper:      math.Round((price*1.1)*100) / 100,
		BollingerLower:      math.Round((price*0.9)*100) / 100,
		ATR:                 math.Round(rand.Float64()*5*100) / 100,
		Signal:              signal,
		VolumeProfile:       volumeProfiles[rand.Intn(len(volumeProfiles))],
		Trend:               trends[rand.Intn(len(trends))],
		EquilibriumZone:     equilibriumZones[rand.Intn(len(equilibriumZones))],
		EquilibriumStrength: math.Round((0.3+rand.Float64()*0.7)*100) / 100,
		SupportLevel:        math.Round((price*0.95)*100) / 100,
		ResistanceLevel:     math.Round((price*1.05)*100) / 100,
		LastUpdated:         time.Now(),
	}
}

// generateMockChartDataForTest generates mock candlestick data for testing
func generateMockChartDataForTest(currentPrice float64, days int) []models.CandlestickData {
	data := make([]models.CandlestickData, days)
	price := currentPrice * 0.95 // Start 5% below current price
	now := time.Now()

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -(days - i - 1))

		// Random price movement
		change := (rand.Float64() - 0.5) * 0.04 // +/- 2% change
		open := price
		close := price * (1 + change)

		// Generate high and low based on open/close
		volatility := 0.015 // 1.5% intraday volatility
		high := math.Max(open, close) * (1 + rand.Float64()*volatility)
		low := math.Min(open, close) * (1 - rand.Float64()*volatility)

		data[i] = models.CandlestickData{
			Time:  date.Format("2006-01-02"),
			Open:  math.Round(open*100) / 100,
			High:  math.Round(high*100) / 100,
			Low:   math.Round(low*100) / 100,
			Close: math.Round(close*100) / 100,
		}

		price = close
	}

	return data
}

// TestGenerateMockStock tests the mock stock generation
func TestGenerateMockStock(t *testing.T) {
	stock := generateMockStock("TEST", 1)

	if stock.Symbol != "TEST" {
		t.Errorf("Expected symbol TEST, got %s", stock.Symbol)
	}

	if stock.Price <= 0 {
		t.Errorf("Price should be positive, got %f", stock.Price)
	}

	if stock.RSI < 0 || stock.RSI > 100 {
		t.Errorf("RSI should be between 0 and 100, got %f", stock.RSI)
	}
}

// TestGenerateMockChartData tests the mock chart data generation
func TestGenerateMockChartData(t *testing.T) {
	data := generateMockChartDataForTest(100.0, 30)

	if len(data) != 30 {
		t.Errorf("Expected 30 data points, got %d", len(data))
	}

	for _, candle := range data {
		if candle.High < candle.Low {
			t.Errorf("High should be greater than Low")
		}
		if candle.High < candle.Open || candle.High < candle.Close {
			t.Errorf("High should be the highest price")
		}
		if candle.Low > candle.Open || candle.Low > candle.Close {
			t.Errorf("Low should be the lowest price")
		}
	}
}

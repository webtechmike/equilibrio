package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"equilibrio-backend/internal/models"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
)

// YahooFinanceProvider implements MarketDataProvider using Yahoo Finance
type YahooFinanceProvider struct {
	// No API key needed for Yahoo Finance
}

// NewYahooFinanceProvider creates a new Yahoo Finance provider
func NewYahooFinanceProvider() *YahooFinanceProvider {
	return &YahooFinanceProvider{}
}

// GetQuote fetches real-time quote for a symbol from Yahoo Finance
func (p *YahooFinanceProvider) GetQuote(ctx context.Context, symbol string) (*models.Quote, error) {
	q, err := quote.Get(symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote for %s: %w", symbol, err)
	}

	// Calculate change and change percent
	change := q.RegularMarketPrice - q.RegularMarketPreviousClose
	changePercent := (change / q.RegularMarketPreviousClose) * 100

	return &models.Quote{
		Symbol:        q.Symbol,
		Name:          q.ShortName,
		Price:         q.RegularMarketPrice,
		Change:        change,
		ChangePercent: changePercent,
		Volume:        int64(q.RegularMarketVolume),
		MarketCap:     0, // Not directly available, would need separate API call
		PERatio:       0, // Not directly available
		DividendYield: 0, // Not directly available
		Week52High:    q.FiftyTwoWeekHigh,
		Week52Low:     q.FiftyTwoWeekLow,
		Open:          q.RegularMarketOpen,
		High:          q.RegularMarketDayHigh,
		Low:           q.RegularMarketDayLow,
		PreviousClose: q.RegularMarketPreviousClose,
		Sector:        p.getSectorName(string(q.QuoteType)),
		Industry:      p.getIndustryFromSymbol(q.Symbol),
	}, nil
}

// GetHistoricalPrices fetches historical price data from Yahoo Finance
func (p *YahooFinanceProvider) GetHistoricalPrices(ctx context.Context, symbol string, days int) ([]models.CandlestickData, error) {
	// Calculate date range
	end := time.Now()
	start := end.AddDate(0, 0, -days)

	// Convert to Yahoo Finance datetime format
	params := &chart.Params{
		Symbol:   symbol,
		Start:    datetime.New(&start),
		End:      datetime.New(&end),
		Interval: datetime.OneDay,
	}

	iter := chart.Get(params)

	var candles []models.CandlestickData
	for iter.Next() {
		bar := iter.Bar()
		
		// Convert timestamp to date string
		date := time.Unix(int64(bar.Timestamp), 0).Format("2006-01-02")

		// Convert decimal.Decimal to float64
		open, _ := bar.Open.Float64()
		high, _ := bar.High.Float64()
		low, _ := bar.Low.Float64()
		close, _ := bar.Close.Float64()

		candles = append(candles, models.CandlestickData{
			Time:  date,
			Open:  p.round(open),
			High:  p.round(high),
			Low:   p.round(low),
			Close: p.round(close),
		})
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to fetch historical data for %s: %w", symbol, err)
	}

	if len(candles) == 0 {
		return nil, fmt.Errorf("no historical data found for %s", symbol)
	}

	return candles, nil
}

// SearchSymbols searches for symbols matching criteria (basic implementation)
func (p *YahooFinanceProvider) SearchSymbols(ctx context.Context, query string) ([]string, error) {
	// Yahoo Finance doesn't have a direct search API in this library
	// For now, return empty - could be extended with web scraping or another API
	return []string{}, nil
}

// GetMarketSnapshot gets current market overview for multiple symbols
func (p *YahooFinanceProvider) GetMarketSnapshot(ctx context.Context, symbols []string) (map[string]*models.Quote, error) {
	snapshot := make(map[string]*models.Quote)

	for _, symbol := range symbols {
		quote, err := p.GetQuote(ctx, symbol)
		if err != nil {
			// Log error but continue with other symbols
			continue
		}
		snapshot[symbol] = quote
	}

	return snapshot, nil
}

// Helper methods

// round rounds a float64 to 2 decimal places
func (p *YahooFinanceProvider) round(value float64) float64 {
	return math.Round(value*100) / 100
}

// getSectorName maps Yahoo quote type to sector name
func (p *YahooFinanceProvider) getSectorName(quoteType string) string {
	// This is a simple mapping - could be enhanced
	sectorMap := map[string]string{
		"EQUITY":      "Technology",
		"ETF":         "ETF",
		"MUTUALFUND":  "Mutual Fund",
		"INDEX":       "Index",
	}

	if sector, ok := sectorMap[quoteType]; ok {
		return sector
	}
	return "Technology" // Default sector
}

// getIndustryFromSymbol provides a basic industry mapping for known symbols
func (p *YahooFinanceProvider) getIndustryFromSymbol(symbol string) string {
	industryMap := map[string]string{
		"AAPL":  "Consumer Electronics",
		"MSFT":  "Software",
		"GOOGL": "Internet Services",
		"AMZN":  "E-Commerce",
		"NVDA":  "Semiconductors",
		"TSLA":  "Automotive",
		"META":  "Social Media",
		"BRK-B": "Insurance",
		"JNJ":   "Pharmaceuticals",
		"JPM":   "Banking",
		"V":     "Financial Services",
		"PG":    "Consumer Goods",
		"MA":    "Financial Services",
		"HD":    "Retail",
		"BAC":   "Banking",
		"XOM":   "Oil & Gas",
		"CVX":   "Oil & Gas",
		"ABBV":  "Biotechnology",
		"KO":    "Beverages",
		"PFE":   "Pharmaceuticals",
	}

	if industry, ok := industryMap[symbol]; ok {
		return industry
	}
	return "General"
}

// CalculateTechnicalIndicators calculates technical indicators from historical data
func (p *YahooFinanceProvider) CalculateTechnicalIndicators(candles []models.CandlestickData) models.TechnicalIndicators {
	if len(candles) == 0 {
		return models.TechnicalIndicators{}
	}

	// Calculate RSI (14-period)
	rsi := p.calculateRSI(candles, 14)

	// Calculate SMAs
	sma50 := p.calculateSMA(candles, 50)
	sma200 := p.calculateSMA(candles, 200)

	// Calculate EMA
	ema20 := p.calculateEMA(candles, 20)

	// Calculate MACD
	macd, macdSignal, macdHist := p.calculateMACD(candles)

	return models.TechnicalIndicators{
		RSI:            rsi,
		StochRSI:       0, // TODO: Implement
		HistoricRSIAvg: rsi, // Simplified
		SMA50:          sma50,
		SMA200:         sma200,
		EMA20:          ema20,
		MACD:           macd,
		MACDSignal:     macdSignal,
		MACDHistogram:  macdHist,
	}
}

// calculateRSI calculates the Relative Strength Index
func (p *YahooFinanceProvider) calculateRSI(candles []models.CandlestickData, period int) float64 {
	if len(candles) < period+1 {
		return 50.0 // Default neutral RSI
	}

	var gains, losses float64
	recentCandles := candles[len(candles)-period-1:]

	for i := 1; i < len(recentCandles); i++ {
		change := recentCandles[i].Close - recentCandles[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses += math.Abs(change)
		}
	}

	if losses == 0 {
		return 100.0
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)
	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return p.round(rsi)
}

// calculateSMA calculates Simple Moving Average
func (p *YahooFinanceProvider) calculateSMA(candles []models.CandlestickData, period int) float64 {
	if len(candles) < period {
		period = len(candles)
	}

	recentCandles := candles[len(candles)-period:]
	var sum float64
	for _, candle := range recentCandles {
		sum += candle.Close
	}

	return p.round(sum / float64(len(recentCandles)))
}

// calculateEMA calculates Exponential Moving Average
func (p *YahooFinanceProvider) calculateEMA(candles []models.CandlestickData, period int) float64 {
	if len(candles) < period {
		return p.calculateSMA(candles, len(candles))
	}

	multiplier := 2.0 / float64(period+1)
	ema := p.calculateSMA(candles[:period], period)

	for i := period; i < len(candles); i++ {
		ema = (candles[i].Close-ema)*multiplier + ema
	}

	return p.round(ema)
}

// calculateMACD calculates Moving Average Convergence Divergence
func (p *YahooFinanceProvider) calculateMACD(candles []models.CandlestickData) (float64, float64, float64) {
	if len(candles) < 26 {
		return 0, 0, 0
	}

	ema12 := p.calculateEMA(candles, 12)
	ema26 := p.calculateEMA(candles, 26)
	macd := ema12 - ema26

	// For signal line, we'd need to calculate EMA of MACD values
	// Simplified: use a fraction of MACD
	signal := macd * 0.9
	histogram := macd - signal

	return p.round(macd), p.round(signal), p.round(histogram)
}


package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"equilibrio-backend/internal/models"
)

// FinnhubProvider implements MarketDataProvider using Finnhub API
type FinnhubProvider struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// NewFinnhubProvider creates a new Finnhub provider
func NewFinnhubProvider(apiKey string) *FinnhubProvider {
	return &FinnhubProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: "https://finnhub.io/api/v1",
	}
}

// FinnhubQuoteResponse represents Finnhub quote API response
type FinnhubQuoteResponse struct {
	C  float64 `json:"c"`  // Current price
	D  float64 `json:"d"`  // Change
	DP float64 `json:"dp"` // Percent change
	H  float64 `json:"h"`  // High price of the day
	L  float64 `json:"l"`  // Low price of the day
	O  float64 `json:"o"`  // Open price of the day
	PC float64 `json:"pc"` // Previous close price
	T  int64   `json:"t"`  // Timestamp
}

// FinnhubCompanyProfileResponse represents company profile API response
type FinnhubCompanyProfileResponse struct {
	Country          string  `json:"country"`
	Currency         string  `json:"currency"`
	Exchange         string  `json:"exchange"`
	FinnhubIndustry  string  `json:"finnhubIndustry"`
	MarketCap        float64 `json:"marketCapitalization"`
	Name             string  `json:"name"`
	ShareOutstanding float64 `json:"shareOutstanding"`
	Ticker           string  `json:"ticker"`
	WebURL           string  `json:"weburl"`
}

// FinnhubCandleResponse represents historical candle data
type FinnhubCandleResponse struct {
	C []float64 `json:"c"` // Close prices
	H []float64 `json:"h"` // High prices
	L []float64 `json:"l"` // Low prices
	O []float64 `json:"o"` // Open prices
	S string    `json:"s"` // Status
	T []int64   `json:"t"` // Timestamps
	V []int64   `json:"v"` // Volumes
}

// GetQuote fetches real-time quote from Finnhub
func (p *FinnhubProvider) GetQuote(ctx context.Context, symbol string) (*models.Quote, error) {
	// Get quote data
	quoteURL := fmt.Sprintf("%s/quote?symbol=%s&token=%s", p.baseURL, symbol, p.apiKey)

	quote, err := p.makeRequest(ctx, quoteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}

	var quoteResp FinnhubQuoteResponse
	if err := json.Unmarshal(quote, &quoteResp); err != nil {
		return nil, fmt.Errorf("failed to parse quote response: %w", err)
	}

	// Get company profile for additional data
	profileURL := fmt.Sprintf("%s/stock/profile2?symbol=%s&token=%s", p.baseURL, symbol, p.apiKey)

	profile, err := p.makeRequest(ctx, profileURL)
	if err != nil {
		// If profile fails, continue with basic quote data
		fmt.Printf("Warning: failed to fetch profile for %s: %v\n", symbol, err)
	}

	var profileResp FinnhubCompanyProfileResponse
	if err := json.Unmarshal(profile, &profileResp); err != nil {
		// If profile parsing fails, continue with basic quote data
		fmt.Printf("Warning: failed to parse profile for %s: %v\n", symbol, err)
	}

	// Calculate 52-week high/low (not available in quote, use current price as approximation)
	week52High := quoteResp.C * 1.2 // Approximate
	week52Low := quoteResp.C * 0.8  // Approximate

	return &models.Quote{
		Symbol:        symbol,
		Name:          profileResp.Name,
		Price:         quoteResp.C,
		Change:        quoteResp.D,
		ChangePercent: quoteResp.DP,
		Volume:        0, // Not available in quote endpoint
		MarketCap:     int64(profileResp.MarketCap),
		PERatio:       0, // Not available
		DividendYield: 0, // Not available
		Week52High:    week52High,
		Week52Low:     week52Low,
		Open:          quoteResp.O,
		High:          quoteResp.H,
		Low:           quoteResp.L,
		PreviousClose: quoteResp.PC,
		Sector:        p.mapIndustryToSector(profileResp.FinnhubIndustry),
		Industry:      profileResp.FinnhubIndustry,
	}, nil
}

// GetHistoricalPrices fetches historical price data from Finnhub
func (p *FinnhubProvider) GetHistoricalPrices(ctx context.Context, symbol string, days int) ([]models.CandlestickData, error) {
	// Calculate timestamp range
	end := time.Now().Unix()
	start := time.Now().AddDate(0, 0, -days).Unix()

	url := fmt.Sprintf("%s/stock/candle?symbol=%s&resolution=D&from=%d&to=%d&token=%s",
		p.baseURL, symbol, start, end, p.apiKey)

	body, err := p.makeRequest(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical data: %w", err)
	}

	var candleResp FinnhubCandleResponse
	if err := json.Unmarshal(body, &candleResp); err != nil {
		return nil, fmt.Errorf("failed to parse candle response: %w", err)
	}

	if candleResp.S != "ok" {
		return nil, fmt.Errorf("API returned error status: %s", candleResp.S)
	}

	if len(candleResp.T) == 0 {
		return nil, fmt.Errorf("no historical data found for %s", symbol)
	}

	var candles []models.CandlestickData
	for i := range candleResp.T {
		date := time.Unix(candleResp.T[i], 0).Format("2006-01-02")

		candles = append(candles, models.CandlestickData{
			Time:  date,
			Open:  p.round(candleResp.O[i]),
			High:  p.round(candleResp.H[i]),
			Low:   p.round(candleResp.L[i]),
			Close: p.round(candleResp.C[i]),
		})
	}

	return candles, nil
}

// SearchSymbols searches for symbols (basic implementation)
func (p *FinnhubProvider) SearchSymbols(ctx context.Context, query string) ([]string, error) {
	// Finnhub has a symbol search endpoint, but for now return empty
	// Could be implemented with: /v1/search?q={query}&token={token}
	return []string{}, nil
}

// GetMarketSnapshot gets current market overview for multiple symbols
func (p *FinnhubProvider) GetMarketSnapshot(ctx context.Context, symbols []string) (map[string]*models.Quote, error) {
	snapshot := make(map[string]*models.Quote)

	for _, symbol := range symbols {
		quote, err := p.GetQuote(ctx, symbol)
		if err != nil {
			// Log error but continue with other symbols
			fmt.Printf("Failed to fetch %s: %v\n", symbol, err)
			continue
		}
		snapshot[symbol] = quote
	}

	return snapshot, nil
}

// CalculateTechnicalIndicators calculates technical indicators from historical data
func (p *FinnhubProvider) CalculateTechnicalIndicators(candles []models.CandlestickData) models.TechnicalIndicators {
	if len(candles) == 0 {
		return models.TechnicalIndicators{}
	}

	rsi := p.calculateRSI(candles, 14)
	sma50 := p.calculateSMA(candles, 50)
	sma200 := p.calculateSMA(candles, 200)
	ema20 := p.calculateEMA(candles, 20)
	macd, macdSignal, macdHist := p.calculateMACD(candles)

	return models.TechnicalIndicators{
		RSI:            rsi,
		StochRSI:       0, // TODO: Implement
		HistoricRSIAvg: rsi,
		SMA50:          sma50,
		SMA200:         sma200,
		EMA20:          ema20,
		MACD:           macd,
		MACDSignal:     macdSignal,
		MACDHistogram:  macdHist,
	}
}

// Helper methods

// makeRequest makes an HTTP request to Finnhub API
func (p *FinnhubProvider) makeRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Equilibrio-Backend/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// round rounds to 2 decimal places
func (p *FinnhubProvider) round(value float64) float64 {
	return math.Round(value*100) / 100
}

// mapIndustryToSector maps Finnhub industry to sector
func (p *FinnhubProvider) mapIndustryToSector(industry string) string {
	sectorMap := map[string]string{
		"Technology":             "Technology",
		"Healthcare":             "Healthcare",
		"Financial Services":     "Financial",
		"Consumer Cyclical":      "Consumer Cyclical",
		"Energy":                 "Energy",
		"Industrials":            "Industrials",
		"Consumer Defensive":     "Consumer Defensive",
		"Real Estate":            "Real Estate",
		"Communication Services": "Communication Services",
		"Utilities":              "Utilities",
		"Basic Materials":        "Basic Materials",
	}

	if sector, ok := sectorMap[industry]; ok {
		return sector
	}
	return "Technology" // Default sector
}

// calculateRSI calculates the Relative Strength Index
func (p *FinnhubProvider) calculateRSI(candles []models.CandlestickData, period int) float64 {
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
func (p *FinnhubProvider) calculateSMA(candles []models.CandlestickData, period int) float64 {
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
func (p *FinnhubProvider) calculateEMA(candles []models.CandlestickData, period int) float64 {
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
func (p *FinnhubProvider) calculateMACD(candles []models.CandlestickData) (float64, float64, float64) {
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

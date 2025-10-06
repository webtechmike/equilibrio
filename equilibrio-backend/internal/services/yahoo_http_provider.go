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

// YahooHTTPProvider uses direct HTTP requests to Yahoo Finance API
type YahooHTTPProvider struct {
	client  *http.Client
	baseURL string
}

// NewYahooHTTPProvider creates a new HTTP-based Yahoo Finance provider
func NewYahooHTTPProvider() *YahooHTTPProvider {
	return &YahooHTTPProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://query1.finance.yahoo.com/v8/finance",
	}
}

// YahooQuoteResponse represents Yahoo Finance API response
type YahooQuoteResponse struct {
	QuoteResponse struct {
		Result []struct {
			Symbol                     string  `json:"symbol"`
			ShortName                  string  `json:"shortName"`
			RegularMarketPrice         float64 `json:"regularMarketPrice"`
			RegularMarketChange        float64 `json:"regularMarketChange"`
			RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
			RegularMarketVolume        int64   `json:"regularMarketVolume"`
			MarketCap                  int64   `json:"marketCap"`
			FiftyTwoWeekHigh           float64 `json:"fiftyTwoWeekHigh"`
			FiftyTwoWeekLow            float64 `json:"fiftyTwoWeekLow"`
			RegularMarketOpen          float64 `json:"regularMarketOpen"`
			RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
			RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
			RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
		} `json:"result"`
	} `json:"quoteResponse"`
}

// YahooChartResponse represents historical data response
type YahooChartResponse struct {
	Chart struct {
		Result []struct {
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

// GetQuote fetches real-time quote via HTTP
func (p *YahooHTTPProvider) GetQuote(ctx context.Context, symbol string) (*models.Quote, error) {
	url := fmt.Sprintf("%s/quote?symbols=%s", p.baseURL, symbol)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
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

	var yahooResp YahooQuoteResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(yahooResp.QuoteResponse.Result) == 0 {
		return nil, fmt.Errorf("no data returned for symbol %s", symbol)
	}

	result := yahooResp.QuoteResponse.Result[0]

	return &models.Quote{
		Symbol:        result.Symbol,
		Name:          result.ShortName,
		Price:         result.RegularMarketPrice,
		Change:        result.RegularMarketChange,
		ChangePercent: result.RegularMarketChangePercent,
		Volume:        result.RegularMarketVolume,
		MarketCap:     result.MarketCap,
		Week52High:    result.FiftyTwoWeekHigh,
		Week52Low:     result.FiftyTwoWeekLow,
		Open:          result.RegularMarketOpen,
		High:          result.RegularMarketDayHigh,
		Low:           result.RegularMarketDayLow,
		PreviousClose: result.RegularMarketPreviousClose,
		Sector:        p.getSectorForSymbol(symbol),
		Industry:      p.getIndustryForSymbol(symbol),
	}, nil
}

// GetHistoricalPrices fetches historical data via HTTP
func (p *YahooHTTPProvider) GetHistoricalPrices(ctx context.Context, symbol string, days int) ([]models.CandlestickData, error) {
	end := time.Now().Unix()
	start := time.Now().AddDate(0, 0, -days).Unix()

	url := fmt.Sprintf("%s/chart/%s?period1=%d&period2=%d&interval=1d",
		p.baseURL, symbol, start, end)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch historical data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chartResp YahooChartResponse
	if err := json.Unmarshal(body, &chartResp); err != nil {
		return nil, fmt.Errorf("failed to parse chart response: %w", err)
	}

	if len(chartResp.Chart.Result) == 0 {
		return nil, fmt.Errorf("no chart data returned")
	}

	result := chartResp.Chart.Result[0]
	if len(result.Indicators.Quote) == 0 {
		return nil, fmt.Errorf("no quote data in chart")
	}

	timestamps := result.Timestamp
	quote := result.Indicators.Quote[0]

	var candles []models.CandlestickData
	for i := range timestamps {
		date := time.Unix(timestamps[i], 0).Format("2006-01-02")

		candles = append(candles, models.CandlestickData{
			Time:  date,
			Open:  p.round(quote.Open[i]),
			High:  p.round(quote.High[i]),
			Low:   p.round(quote.Low[i]),
			Close: p.round(quote.Close[i]),
		})
	}

	return candles, nil
}

// round rounds to 2 decimal places
func (p *YahooHTTPProvider) round(value float64) float64 {
	return math.Round(value*100) / 100
}

// getSectorForSymbol maps symbols to sectors
func (p *YahooHTTPProvider) getSectorForSymbol(symbol string) string {
	sectorMap := map[string]string{
		"AAPL": "Technology", "MSFT": "Technology", "GOOGL": "Communication Services",
		"AMZN": "Consumer Cyclical", "NVDA": "Technology", "TSLA": "Consumer Cyclical",
		"META": "Communication Services", "BRK-B": "Financial", "JNJ": "Healthcare",
		"JPM": "Financial", "V": "Financial", "PG": "Consumer Defensive",
		"MA": "Financial", "HD": "Consumer Cyclical", "BAC": "Financial",
		"XOM": "Energy", "CVX": "Energy", "ABBV": "Healthcare",
		"KO": "Consumer Defensive", "PFE": "Healthcare",
	}
	if sector, ok := sectorMap[symbol]; ok {
		return sector
	}
	return "Technology"
}

// getIndustryForSymbol maps symbols to industries
func (p *YahooHTTPProvider) getIndustryForSymbol(symbol string) string {
	industryMap := map[string]string{
		"AAPL": "Consumer Electronics", "MSFT": "Software", "GOOGL": "Internet Services",
		"AMZN": "E-Commerce", "NVDA": "Semiconductors", "TSLA": "Automotive",
		"META": "Social Media", "BRK-B": "Insurance", "JNJ": "Pharmaceuticals",
		"JPM": "Banking", "V": "Financial Services", "PG": "Consumer Goods",
		"MA": "Financial Services", "HD": "Retail", "BAC": "Banking",
		"XOM": "Oil & Gas", "CVX": "Oil & Gas", "ABBV": "Biotechnology",
		"KO": "Beverages", "PFE": "Pharmaceuticals",
	}
	if industry, ok := industryMap[symbol]; ok {
		return industry
	}
	return "General"
}

// CalculateTechnicalIndicators calculates indicators from historical data
func (p *YahooHTTPProvider) CalculateTechnicalIndicators(candles []models.CandlestickData) models.TechnicalIndicators {
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
		StochRSI:       0,
		HistoricRSIAvg: rsi,
		SMA50:          sma50,
		SMA200:         sma200,
		EMA20:          ema20,
		MACD:           macd,
		MACDSignal:     macdSignal,
		MACDHistogram:  macdHist,
	}
}

func (p *YahooHTTPProvider) calculateRSI(candles []models.CandlestickData, period int) float64 {
	if len(candles) < period+1 {
		return 50.0
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

func (p *YahooHTTPProvider) calculateSMA(candles []models.CandlestickData, period int) float64 {
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

func (p *YahooHTTPProvider) calculateEMA(candles []models.CandlestickData, period int) float64 {
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

func (p *YahooHTTPProvider) calculateMACD(candles []models.CandlestickData) (float64, float64, float64) {
	if len(candles) < 26 {
		return 0, 0, 0
	}

	ema12 := p.calculateEMA(candles, 12)
	ema26 := p.calculateEMA(candles, 26)
	macd := ema12 - ema26
	signal := macd * 0.9
	histogram := macd - signal

	return p.round(macd), p.round(signal), p.round(histogram)
}

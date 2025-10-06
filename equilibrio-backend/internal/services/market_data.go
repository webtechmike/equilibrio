package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"equilibrio-backend/internal/config"
	"equilibrio-backend/internal/models"

	"github.com/redis/go-redis/v9"
)

type MarketDataService struct {
	config        *config.Config
	cache         *redis.Client
	cacheStrategy *MarketCacheStrategy
	yahooProvider *YahooHTTPProvider
	indicators    *EquilibriumCalculator
}

func NewMarketDataService(cfg *config.Config) *MarketDataService {
	// Initialize Redis client
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		// Fallback to default Redis connection
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	rdb := redis.NewClient(opt)

	return &MarketDataService{
		config:        cfg,
		cache:         rdb,
		cacheStrategy: NewMarketCacheStrategy(rdb),
		yahooProvider: NewYahooHTTPProvider(),       // Use HTTP provider
		indicators:    NewEquilibriumCalculator(20), // 20-day lookback
	}
}

// GetStocks retrieves and filters stocks based on the request
func (s *MarketDataService) GetStocks(req models.StockListRequest) ([]models.StockData, int, error) {
	ctx := context.Background()

	// Log market status
	marketStatus := s.cacheStrategy.GetMarketStatus()
	cacheTTL := s.cacheStrategy.GetCacheTTL()
	fmt.Printf("Market Status: %s | Cache TTL: %v\n", marketStatus, cacheTTL)

	// Try to get from daily snapshot first (if market is closed)
	var stocks []models.StockData
	var err2 error

	if !s.cacheStrategy.IsMarketOpen() {
		// Market closed - try daily snapshot
		snapshot, err := s.cacheStrategy.GetDailySnapshot(ctx)
		if err == nil && len(snapshot) > 0 {
			fmt.Println("Using daily snapshot (market closed)")
			stocks = snapshot
		}
	}

	// If no snapshot or market is open, get fresh data
	if len(stocks) == 0 {
		fmt.Printf("USE_MOCK_DATA config: %v\n", s.config.UseMockData)

		if s.config.UseMockData {
			fmt.Println("Using MOCK data")
			stocks = s.generateMockStockData()
		} else {
			fmt.Println("Fetching REAL data from Yahoo Finance...")
			stocks, err2 = s.getRealStockData()
			if err2 != nil {
				// Fallback to mock data if real data fails
				fmt.Printf("Failed to fetch real data, using mock: %v\n", err2)
				stocks = s.generateMockStockData()
			} else {
				fmt.Printf("Successfully fetched %d real stocks\n", len(stocks))

				// Cache daily snapshot if market just closed or during trading hours
				if err := s.cacheStrategy.CacheDailySnapshot(ctx, stocks); err == nil {
					fmt.Println("Cached daily snapshot for reuse")
				}
			}
		}
	}

	// Create filter from request
	filter := models.StockFilter{
		SearchTerm:      req.SearchTerm,
		Sectors:         req.Sectors,
		RSIMin:          req.RSIMin,
		RSIMax:          req.RSIMax,
		PriceMin:        req.PriceMin,
		PriceMax:        req.PriceMax,
		VolumeProfile:   req.VolumeProfile,
		Signals:         req.Signals,
		Trend:           req.Trend,
		EquilibriumZone: req.EquilibriumZone,
	}

	// Apply filters
	filteredStocks := s.applyFilters(stocks, filter)

	// Apply sorting
	sortedStocks := s.applySorting(filteredStocks, req.SortField, req.SortOrder)

	// Apply pagination
	total := len(sortedStocks)
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	if start >= total {
		return []models.StockData{}, total, nil
	}

	if end > total {
		end = total
	}

	paginatedStocks := sortedStocks[start:end]

	// Cache the result with smart TTL based on market hours
	cacheData := struct {
		Stocks []models.StockData `json:"stocks"`
		Total  int                `json:"total"`
	}{
		Stocks: paginatedStocks,
		Total:  total,
	}

	// Use market-aware caching
	if err := s.cacheStrategy.CacheStockData(ctx, fmt.Sprintf("stocks:%s", s.generateCacheKey(req)), cacheData); err != nil {
		fmt.Printf("Failed to cache results: %v\n", err)
	}

	return paginatedStocks, total, nil
}

// GetStock retrieves a single stock by symbol
func (s *MarketDataService) GetStock(symbol string) (*models.StockData, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("stock:%s", strings.ToUpper(symbol))
	cached, err := s.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var stock models.StockData
		if json.Unmarshal([]byte(cached), &stock) == nil {
			return &stock, nil
		}
	}

	// Generate mock data for the symbol
	stocks := s.generateMockStockData()
	for _, stock := range stocks {
		if strings.ToUpper(stock.Symbol) == strings.ToUpper(symbol) {
			// Cache the result
			if data, err := json.Marshal(stock); err == nil {
				s.cache.Set(context.Background(), cacheKey, data, 30*time.Second)
			}
			return &stock, nil
		}
	}

	return nil, fmt.Errorf("stock not found: %s", symbol)
}

// SearchStock searches for any ticker symbol and fetches from Yahoo Finance
func (s *MarketDataService) SearchStock(symbol string) (*models.StockData, error) {
	ctx := context.Background()
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	
	// Check cache first
	cacheKey := fmt.Sprintf("stock:%s", symbol)
	cached, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var stock models.StockData
		if json.Unmarshal([]byte(cached), &stock) == nil {
			fmt.Printf("Cache hit for symbol: %s\n", symbol)
			return &stock, nil
		}
	}
	
	fmt.Printf("Searching for symbol: %s\n", symbol)
	
	// If using mock data, check if symbol exists in mock list
	if s.config.UseMockData {
		stocks := s.generateMockStockData()
		for _, stock := range stocks {
			if stock.Symbol == symbol {
				// Cache the result with market-aware TTL
				if data, err := json.Marshal(stock); err == nil {
					ttl := s.cacheStrategy.GetCacheTTL()
					s.cache.Set(ctx, cacheKey, data, ttl)
				}
				return &stock, nil
			}
		}
		return nil, fmt.Errorf("symbol not found in mock data: %s", symbol)
	}
	
	// Fetch from Yahoo Finance
	quote, err := s.yahooProvider.GetQuote(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch symbol %s: %w", symbol, err)
	}
	
	// Fetch historical data for indicators (90 days)
	historical, err := s.yahooProvider.GetHistoricalPrices(ctx, symbol, 90)
	if err != nil || len(historical) == 0 {
		return nil, fmt.Errorf("failed to fetch historical data for %s: %w", symbol, err)
	}
	
	// Calculate technical indicators
	indicators := s.yahooProvider.CalculateTechnicalIndicators(historical)
	
	// Calculate equilibrium
	equilibrium := s.indicators.CalculateEquilibrium(historical, quote.Price)
	
	// Determine trend based on moving averages
	trend := "neutral"
	if quote.Price > indicators.SMA50 && indicators.SMA50 > indicators.SMA200 {
		trend = "bullish"
	} else if quote.Price < indicators.SMA50 && indicators.SMA50 < indicators.SMA200 {
		trend = "bearish"
	}
	
	// Determine signal based on RSI and equilibrium
	signal := "hold"
	if indicators.RSI < 30 {
		signal = "buy"
	} else if indicators.RSI > 70 {
		signal = "sell"
	} else if equilibrium.Support > 0 {
		priceToEquilibrium := ((quote.Price - equilibrium.Support) / equilibrium.Support) * 100
		if priceToEquilibrium < -15 {
			signal = "buy"
		} else if priceToEquilibrium > 15 {
			signal = "sell"
		}
	}
	
	// Volume profile
	volumeProfile := "medium"
	if quote.Volume > 50000000 {
		volumeProfile = "high"
	} else if quote.Volume < 10000000 {
		volumeProfile = "low"
	}
	
	// Create stock data
	stock := models.StockData{
		Symbol:                 quote.Symbol,
		Name:                   quote.Name,
		Price:                  quote.Price,
		Change:                 quote.Change,
		ChangePercent:          quote.ChangePercent,
		Volume:                 quote.Volume,
		Sector:                 quote.Sector,
		Industry:               quote.Industry,
		MarketCap:              float64(quote.MarketCap),
		RSI:                    indicators.RSI,
		StochRSI:               indicators.StochRSI,
		HistoricRSIAvg:         indicators.HistoricRSIAvg,
		SMA50:                  indicators.SMA50,
		SMA200:                 indicators.SMA200,
		EMA20:                  indicators.EMA20,
		MACD:                   indicators.MACD,
		MACDSignal:             indicators.MACDSignal,
		MACDHistogram:          indicators.MACDHistogram,
		EquilibriumLevel:       (equilibrium.Support + equilibrium.Resistance) / 2,
		PriceToEquilibrium:     ((quote.Price - equilibrium.Support) / equilibrium.Support) * 100,
		Trend:                  trend,
		Signal:                 signal,
		VolumeProfile:          volumeProfile,
		DistanceFrom52WeekHigh: ((quote.Price - quote.Week52High) / quote.Week52High) * 100,
		DistanceFrom52WeekLow:  ((quote.Price - quote.Week52Low) / quote.Week52Low) * 100,
		LastUpdated:            time.Now(),
	}
	
	// Cache the result with market-aware TTL
	if data, err := json.Marshal(stock); err == nil {
		ttl := s.cacheStrategy.GetCacheTTL()
		s.cache.Set(ctx, cacheKey, data, ttl)
		fmt.Printf("Cached %s for %v\n", symbol, ttl)
	}
	
	return &stock, nil
}

// GetSectors returns all available sectors
func (s *MarketDataService) GetSectors() ([]string, error) {
	sectors := []string{
		"Technology", "Healthcare", "Financial", "Consumer Cyclical",
		"Energy", "Industrials", "Consumer Defensive", "Real Estate",
		"Communication Services", "Utilities", "Basic Materials",
	}
	return sectors, nil
}

// GetStockChart returns candlestick chart data for a stock (default 90 days)
func (s *MarketDataService) GetStockChart(symbol string) (*models.ChartDataResponse, error) {
	return s.GetStockChartWithDays(symbol, 90)
}

// GetStockChartWithDays returns candlestick chart data for a stock with specified days
func (s *MarketDataService) GetStockChartWithDays(symbol string, days int) (*models.ChartDataResponse, error) {
	var data []models.CandlestickData
	var err error

	if s.config.UseMockData {
		// Get the stock to get its current price
		stock, err2 := s.GetStock(symbol)
		if err2 != nil {
			return nil, err2
		}
		// Generate mock candlestick data for specified days
		data = s.generateMockChartData(stock.Price, days)
	} else {
		// Fetch real historical data from Yahoo Finance
		data, err = s.yahooProvider.GetHistoricalPrices(context.Background(), symbol, days)
		if err != nil {
			// Fallback to mock data
			fmt.Printf("Failed to fetch real chart data, using mock: %v\n", err)
			stock, err2 := s.GetStock(symbol)
			if err2 != nil {
				return nil, err2
			}
			data = s.generateMockChartData(stock.Price, days)
		}
	}

	response := &models.ChartDataResponse{
		Symbol: symbol,
		Data:   data,
	}

	return response, nil
}

// generateMockChartData generates realistic candlestick data that ends at currentPrice
func (s *MarketDataService) generateMockChartData(currentPrice float64, days int) []models.CandlestickData {
	data := make([]models.CandlestickData, days)
	now := time.Now()

	// Use a consistent seed based on the current price to ensure repeatability
	// This ensures the same stock always generates the same historical pattern
	seed := int64(currentPrice * 1000)
	rng := rand.New(rand.NewSource(seed))

	// First pass: generate random walk backwards, then we'll normalize
	prices := make([]float64, days+1)
	prices[days] = currentPrice // End at current price

	// Work backwards from current price
	for i := days - 1; i >= 0; i-- {
		// Random daily change: +/- 2%
		change := (rng.Float64() - 0.5) * 0.04
		prices[i] = prices[i+1] / (1 + change) // Reverse the change
	}

	// Second pass: create candlesticks
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -(days - i - 1))

		open := prices[i]
		close := prices[i+1]

		// High and low based on volatility
		volatility := 0.015 // 1.5% intraday volatility
		high := math.Max(open, close) * (1 + rng.Float64()*volatility)
		low := math.Min(open, close) * (1 - rng.Float64()*volatility)

		data[i] = models.CandlestickData{
			Time:  date.Format("2006-01-02"),
			Open:  math.Round(open*100) / 100,
			High:  math.Round(high*100) / 100,
			Low:   math.Round(low*100) / 100,
			Close: math.Round(close*100) / 100,
		}
	}

	return data
}

// RefreshAllData refreshes all stock data
func (s *MarketDataService) RefreshAllData() error {
	// Clear cache
	s.cache.FlushDB(context.Background())

	// In a real implementation, this would fetch fresh data from APIs
	// For now, we'll just return success
	return nil
}

// generateMockStockData creates mock stock data (replace with real API integration)
func (s *MarketDataService) generateMockStockData() []models.StockData {
	tickers := []struct {
		symbol string
		name   string
		sector string
	}{
		{"AAPL", "Apple Inc.", "Technology"},
		{"MSFT", "Microsoft Corp.", "Technology"},
		{"GOOGL", "Alphabet Inc.", "Communication Services"},
		{"AMZN", "Amazon.com Inc.", "Consumer Cyclical"},
		{"NVDA", "NVIDIA Corp.", "Technology"},
		{"TSLA", "Tesla Inc.", "Consumer Cyclical"},
		{"META", "Meta Platforms", "Communication Services"},
		{"BRK.B", "Berkshire Hathaway", "Financial"},
		{"JNJ", "Johnson & Johnson", "Healthcare"},
		{"JPM", "JPMorgan Chase", "Financial"},
		{"V", "Visa Inc.", "Financial"},
		{"PG", "Procter & Gamble", "Consumer Defensive"},
		{"MA", "Mastercard Inc.", "Financial"},
		{"HD", "Home Depot", "Consumer Cyclical"},
		{"BAC", "Bank of America", "Financial"},
		{"XOM", "Exxon Mobil", "Energy"},
		{"CVX", "Chevron Corp.", "Energy"},
		{"ABBV", "AbbVie Inc.", "Healthcare"},
		{"KO", "Coca-Cola Co.", "Consumer Defensive"},
		{"PFE", "Pfizer Inc.", "Healthcare"},
	}

	var stocks []models.StockData
	for i, ticker := range tickers {
		// Use consistent seed per symbol for repeatable data
		symbolSeed := int64(len(ticker.symbol)*100 + i)
		rng := rand.New(rand.NewSource(symbolSeed))

		basePrice := rng.Float64()*500 + 50
		changePercent := (rng.Float64() - 0.5) * 10
		rsi := rng.Float64() * 100
		sma50 := basePrice * (0.9 + rng.Float64()*0.2)
		sma200 := basePrice * (0.85 + rng.Float64()*0.3)
		high52Week := basePrice * (1 + rng.Float64()*0.3)
		low52Week := basePrice * (0.7 + rng.Float64()*0.2)

		// Calculate equilibrium (50% retracement from low to high)
		equilibriumLevel := (high52Week + low52Week) / 2
		priceToEquilibrium := ((basePrice - equilibriumLevel) / equilibriumLevel) * 100

		macd := (rng.Float64() - 0.5) * 5
		macdSignal := macd + (rng.Float64()-0.5)*2

		// Determine trend based on moving averages
		var trend string = "neutral"
		if basePrice > sma50 && sma50 > sma200 {
			trend = "bullish"
		} else if basePrice < sma50 && sma50 < sma200 {
			trend = "bearish"
		}

		// Determine signal based on RSI and equilibrium (more varied)
		var signal string = "hold"
		if rsi < 30 {
			signal = "buy" // Oversold
		} else if rsi > 70 {
			signal = "sell" // Overbought
		} else if priceToEquilibrium < -15 {
			signal = "buy" // Strong discount
		} else if priceToEquilibrium > 15 {
			signal = "sell" // Strong premium
		} else {
			// Random distribution for more variety
			randSignal := rng.Float64()
			if randSignal < 0.2 {
				signal = "buy"
			} else if randSignal < 0.4 {
				signal = "sell"
			} else {
				signal = "hold"
			}
		}

		// Volume profile based on volume
		avgVolume := rng.Float64() * 100000000
		var volumeProfile string = "medium"
		if avgVolume > 50000000 {
			volumeProfile = "high"
		} else if avgVolume < 10000000 {
			volumeProfile = "low"
		}

		stock := models.StockData{
			Symbol:                 ticker.symbol,
			Name:                   ticker.name,
			Price:                  basePrice,
			Change:                 basePrice * (changePercent / 100),
			ChangePercent:          changePercent,
			Volume:                 int64(avgVolume),
			Sector:                 ticker.sector,
			Industry:               s.getIndustryForSectorWithSeed(ticker.sector, symbolSeed),
			MarketCap:              basePrice * (rng.Float64()*1000000000 + 100000000),
			RSI:                    rsi,
			StochRSI:               rng.Float64() * 100,
			HistoricRSIAvg:         50 + (rng.Float64()-0.5)*20,
			SMA50:                  sma50,
			SMA200:                 sma200,
			EMA20:                  basePrice * (0.95 + rng.Float64()*0.1),
			MACD:                   macd,
			MACDSignal:             macdSignal,
			MACDHistogram:          macd - macdSignal,
			EquilibriumLevel:       equilibriumLevel,
			PriceToEquilibrium:     priceToEquilibrium,
			Trend:                  trend,
			Signal:                 signal,
			VolumeProfile:          volumeProfile,
			DistanceFrom52WeekHigh: ((basePrice - high52Week) / high52Week) * 100,
			DistanceFrom52WeekLow:  ((basePrice - low52Week) / low52Week) * 100,
			LastUpdated:            time.Now(),
		}

		stocks = append(stocks, stock)
	}

	return stocks
}

// getIndustryForSector returns a random industry for a given sector
func (s *MarketDataService) getIndustryForSector(sector string) string {
	return s.getIndustryForSectorWithSeed(sector, time.Now().UnixNano())
}

// getIndustryForSectorWithSeed returns an industry for a given sector using a seed for consistency
func (s *MarketDataService) getIndustryForSectorWithSeed(sector string, seed int64) string {
	industries := map[string][]string{
		"Technology":             {"Software", "Semiconductors", "Hardware", "IT Services"},
		"Healthcare":             {"Biotechnology", "Pharmaceuticals", "Medical Devices", "Healthcare Plans"},
		"Financial":              {"Banks", "Insurance", "Asset Management", "Capital Markets"},
		"Consumer Cyclical":      {"Retail", "Automotive", "Apparel", "Restaurants"},
		"Energy":                 {"Oil & Gas", "Renewable Energy", "Utilities"},
		"Industrials":            {"Aerospace", "Construction", "Manufacturing", "Transportation"},
		"Consumer Defensive":     {"Food Products", "Beverages", "Household Products"},
		"Real Estate":            {"REITs", "Real Estate Services", "Development"},
		"Communication Services": {"Telecom", "Media", "Entertainment"},
		"Utilities":              {"Electric", "Gas", "Water"},
		"Basic Materials":        {"Chemicals", "Metals & Mining", "Paper & Forest Products"},
	}

	if sectorIndustries, exists := industries[sector]; exists {
		rng := rand.New(rand.NewSource(seed))
		return sectorIndustries[rng.Intn(len(sectorIndustries))]
	}
	return "General"
}

// applyFilters applies the filter criteria to the stock list
// All filters are combined with AND logic - stock must pass all active filters
func (s *MarketDataService) applyFilters(stocks []models.StockData, filter models.StockFilter) []models.StockData {
	var filtered []models.StockData

	for _, stock := range stocks {
		if !s.matchesFilters(stock, filter) {
			continue
		}
		filtered = append(filtered, stock)
	}

	return filtered
}

// matchesFilters checks if a stock matches all active filter criteria
func (s *MarketDataService) matchesFilters(stock models.StockData, filter models.StockFilter) bool {
	// Search term filter (if provided)
	if filter.SearchTerm != "" {
		if !s.matchesSearchTerm(stock, filter.SearchTerm) {
			return false
		}
	}

	// Sector filter (if any sectors selected)
	if len(filter.Sectors) > 0 {
		if !s.containsString(filter.Sectors, stock.Sector) {
			return false
		}
	}

	// RSI range filter (always applied with defaults 0-100)
	if stock.RSI < filter.RSIMin || stock.RSI > filter.RSIMax {
		return false
	}

	// Price range filter (always applied with defaults)
	if stock.Price < filter.PriceMin || stock.Price > filter.PriceMax {
		return false
	}

	// Volume profile filter (if any selected)
	if len(filter.VolumeProfile) > 0 {
		if !s.containsString(filter.VolumeProfile, stock.VolumeProfile) {
			return false
		}
	}

	// Signal filter (if any selected)
	if len(filter.Signals) > 0 {
		if !s.containsString(filter.Signals, stock.Signal) {
			return false
		}
	}

	// Trend filter (if any selected)
	if len(filter.Trend) > 0 {
		if !s.containsString(filter.Trend, stock.Trend) {
			return false
		}
	}

	// Equilibrium zone filter (if any selected)
	if len(filter.EquilibriumZone) > 0 {
		if !s.matchesEquilibriumZone(stock, filter.EquilibriumZone) {
			return false
		}
	}

	// All filters passed
	return true
}

// matchesSearchTerm checks if stock matches search term in symbol or name
func (s *MarketDataService) matchesSearchTerm(stock models.StockData, searchTerm string) bool {
	searchLower := strings.ToLower(searchTerm)
	return strings.Contains(strings.ToLower(stock.Symbol), searchLower) ||
		strings.Contains(strings.ToLower(stock.Name), searchLower)
}

// containsString checks if a slice contains a string
func (s *MarketDataService) containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// matchesEquilibriumZone checks if stock is in any of the specified equilibrium zones
func (s *MarketDataService) matchesEquilibriumZone(stock models.StockData, zones []string) bool {
	// Determine which zone the stock is in
	var stockZone string
	if stock.PriceToEquilibrium < -5 {
		stockZone = "discount"
	} else if stock.PriceToEquilibrium > 5 {
		stockZone = "premium"
	} else {
		stockZone = "equilibrium"
	}

	// Check if stock's zone is in the requested zones
	return s.containsString(zones, stockZone)
}

// applySorting applies sorting to the stock list
func (s *MarketDataService) applySorting(stocks []models.StockData, sortField, sortOrder string) []models.StockData {
	sort.Slice(stocks, func(i, j int) bool {
		var aVal, bVal interface{}

		switch sortField {
		case "symbol":
			aVal, bVal = stocks[i].Symbol, stocks[j].Symbol
		case "name":
			aVal, bVal = stocks[i].Name, stocks[j].Name
		case "price":
			aVal, bVal = stocks[i].Price, stocks[j].Price
		case "changePercent":
			aVal, bVal = stocks[i].ChangePercent, stocks[j].ChangePercent
		case "rsi":
			aVal, bVal = stocks[i].RSI, stocks[j].RSI
		case "trend":
			aVal, bVal = stocks[i].Trend, stocks[j].Trend
		case "signal":
			aVal, bVal = stocks[i].Signal, stocks[j].Signal
		case "sector":
			aVal, bVal = stocks[i].Sector, stocks[j].Sector
		default:
			aVal, bVal = stocks[i].Symbol, stocks[j].Symbol
		}

		// Handle numeric comparison
		if aNum, aOk := aVal.(float64); aOk {
			if bNum, bOk := bVal.(float64); bOk {
				if sortOrder == "desc" {
					return aNum > bNum
				}
				return aNum < bNum
			}
		}

		// Handle string comparison
		aStr := fmt.Sprintf("%v", aVal)
		bStr := fmt.Sprintf("%v", bVal)
		if sortOrder == "desc" {
			return aStr > bStr
		}
		return aStr < bStr
	})

	return stocks
}

// generateCacheKey creates a cache key from the request
func (s *MarketDataService) generateCacheKey(req models.StockListRequest) string {
	// Create a hash of the request parameters for caching
	key := fmt.Sprintf("%s_%s_%d_%d_%s_%.1f_%.1f_%.1f_%.1f_%s_%s_%s_%s_%s",
		req.SortField,
		req.SortOrder,
		req.Page,
		req.PageSize,
		req.SearchTerm,
		req.RSIMin,
		req.RSIMax,
		req.PriceMin,
		req.PriceMax,
		strings.Join(req.Sectors, ","),
		strings.Join(req.Signals, ","),
		strings.Join(req.Trend, ","),
		strings.Join(req.VolumeProfile, ","),
		strings.Join(req.EquilibriumZone, ","),
	)
	return key
}

// getRealStockData fetches real stock data from Yahoo Finance
func (s *MarketDataService) getRealStockData() ([]models.StockData, error) {
	// List of popular symbols to scan
	symbols := []string{
		"AAPL", "MSFT", "GOOGL", "AMZN", "NVDA", "TSLA", "META", "BRK-B",
		"JNJ", "JPM", "V", "PG", "MA", "HD", "BAC", "XOM",
		"CVX", "ABBV", "KO", "PFE", "COST", "AVGO", "WMT", "DIS",
	}

	var stocks []models.StockData
	ctx := context.Background()

	for _, symbol := range symbols {
		// Fetch quote
		quote, err := s.yahooProvider.GetQuote(ctx, symbol)
		if err != nil {
			fmt.Printf("Failed to fetch %s: %v\n", symbol, err)
			continue
		}

		// Fetch historical data for indicators (90 days)
		historical, err := s.yahooProvider.GetHistoricalPrices(ctx, symbol, 90)
		if err != nil || len(historical) == 0 {
			fmt.Printf("Failed to fetch historical data for %s: %v\n", symbol, err)
			continue
		}

		// Calculate technical indicators
		indicators := s.yahooProvider.CalculateTechnicalIndicators(historical)

		// Calculate equilibrium
		equilibrium := s.indicators.CalculateEquilibrium(historical, quote.Price)

		// Determine trend based on moving averages
		trend := "neutral"
		if quote.Price > indicators.SMA50 && indicators.SMA50 > indicators.SMA200 {
			trend = "bullish"
		} else if quote.Price < indicators.SMA50 && indicators.SMA50 < indicators.SMA200 {
			trend = "bearish"
		}

		// Determine signal based on RSI and equilibrium
		signal := "hold"
		if indicators.RSI < 30 {
			signal = "buy"
		} else if indicators.RSI > 70 {
			signal = "sell"
		} else if equilibrium.Support > 0 {
			priceToEquilibrium := ((quote.Price - equilibrium.Support) / equilibrium.Support) * 100
			if priceToEquilibrium < -15 {
				signal = "buy"
			} else if priceToEquilibrium > 15 {
				signal = "sell"
			}
		}

		// Volume profile
		volumeProfile := "medium"
		if quote.Volume > 50000000 {
			volumeProfile = "high"
		} else if quote.Volume < 10000000 {
			volumeProfile = "low"
		}

		// Create stock data
		stock := models.StockData{
			Symbol:                 quote.Symbol,
			Name:                   quote.Name,
			Price:                  quote.Price,
			Change:                 quote.Change,
			ChangePercent:          quote.ChangePercent,
			Volume:                 quote.Volume,
			Sector:                 quote.Sector,
			Industry:               quote.Industry,
			MarketCap:              float64(quote.MarketCap),
			RSI:                    indicators.RSI,
			StochRSI:               indicators.StochRSI,
			HistoricRSIAvg:         indicators.HistoricRSIAvg,
			SMA50:                  indicators.SMA50,
			SMA200:                 indicators.SMA200,
			EMA20:                  indicators.EMA20,
			MACD:                   indicators.MACD,
			MACDSignal:             indicators.MACDSignal,
			MACDHistogram:          indicators.MACDHistogram,
			EquilibriumLevel:       (equilibrium.Support + equilibrium.Resistance) / 2,
			PriceToEquilibrium:     ((quote.Price - equilibrium.Support) / equilibrium.Support) * 100,
			Trend:                  trend,
			Signal:                 signal,
			VolumeProfile:          volumeProfile,
			DistanceFrom52WeekHigh: ((quote.Price - quote.Week52High) / quote.Week52High) * 100,
			DistanceFrom52WeekLow:  ((quote.Price - quote.Week52Low) / quote.Week52Low) * 100,
			LastUpdated:            time.Now(),
		}

		stocks = append(stocks, stock)
	}

	if len(stocks) == 0 {
		return nil, fmt.Errorf("no stock data retrieved")
	}

	return stocks, nil
}

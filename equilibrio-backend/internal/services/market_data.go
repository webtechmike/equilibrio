package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"equilibrio-backend/internal/config"
	"equilibrio-backend/internal/models"

	"github.com/redis/go-redis/v9"
)

type MarketDataService struct {
	config *config.Config
	cache  *redis.Client
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
		config: cfg,
		cache:  rdb,
	}
}

// GetStocks retrieves and filters stocks based on the request
func (s *MarketDataService) GetStocks(req models.StockListRequest) ([]models.StockData, int, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("stocks:%s", s.generateCacheKey(req))
	cached, err := s.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var cachedData struct {
			Stocks []models.StockData `json:"stocks"`
			Total  int                `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &cachedData) == nil {
			return cachedData.Stocks, cachedData.Total, nil
		}
	}

	// Generate mock data (replace with real API calls)
	stocks := s.generateMockStockData()

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

	// Cache the result
	cacheData := struct {
		Stocks []models.StockData `json:"stocks"`
		Total  int                `json:"total"`
	}{
		Stocks: paginatedStocks,
		Total:  total,
	}

	if data, err := json.Marshal(cacheData); err == nil {
		s.cache.Set(context.Background(), cacheKey, data, 30*time.Second)
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

// GetSectors returns all available sectors
func (s *MarketDataService) GetSectors() ([]string, error) {
	sectors := []string{
		"Technology", "Healthcare", "Financial", "Consumer Cyclical",
		"Energy", "Industrials", "Consumer Defensive", "Real Estate",
		"Communication Services", "Utilities", "Basic Materials",
	}
	return sectors, nil
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
	for _, ticker := range tickers {
		basePrice := rand.Float64()*500 + 50
		changePercent := (rand.Float64() - 0.5) * 10
		rsi := rand.Float64() * 100
		sma50 := basePrice * (0.9 + rand.Float64()*0.2)
		sma200 := basePrice * (0.85 + rand.Float64()*0.3)
		high52Week := basePrice * (1 + rand.Float64()*0.3)
		low52Week := basePrice * (0.7 + rand.Float64()*0.2)

		// Calculate equilibrium (50% retracement from low to high)
		equilibriumLevel := (high52Week + low52Week) / 2
		priceToEquilibrium := ((basePrice - equilibriumLevel) / equilibriumLevel) * 100

		macd := (rand.Float64() - 0.5) * 5
		macdSignal := macd + (rand.Float64()-0.5)*2

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
			randSignal := rand.Float64()
			if randSignal < 0.2 {
				signal = "buy"
			} else if randSignal < 0.4 {
				signal = "sell"
			} else {
				signal = "hold"
			}
		}

		// Volume profile based on volume
		avgVolume := rand.Float64() * 100000000
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
			Industry:               s.getIndustryForSector(ticker.sector),
			MarketCap:              basePrice * (rand.Float64()*1000000000 + 100000000),
			RSI:                    rsi,
			StochRSI:               rand.Float64() * 100,
			HistoricRSIAvg:         50 + (rand.Float64()-0.5)*20,
			SMA50:                  sma50,
			SMA200:                 sma200,
			EMA20:                  basePrice * (0.95 + rand.Float64()*0.1),
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
		return sectorIndustries[rand.Intn(len(sectorIndustries))]
	}
	return "General"
}

// applyFilters applies the filter criteria to the stock list
func (s *MarketDataService) applyFilters(stocks []models.StockData, filter models.StockFilter) []models.StockData {
	var filtered []models.StockData

	for _, stock := range stocks {
		// Search term filter
		if filter.SearchTerm != "" {
			searchLower := strings.ToLower(filter.SearchTerm)
			if !strings.Contains(strings.ToLower(stock.Symbol), searchLower) &&
				!strings.Contains(strings.ToLower(stock.Name), searchLower) {
				continue
			}
		}

		// Sector filter
		if len(filter.Sectors) > 0 {
			found := false
			for _, sector := range filter.Sectors {
				if stock.Sector == sector {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// RSI filter
		if stock.RSI < filter.RSIMin || stock.RSI > filter.RSIMax {
			continue
		}

		// Price filter
		if stock.Price < filter.PriceMin || stock.Price > filter.PriceMax {
			continue
		}

		// Volume profile filter
		if len(filter.VolumeProfile) > 0 {
			found := false
			for _, profile := range filter.VolumeProfile {
				if stock.VolumeProfile == profile {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Signal filter
		if len(filter.Signals) > 0 {
			found := false
			for _, signal := range filter.Signals {
				if stock.Signal == signal {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Trend filter
		if len(filter.Trend) > 0 {
			found := false
			for _, trend := range filter.Trend {
				if stock.Trend == trend {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Equilibrium zone filter
		if len(filter.EquilibriumZone) > 0 {
			inDiscount := stock.PriceToEquilibrium < -5
			inEquilibrium := stock.PriceToEquilibrium >= -5 && stock.PriceToEquilibrium <= 5
			inPremium := stock.PriceToEquilibrium > 5

			found := false
			for _, zone := range filter.EquilibriumZone {
				if (zone == "discount" && inDiscount) ||
					(zone == "equilibrium" && inEquilibrium) ||
					(zone == "premium" && inPremium) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, stock)
	}

	return filtered
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

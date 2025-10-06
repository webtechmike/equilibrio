package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"equilibrio-backend/internal/models"

	"github.com/redis/go-redis/v9"
)

// MarketCacheStrategy handles intelligent caching based on market hours
type MarketCacheStrategy struct {
	cache *redis.Client
}

// NewMarketCacheStrategy creates a new market-aware cache strategy
func NewMarketCacheStrategy(cache *redis.Client) *MarketCacheStrategy {
	return &MarketCacheStrategy{
		cache: cache,
	}
}

// MarketHours represents US stock market hours
type MarketHours struct {
	OpenHour   int // 9 AM ET (market open)
	CloseHour  int // 4 PM ET (market close)
	TimeZone   string
}

var USMarketHours = MarketHours{
	OpenHour:  9,
	CloseHour: 16,
	TimeZone:  "America/New_York",
}

// IsMarketOpen checks if the US stock market is currently open
func (m *MarketCacheStrategy) IsMarketOpen() bool {
	// Load Eastern Time
	loc, err := time.LoadLocation(USMarketHours.TimeZone)
	if err != nil {
		// Fallback to UTC if timezone load fails
		loc = time.UTC
	}

	now := time.Now().In(loc)
	
	// Check if it's a weekday (Monday = 1, Sunday = 0)
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}

	// Check if within market hours (9:30 AM - 4:00 PM ET)
	hour := now.Hour()
	minute := now.Minute()
	
	// Market opens at 9:30 AM
	if hour < 9 || (hour == 9 && minute < 30) {
		return false
	}
	
	// Market closes at 4:00 PM
	if hour >= 16 {
		return false
	}

	return true
}

// GetCacheTTL returns the appropriate cache TTL based on market status
func (m *MarketCacheStrategy) GetCacheTTL() time.Duration {
	if m.IsMarketOpen() {
		// During market hours: refresh every 5 minutes
		return 5 * time.Minute
	}

	// After market close: cache until next market open
	return m.TimeUntilNextMarketOpen()
}

// TimeUntilNextMarketOpen calculates time until next market open
func (m *MarketCacheStrategy) TimeUntilNextMarketOpen() time.Duration {
	loc, err := time.LoadLocation(USMarketHours.TimeZone)
	if err != nil {
		// Fallback: cache for 12 hours
		return 12 * time.Hour
	}

	now := time.Now().In(loc)
	
	// Calculate next market open (9:30 AM ET)
	var nextOpen time.Time
	
	switch now.Weekday() {
	case time.Saturday:
		// Next open is Monday
		nextOpen = time.Date(now.Year(), now.Month(), now.Day()+2, 9, 30, 0, 0, loc)
	case time.Sunday:
		// Next open is Monday
		nextOpen = time.Date(now.Year(), now.Month(), now.Day()+1, 9, 30, 0, 0, loc)
	case time.Friday:
		if now.Hour() >= 16 {
			// After Friday close, next open is Monday
			nextOpen = time.Date(now.Year(), now.Month(), now.Day()+3, 9, 30, 0, 0, loc)
		} else {
			// Before Friday close, next open is today
			nextOpen = time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, loc)
		}
	default:
		// Weekday
		if now.Hour() >= 16 {
			// After market close, next open is tomorrow
			nextOpen = time.Date(now.Year(), now.Month(), now.Day()+1, 9, 30, 0, 0, loc)
		} else if now.Hour() < 9 || (now.Hour() == 9 && now.Minute() < 30) {
			// Before market open, next open is today
			nextOpen = time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, loc)
		} else {
			// During market hours (shouldn't reach here)
			return 5 * time.Minute
		}
	}

	duration := nextOpen.Sub(now)
	
	// Ensure minimum TTL of 1 hour to handle timezone edge cases
	if duration < time.Hour {
		return time.Hour
	}

	return duration
}

// CacheStockData caches stock data with market-aware TTL
func (m *MarketCacheStrategy) CacheStockData(ctx context.Context, key string, data interface{}) error {
	ttl := m.GetCacheTTL()
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return m.cache.Set(ctx, key, jsonData, ttl).Err()
}

// CacheDailySnapshot stores a snapshot of all stocks for the trading day
func (m *MarketCacheStrategy) CacheDailySnapshot(ctx context.Context, stocks []models.StockData) error {
	loc, _ := time.LoadLocation(USMarketHours.TimeZone)
	now := time.Now().In(loc)
	
	// Create daily snapshot key (e.g., "daily_snapshot:2025-10-06")
	snapshotKey := fmt.Sprintf("daily_snapshot:%s", now.Format("2006-01-02"))
	
	// Cache until next market open
	ttl := m.TimeUntilNextMarketOpen()
	
	jsonData, err := json.Marshal(stocks)
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	return m.cache.Set(ctx, snapshotKey, jsonData, ttl).Err()
}

// GetDailySnapshot retrieves the daily snapshot if available
func (m *MarketCacheStrategy) GetDailySnapshot(ctx context.Context) ([]models.StockData, error) {
	loc, _ := time.LoadLocation(USMarketHours.TimeZone)
	now := time.Now().In(loc)
	
	snapshotKey := fmt.Sprintf("daily_snapshot:%s", now.Format("2006-01-02"))
	
	data, err := m.cache.Get(ctx, snapshotKey).Result()
	if err != nil {
		return nil, err
	}

	var stocks []models.StockData
	if err := json.Unmarshal([]byte(data), &stocks); err != nil {
		return nil, err
	}

	return stocks, nil
}

// WarmCache fetches and caches all stock data for the day
func (m *MarketCacheStrategy) WarmCache(ctx context.Context, fetchFunc func() ([]models.StockData, error)) error {
	// Check if we already have today's snapshot
	if _, err := m.GetDailySnapshot(ctx); err == nil {
		// Snapshot exists, no need to warm
		return nil
	}

	// Fetch fresh data
	stocks, err := fetchFunc()
	if err != nil {
		return fmt.Errorf("failed to fetch data for cache warming: %w", err)
	}

	// Store daily snapshot
	return m.CacheDailySnapshot(ctx, stocks)
}

// GetMarketStatus returns a human-readable market status
func (m *MarketCacheStrategy) GetMarketStatus() string {
	if m.IsMarketOpen() {
		return "Market Open"
	}

	loc, _ := time.LoadLocation(USMarketHours.TimeZone)
	now := time.Now().In(loc)
	
	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return "Market Closed (Weekend)"
	}

	if now.Hour() < 9 || (now.Hour() == 9 && now.Minute() < 30) {
		return "Market Closed (Pre-Market)"
	}

	if now.Hour() >= 16 {
		return "Market Closed (After Hours)"
	}

	return "Market Closed"
}


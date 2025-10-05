package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"equilibrio-backend/internal/config"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
	config *config.Config
}

func NewCacheService(cfg *config.Config) *CacheService {
	// Initialize Redis client
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		// Fallback to default Redis connection
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	rdb := redis.NewClient(opt)

	return &CacheService{
		client: rdb,
		config: cfg,
	}
}

// Set stores a value in the cache with expiration
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value from the cache
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete removes a key from the cache
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in the cache
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	return result > 0, err
}

// FlushDB clears all keys in the current database
func (c *CacheService) FlushDB(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// SetStockData caches stock data with a standard expiration
func (c *CacheService) SetStockData(ctx context.Context, symbol string, data interface{}) error {
	key := fmt.Sprintf("stock:%s", symbol)
	return c.Set(ctx, key, data, 5*time.Minute)
}

// GetStockData retrieves cached stock data
func (c *CacheService) GetStockData(ctx context.Context, symbol string, dest interface{}) error {
	key := fmt.Sprintf("stock:%s", symbol)
	return c.Get(ctx, key, dest)
}

// SetStocksList caches a list of stocks with filters
func (c *CacheService) SetStocksList(ctx context.Context, filterKey string, data interface{}) error {
	key := fmt.Sprintf("stocks:%s", filterKey)
	return c.Set(ctx, key, data, 30*time.Second)
}

// GetStocksList retrieves cached stocks list
func (c *CacheService) GetStocksList(ctx context.Context, filterKey string, dest interface{}) error {
	key := fmt.Sprintf("stocks:%s", filterKey)
	return c.Get(ctx, key, dest)
}

// SetIndicators caches technical indicators
func (c *CacheService) SetIndicators(ctx context.Context, symbol string, period int, data interface{}) error {
	key := fmt.Sprintf("indicators:%s:%d", symbol, period)
	return c.Set(ctx, key, data, 10*time.Minute)
}

// GetIndicators retrieves cached technical indicators
func (c *CacheService) GetIndicators(ctx context.Context, symbol string, period int, dest interface{}) error {
	key := fmt.Sprintf("indicators:%s:%d", symbol, period)
	return c.Get(ctx, key, dest)
}

// Close closes the Redis connection
func (c *CacheService) Close() error {
	return c.client.Close()
}

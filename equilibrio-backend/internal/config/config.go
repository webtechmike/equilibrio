package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port               string
	RedisURL           string
	AlphaVantageKey    string
	IEXCloudKey        string
	FinnhubAPIKey      string
	MarketDataProvider string
	Environment        string
	UseMockData        bool
	CacheTTL           int // Cache TTL in seconds
	MaxRequestsPerMin  int // Rate limiting
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		AlphaVantageKey:    getEnv("ALPHA_VANTAGE_API_KEY", ""),
		IEXCloudKey:        getEnv("IEX_CLOUD_API_KEY", ""),
		FinnhubAPIKey:      getEnv("FINNHUB_API_KEY", ""),
		MarketDataProvider: getEnv("MARKET_DATA_PROVIDER", "yahoo"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		UseMockData:        getEnvAsBool("USE_MOCK_DATA", false),
		CacheTTL:           getEnvAsInt("CACHE_TTL_SECONDS", 300), // 5 minutes default
		MaxRequestsPerMin:  getEnvAsInt("MAX_REQUESTS_PER_MINUTE", 60),
	}
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "true" || value == "1" || value == "yes" {
			return true
		}
		if value == "false" || value == "0" || value == "no" {
			return false
		}
	}
	return defaultValue
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

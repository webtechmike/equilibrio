package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            string
	RedisURL        string
	AlphaVantageKey string
	IEXCloudKey     string
	Environment     string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		AlphaVantageKey: getEnv("ALPHA_VANTAGE_API_KEY", ""),
		IEXCloudKey:     getEnv("IEX_CLOUD_API_KEY", ""),
		Environment:     getEnv("ENVIRONMENT", "development"),
	}
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

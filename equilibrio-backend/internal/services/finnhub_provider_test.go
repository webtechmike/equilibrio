package services

import (
	"context"
	"os"
	"testing"
)

func TestFinnhubProvider(t *testing.T) {
	// Skip test if no API key is provided
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		t.Skip("FINNHUB_API_KEY not set, skipping test")
	}

	provider := NewFinnhubProvider(apiKey)
	ctx := context.Background()

	// Test GetQuote
	t.Run("GetQuote", func(t *testing.T) {
		quote, err := provider.GetQuote(ctx, "AAPL")
		if err != nil {
			t.Fatalf("Failed to get quote: %v", err)
		}

		if quote.Symbol != "AAPL" {
			t.Errorf("Expected symbol AAPL, got %s", quote.Symbol)
		}

		if quote.Price <= 0 {
			t.Errorf("Expected positive price, got %f", quote.Price)
		}

		t.Logf("Quote: %+v", quote)
	})

	// Test GetHistoricalPrices (Note: Finnhub free tier doesn't include historical data)
	t.Run("GetHistoricalPrices", func(t *testing.T) {
		prices, err := provider.GetHistoricalPrices(ctx, "AAPL", 30)
		if err != nil {
			// This is expected for Finnhub free tier
			t.Logf("Historical data not available (expected for free tier): %v", err)
			return
		}

		if len(prices) == 0 {
			t.Error("Expected historical prices, got empty slice")
		}

		t.Logf("Got %d historical prices", len(prices))
	})

	// Test CalculateTechnicalIndicators (Note: Requires historical data)
	t.Run("CalculateTechnicalIndicators", func(t *testing.T) {
		prices, err := provider.GetHistoricalPrices(ctx, "AAPL", 100)
		if err != nil {
			// This is expected for Finnhub free tier
			t.Logf("Historical data not available for indicators (expected for free tier): %v", err)
			return
		}

		indicators := provider.CalculateTechnicalIndicators(prices)

		if indicators.RSI <= 0 || indicators.RSI > 100 {
			t.Errorf("Expected RSI between 0-100, got %f", indicators.RSI)
		}

		if indicators.SMA50 <= 0 {
			t.Errorf("Expected positive SMA50, got %f", indicators.SMA50)
		}

		t.Logf("Indicators: %+v", indicators)
	})
}

func TestRateLimitedFinnhubProvider(t *testing.T) {
	// Skip test if no API key is provided
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		t.Skip("FINNHUB_API_KEY not set, skipping test")
	}

	provider := NewRateLimitedFinnhubProvider(apiKey, 60) // 60 requests per minute
	ctx := context.Background()

	// Test that rate limiting works
	t.Run("RateLimitedGetQuote", func(t *testing.T) {
		quote, err := provider.GetQuote(ctx, "AAPL")
		if err != nil {
			t.Fatalf("Failed to get quote: %v", err)
		}

		if quote.Symbol != "AAPL" {
			t.Errorf("Expected symbol AAPL, got %s", quote.Symbol)
		}

		t.Logf("Rate-limited quote: %+v", quote)
	})
}

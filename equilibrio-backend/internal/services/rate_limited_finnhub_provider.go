package services

import (
	"context"

	"equilibrio-backend/internal/models"
)

// RateLimitedFinnhubProvider wraps FinnhubProvider with rate limiting
type RateLimitedFinnhubProvider struct {
	provider    *FinnhubProvider
	rateLimiter *RateLimiter
}

// NewRateLimitedFinnhubProvider creates a new rate-limited Finnhub provider
func NewRateLimitedFinnhubProvider(apiKey string, requestsPerMinute int) *RateLimitedFinnhubProvider {
	return &RateLimitedFinnhubProvider{
		provider:    NewFinnhubProvider(apiKey),
		rateLimiter: NewRateLimiter(requestsPerMinute),
	}
}

// GetQuote fetches real-time quote with rate limiting
func (p *RateLimitedFinnhubProvider) GetQuote(ctx context.Context, symbol string) (*models.Quote, error) {
	if err := p.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	return p.provider.GetQuote(ctx, symbol)
}

// GetHistoricalPrices fetches historical data with rate limiting
func (p *RateLimitedFinnhubProvider) GetHistoricalPrices(ctx context.Context, symbol string, days int) ([]models.CandlestickData, error) {
	if err := p.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	return p.provider.GetHistoricalPrices(ctx, symbol, days)
}

// SearchSymbols searches for symbols with rate limiting
func (p *RateLimitedFinnhubProvider) SearchSymbols(ctx context.Context, query string) ([]string, error) {
	if err := p.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	return p.provider.SearchSymbols(ctx, query)
}

// GetMarketSnapshot gets market snapshot with rate limiting
func (p *RateLimitedFinnhubProvider) GetMarketSnapshot(ctx context.Context, symbols []string) (map[string]*models.Quote, error) {
	snapshot := make(map[string]*models.Quote)

	for _, symbol := range symbols {
		// Rate limit each individual request
		if err := p.rateLimiter.Wait(ctx); err != nil {
			continue // Skip this symbol if rate limited
		}

		quote, err := p.provider.GetQuote(ctx, symbol)
		if err != nil {
			continue // Skip this symbol if error
		}
		snapshot[symbol] = quote
	}

	return snapshot, nil
}

// CalculateTechnicalIndicators calculates technical indicators (no rate limiting needed)
func (p *RateLimitedFinnhubProvider) CalculateTechnicalIndicators(candles []models.CandlestickData) models.TechnicalIndicators {
	return p.provider.CalculateTechnicalIndicators(candles)
}

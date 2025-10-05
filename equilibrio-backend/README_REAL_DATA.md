# Real Market Data Integration Guide

## Current State
The application currently uses **mock data** for demonstration purposes. The mock data:
- ✅ Generates consistent data for the same symbols
- ✅ Price changes are deterministic based on seed
- ⚠️ **Does NOT match real market data from TradingView or other sources**
- ⚠️ Data is randomly generated, not real-time

## Why Mock Data?
Mock data allows us to:
1. Develop and test the UI/UX without API dependencies
2. Avoid rate limits during development
3. Work offline
4. Test edge cases with controlled data

## Integrating Real Market Data

### Option 1: Yahoo Finance (Recommended for Development)
**Pros:**
- Free, no API key required
- Good data quality
- Go library available: `github.com/piquette/finance-go`

**Cons:**
- Unofficial API (could break)
- No guaranteed SLA

**Implementation:**
```go
import (
    "github.com/piquette/finance-go/quote"
)

func (s *MarketDataService) GetRealQuote(symbol string) (*models.Quote, error) {
    q, err := quote.Get(symbol)
    if err != nil {
        return nil, err
    }
    
    return &models.Quote{
        Symbol:        q.Symbol,
        Price:         q.RegularMarketPrice,
        Change:        q.RegularMarketChange,
        ChangePercent: q.RegularMarketChangePercent,
        Volume:        int64(q.RegularMarketVolume),
        // ... map other fields
    }, nil
}
```

### Option 2: Alpha Vantage
**Pros:**
- Official API with free tier
- 25 requests/day on free tier
- Reliable and well-documented

**Cons:**
- Requires API key
- Low rate limit on free tier

**Setup:**
1. Get API key: https://www.alphavantage.co/support/#api-key
2. Add to `.env`: `ALPHA_VANTAGE_API_KEY=your_key`
3. Use library: `github.com/alpacahq/alpaca-trade-api-go`

### Option 3: Finnhub
**Pros:**
- Real-time data
- 60 calls/minute on free tier
- WebSocket support for live updates

**Cons:**
- Requires API key
- Free tier has limitations

**Setup:**
1. Get API key: https://finnhub.io/
2. Add to `.env`: `FINNHUB_API_KEY=your_key`
3. Use REST API or WebSocket

### Option 4: Polygon.io
**Pros:**
- Very comprehensive data
- Good for production

**Cons:**
- Free tier: only 5 calls/minute
- Paid tier required for real-time data

## Implementation Steps

### 1. Create Real Data Provider
File: `internal/services/real_market_data.go`

```go
package services

import (
    "context"
    "equilibrio-backend/internal/models"
)

type RealMarketDataProvider struct {
    apiKey string
    // Add HTTP client, rate limiter, etc.
}

func NewRealMarketDataProvider(apiKey string) *RealMarketDataProvider {
    return &RealMarketDataProvider{
        apiKey: apiKey,
    }
}

func (p *RealMarketDataProvider) GetQuote(ctx context.Context, symbol string) (*models.Quote, error) {
    // Implement real API call
    // TODO: Replace with actual implementation
    return nil, nil
}

func (p *RealMarketDataProvider) GetHistoricalPrices(ctx context.Context, symbol string, days int) ([]models.CandlestickData, error) {
    // Implement historical data fetch
    // TODO: Replace with actual implementation
    return nil, nil
}
```

### 2. Update MarketDataService
Replace mock data generation with real API calls:

```go
// In internal/services/market_data.go
func (s *MarketDataService) GetStock(symbol string) (*models.StockData, error) {
    // Check cache first
    cached, err := s.getFromCache(symbol)
    if err == nil && cached != nil {
        return cached, nil
    }

    // Fetch from real API
    quote, err := s.marketProvider.GetQuote(context.Background(), symbol)
    if err != nil {
        return nil, err
    }

    // Calculate technical indicators
    indicators := s.indicatorsService.Calculate(quote, historicalPrices)

    // Convert to StockData
    stock := s.convertQuoteToStockData(quote, indicators)

    // Cache it
    s.saveToCache(symbol, stock)

    return stock, nil
}
```

### 3. Add Environment Configuration
File: `.env`

```bash
# Market Data Provider
MARKET_DATA_PROVIDER=yahoo  # or "alphavantage", "finnhub", "polygon"
ALPHA_VANTAGE_API_KEY=your_key_here
FINNHUB_API_KEY=your_key_here
POLYGON_API_KEY=your_key_here

# Rate Limiting
MAX_REQUESTS_PER_MINUTE=60
CACHE_TTL_SECONDS=300  # 5 minutes
```

### 4. Update Config
File: `internal/config/config.go`

```go
type Config struct {
    // ... existing fields
    MarketDataProvider string
    AlphaVantageAPIKey string
    FinnhubAPIKey      string
    PolygonAPIKey      string
    MaxRequestsPerMin  int
    CacheTTL           int
}
```

## Testing Strategy

### 1. Keep Mock Data for Tests
Move mock generation to test files only:
- `internal/services/market_data_test.go` ✅ (already done)
- Use mock data for unit tests
- Use real API (with caching) for integration tests

### 2. Feature Flag
Add environment variable to switch between mock and real:

```go
if os.Getenv("USE_MOCK_DATA") == "true" {
    return s.generateMockStockData()
}
return s.getRealStockData()
```

## Rate Limiting & Caching

### Rate Limiting
Implement rate limiting to avoid hitting API limits:

```go
import "golang.org/x/time/rate"

type RateLimitedProvider struct {
    limiter *rate.Limiter
    provider MarketDataProvider
}

func (p *RateLimitedProvider) GetQuote(ctx context.Context, symbol string) (*models.Quote, error) {
    if err := p.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    return p.provider.GetQuote(ctx, symbol)
}
```

### Caching Strategy
- Cache quotes for 5 minutes (configurable)
- Cache historical data for 1 hour
- Use Redis for distributed caching
- Implement cache warming for popular symbols

## Current Status

- ✅ Mock data generation working
- ✅ Filter logic refactored and simplified
- ✅ Chart price consistency fixed
- ✅ MarketDataProvider interface created
- ⏳ Real API integration (pending API key selection)
- ⏳ Rate limiting (pending)
- ⏳ Error handling for API failures (pending)

## Next Steps

1. Choose a market data provider (recommend Yahoo Finance for dev)
2. Install required Go library
3. Implement RealMarketDataProvider
4. Add feature flag to switch between mock/real
5. Add rate limiting
6. Update error handling
7. Test with real data
8. Deploy to production

## Notes

- **Important**: Real market data will make your filters work against actual market conditions
- The mock data is useful for development but should not be used to make trading decisions
- Consider costs and rate limits when choosing a provider
- Implement proper error handling and fallbacks


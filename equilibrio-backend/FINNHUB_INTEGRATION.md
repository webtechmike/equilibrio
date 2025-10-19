# Finnhub Integration Guide

## Overview

The Equilibrio backend now supports Finnhub as a market data provider alongside Yahoo Finance. Finnhub provides real-time and historical market data with a generous free tier.

## Setup

### 1. Get Finnhub API Key

1. Visit [Finnhub.io](https://finnhub.io/)
2. Sign up for a free account
3. Get your API key from the dashboard

### 2. Configure Environment

Add your Finnhub API key to your `.env` file:

```bash
# Market Data Provider Configuration
MARKET_DATA_PROVIDER=finnhub  # Options: yahoo, finnhub
FINNHUB_API_KEY=your_finnhub_api_key_here

# Rate Limiting (Finnhub free tier: 60 calls/minute)
MAX_REQUESTS_PER_MINUTE=60
CACHE_TTL_SECONDS=300  # 5 minutes
```

### 3. Start the Server

The server will automatically use Finnhub when `MARKET_DATA_PROVIDER=finnhub` is set.

```bash
go run cmd/server/main.go
```

## Features

### Real-time Quotes
- Current price, change, and change percentage
- Volume, market cap, and company information
- 52-week high/low (approximated)
- Sector and industry classification

### Historical Data
- Daily candlestick data (OHLC) - **Uses Yahoo Finance (free)**
- Configurable time periods
- Used for technical indicator calculations

**Note**: Finnhub's free tier doesn't include historical data. The system automatically uses Yahoo Finance for historical data when using Finnhub for real-time quotes.

### Technical Indicators
- RSI (Relative Strength Index)
- SMA (Simple Moving Average) - 50 and 200 day
- EMA (Exponential Moving Average) - 20 day
- MACD (Moving Average Convergence Divergence)

### Rate Limiting
- Built-in rate limiting to respect Finnhub's API limits
- Free tier: 60 calls per minute
- Automatic request queuing and throttling

## API Endpoints

All existing endpoints work with Finnhub:

- `GET /api/stocks` - Get filtered stock list
- `GET /api/stocks/{symbol}` - Get individual stock data
- `GET /api/stocks/{symbol}/chart` - Get historical chart data

## Example Usage

### Using Finnhub Provider Directly

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "equilibrio-backend/internal/services"
)

func main() {
    // Initialize Finnhub provider
    provider := services.NewRateLimitedFinnhubProvider("your_api_key", 60)
    ctx := context.Background()
    
    // Get real-time quote
    quote, err := provider.GetQuote(ctx, "AAPL")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("AAPL: $%.2f (%.2f%%)\n", quote.Price, quote.ChangePercent)
    
    // Get historical data
    prices, err := provider.GetHistoricalPrices(ctx, "AAPL", 30)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Got %d days of historical data\n", len(prices))
    
    // Calculate technical indicators
    indicators := provider.CalculateTechnicalIndicators(prices)
    fmt.Printf("RSI: %.2f, SMA50: %.2f\n", indicators.RSI, indicators.SMA50)
}
```

### Switching Between Providers

The system automatically chooses the provider based on configuration:

```go
// In your config
config := &config.Config{
    MarketDataProvider: "finnhub", // or "yahoo"
    FinnhubAPIKey:     "your_key",
    // ... other config
}

// The MarketDataService will automatically use the correct provider
service := services.NewMarketDataService(config)
```

## Finnhub API Limits

### Free Tier
- **60 calls per minute**
- Real-time quotes
- Historical data
- Company profiles
- Basic technical indicators

### Paid Tiers
- Higher rate limits
- WebSocket real-time data
- News and sentiment data
- Advanced technical indicators

## Error Handling

The system includes robust error handling:

1. **API Failures**: Falls back to mock data if Finnhub is unavailable
2. **Rate Limiting**: Automatically throttles requests
3. **Invalid Symbols**: Returns appropriate error messages
4. **Network Issues**: Retries with exponential backoff

## Caching Strategy

- **Real-time quotes**: Cached for 5 minutes (configurable)
- **Historical data**: Cached for 1 hour
- **Company profiles**: Cached for 24 hours
- **Market snapshots**: Cached based on market hours

## Testing

Run the Finnhub integration tests:

```bash
# Set your API key first
export FINNHUB_API_KEY=your_api_key_here

# Run tests
go test ./internal/services -v -run TestFinnhubProvider
```

## Troubleshooting

### Common Issues

1. **"API key not set"**: Make sure `FINNHUB_API_KEY` is in your `.env` file
2. **Rate limit exceeded**: Reduce `MAX_REQUESTS_PER_MINUTE` or upgrade your Finnhub plan
3. **No data returned**: Check if the symbol exists and is valid
4. **Network timeouts**: Increase timeout values in the HTTP client

### Debug Mode

Enable debug logging to see provider selection:

```bash
ENVIRONMENT=development go run cmd/server/main.go
```

This will show which provider is being used for each request.

## Migration from Yahoo Finance

To switch from Yahoo Finance to Finnhub:

1. Get your Finnhub API key
2. Update `.env`:
   ```bash
   MARKET_DATA_PROVIDER=finnhub
   FINNHUB_API_KEY=your_key
   ```
3. Restart the server

The system will automatically use Finnhub for all market data requests.

## Performance Considerations

- **Rate Limiting**: Finnhub free tier has 60 calls/minute limit
- **Caching**: Aggressive caching reduces API calls
- **Batch Requests**: Multiple symbols are processed efficiently
- **Fallback**: Mock data ensures service availability

## Security

- API keys are loaded from environment variables
- No API keys are logged or exposed in responses
- Rate limiting prevents abuse
- Input validation on all symbol parameters

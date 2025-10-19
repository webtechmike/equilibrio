# Equilibrio Backend

A Go-based backend service for the Equilibrio stock scanner application, providing real-time stock data, technical indicators, and filtering capabilities.

## Features

- **Stock Data API**: Retrieve and filter stock information
- **Multiple Data Providers**: Support for Finnhub and Yahoo Finance
- **Technical Indicators**: Calculate RSI, MACD, moving averages, and equilibrium levels
- **Real-time Updates**: WebSocket support for live data
- **Rate Limiting**: Built-in rate limiting for API providers
- **Caching**: Redis-based caching for improved performance
- **Export**: CSV export functionality
- **RESTful API**: Clean REST API design

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Redis server
- (Optional) Market data API keys (Finnhub, Alpha Vantage, IEX Cloud)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd equilibrio-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp env.example .env
# Edit .env with your configuration
```

4. Start Redis (if not already running):
```bash
redis-server
```

5. Run the application:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Stocks
- `GET /api/stocks` - Get filtered list of stocks
- `GET /api/stocks/:symbol` - Get specific stock data
- `GET /api/sectors` - Get available sectors
- `GET /api/export` - Export stocks to CSV

### Data Management
- `POST /api/refresh` - Refresh all stock data

### Technical Indicators
- `POST /api/indicators` - Calculate technical indicators

## Query Parameters

### GET /api/stocks

- `searchTerm` - Search by symbol or name
- `sectors` - Filter by sectors (comma-separated)
- `rsiMin`, `rsiMax` - RSI range filter
- `priceMin`, `priceMax` - Price range filter
- `volumeProfile` - Volume profile filter (high, medium, low)
- `signals` - Signal filter (buy, sell, hold)
- `trend` - Trend filter (bullish, bearish, neutral)
- `equilibriumZone` - Equilibrium zone filter (discount, equilibrium, premium)
- `sortField` - Sort field (symbol, price, changePercent, rsi, etc.)
- `sortOrder` - Sort order (asc, desc)
- `page` - Page number (default: 1)
- `pageSize` - Items per page (default: 50)

## Market Data Providers

The backend supports multiple market data providers:

### Finnhub (Recommended)
- **Free tier**: 60 calls/minute
- **Real-time quotes**: Current price, change, volume
- **Historical data**: Daily OHLC data
- **Company profiles**: Sector, industry, market cap
- **Setup**: Get API key from [Finnhub.io](https://finnhub.io/)

### Yahoo Finance
- **Free**: No API key required
- **Real-time quotes**: Basic stock information
- **Historical data**: Daily candlestick data
- **Limitations**: Unofficial API, rate limits may apply

### Configuration
Set your preferred provider in `.env`:
```bash
MARKET_DATA_PROVIDER=finnhub  # or "yahoo"
FINNHUB_API_KEY=your_finnhub_api_key
MAX_REQUESTS_PER_MINUTE=60
```

For detailed setup instructions, see [FINNHUB_INTEGRATION.md](FINNHUB_INTEGRATION.md).

## Environment Variables

- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment (development, production)
- `REDIS_URL` - Redis connection URL
- `MARKET_DATA_PROVIDER` - Provider to use (finnhub, yahoo)
- `FINNHUB_API_KEY` - Finnhub API key
- `ALPHA_VANTAGE_API_KEY` - Alpha Vantage API key
- `IEX_CLOUD_API_KEY` - IEX Cloud API key
- `MAX_REQUESTS_PER_MINUTE` - Rate limiting (default: 60)
- `CACHE_TTL_SECONDS` - Cache TTL in seconds (default: 300)
- `CORS_ORIGIN` - CORS origin for frontend

## Docker

Build and run with Docker:

```bash
# Build the image
docker build -t equilibrio-backend .

# Run the container
docker run -p 8080:8080 --env-file .env equilibrio-backend
```

## Development

### Project Structure

```
equilibrio-backend/
├── cmd/server/           # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and routes
│   ├── config/          # Configuration management
│   ├── models/          # Data models
│   └── services/        # Business logic services
├── pkg/                 # Reusable packages
├── go.mod
├── go.sum
└── Dockerfile
```

### Adding New Features

1. Define models in `internal/models/`
2. Implement business logic in `internal/services/`
3. Create API handlers in `internal/api/`
4. Add routes in `internal/api/routes.go`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License

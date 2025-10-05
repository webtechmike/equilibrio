# Equilibrio Stock Scanner

A comprehensive stock analysis application built with React frontend and Go backend, featuring equilibrium-based swing trading analysis.

## Architecture

This project consists of two main components:

- **Frontend**: React SPA with TypeScript, Vite, and Tailwind CSS
- **Backend**: Go API with Gin framework, Redis caching, and technical indicators

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd equilibrio
```

2. Start all services:
```bash
docker-compose up -d
```

3. Access the application:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Redis: localhost:6379

### Local Development

#### Backend Setup

1. Navigate to backend directory:
```bash
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

4. Start Redis:
```bash
docker run -d -p 6379:6379 redis:7-alpine
```

5. Run the backend:
```bash
go run cmd/server/main.go
```

#### Frontend Setup

1. Navigate to frontend directory:
```bash
cd equilibrio-frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

## Features

### Stock Analysis
- Real-time stock data with technical indicators
- RSI, MACD, Moving Averages analysis
- Equilibrium level calculations (50% retracement)
- Trend analysis and trading signals

### Filtering & Search
- Advanced filtering by sector, RSI, price range
- Signal-based filtering (buy/sell/hold)
- Equilibrium zone filtering (discount/equilibrium/premium)
- Real-time search by symbol or company name

### Data Export
- CSV export functionality
- Filtered data export
- Real-time data refresh

### Technical Features
- Responsive design for all devices
- Real-time data updates
- Caching for improved performance
- RESTful API design
- Type-safe frontend with TypeScript

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Stocks
- `GET /api/stocks` - Get filtered stock list
- `GET /api/stocks/:symbol` - Get specific stock data
- `GET /api/sectors` - Get available sectors
- `GET /api/export` - Export stocks to CSV

### Data Management
- `POST /api/refresh` - Refresh all stock data

### Technical Indicators
- `POST /api/indicators` - Calculate technical indicators

## Project Structure

```
equilibrio/
├── equilibrio-backend/          # Go backend service
│   ├── cmd/server/             # Application entry point
│   ├── internal/               # Private application code
│   │   ├── api/               # HTTP handlers and routes
│   │   ├── config/            # Configuration management
│   │   ├── models/            # Data models
│   │   └── services/          # Business logic services
│   ├── pkg/                   # Reusable packages
│   ├── go.mod
│   ├── Dockerfile
│   └── README.md
├── equilibrio-frontend/         # React frontend
│   ├── src/
│   │   ├── components/        # React components
│   │   ├── hooks/            # Custom React hooks
│   │   ├── services/         # API services
│   │   ├── types/            # TypeScript types
│   │   └── App.tsx           # Main app component
│   ├── package.json
│   ├── Dockerfile
│   └── README.md
├── docker-compose.yml          # Docker orchestration
└── README.md                   # This file
```

## Configuration

### Environment Variables

#### Backend (.env)
```bash
PORT=8080
ENVIRONMENT=development
REDIS_URL=redis://localhost:6379
ALPHA_VANTAGE_API_KEY=your_key_here
IEX_CLOUD_API_KEY=your_key_here
CORS_ORIGIN=http://localhost:3000
```

#### Frontend
```bash
VITE_API_URL=/api
VITE_APP_TITLE=Equilibrio Stock Scanner
```

## Deployment

### Production Deployment

1. **Using Docker Compose**:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

2. **Manual Deployment**:
   - Build and deploy backend to your server
   - Build and deploy frontend to a CDN or static hosting
   - Configure reverse proxy (nginx) for API routing

### Cloud Deployment

- **Frontend**: Deploy to Vercel, Netlify, or AWS S3 + CloudFront
- **Backend**: Deploy to AWS ECS, Google Cloud Run, or Railway
- **Database**: Use managed Redis (AWS ElastiCache, Redis Cloud)

## Development

### Adding New Features

1. **Backend**: Add new endpoints in `internal/api/`
2. **Frontend**: Create components in `src/components/`
3. **Types**: Update TypeScript interfaces in `src/types/`
4. **API**: Add new service methods in `src/services/`

### Code Style

- **Go**: Follow standard Go conventions
- **React**: Use functional components with hooks
- **TypeScript**: Strict mode enabled
- **CSS**: Tailwind CSS utility classes

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License

## Support

For questions or issues, please open an issue on GitHub or contact the development team.

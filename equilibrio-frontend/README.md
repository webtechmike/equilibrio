# Equilibrio Frontend

A modern React-based frontend for the Equilibrio stock scanner application, built with TypeScript, Vite, and Tailwind CSS.

## Features

- **Modern React**: Built with React 18 and TypeScript
- **Real-time Data**: Live stock data updates with React Query
- **Advanced Filtering**: Comprehensive filtering and sorting capabilities
- **Responsive Design**: Mobile-first design with Tailwind CSS
- **Component Architecture**: Modular, reusable components
- **Type Safety**: Full TypeScript support
- **Performance**: Optimized with Vite and code splitting

## Quick Start

### Prerequisites

- Node.js 18 or higher
- npm or yarn
- Equilibrio backend running on port 8080

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
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

The application will be available at `http://localhost:3000`

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint
- `npm run lint:fix` - Fix ESLint errors

## Project Structure

```
equilibrio-frontend/
├── src/
│   ├── components/          # React components
│   │   ├── StockHeader.tsx  # Header with search and actions
│   │   ├── StockFilters.tsx # Filter controls
│   │   ├── StockTable.tsx   # Main data table
│   │   ├── TradingInsights.tsx # Trading insights
│   │   └── EquilibriumInfo.tsx # Info component
│   ├── hooks/              # Custom React hooks
│   │   └── useStocks.ts    # Stock data hooks
│   ├── services/           # API services
│   │   └── api.ts          # API client
│   ├── types/              # TypeScript types
│   │   └── index.ts        # Type definitions
│   ├── App.tsx             # Main app component
│   ├── main.tsx            # App entry point
│   └── index.css           # Global styles
├── public/                 # Static assets
├── package.json
├── vite.config.ts          # Vite configuration
├── tailwind.config.js      # Tailwind configuration
└── tsconfig.json           # TypeScript configuration
```

## Components

### StockHeader
- Search functionality
- Refresh and export buttons
- Application title and description

### StockFilters
- RSI range filtering
- Price range filtering
- Sector selection
- Signal and trend filtering
- Equilibrium zone filtering
- Volume profile filtering

### StockTable
- Sortable columns
- Expandable rows with detailed information
- Technical indicators display
- Trading insights
- Responsive design

### TradingInsights
- AI-powered trading recommendations
- Technical analysis insights
- Risk assessment

## API Integration

The frontend communicates with the Go backend through a RESTful API:

- **GET /api/stocks** - Retrieve filtered stock data
- **GET /api/stocks/:symbol** - Get specific stock information
- **GET /api/sectors** - Get available sectors
- **POST /api/refresh** - Refresh all data
- **GET /api/export** - Export data to CSV
- **POST /api/indicators** - Calculate technical indicators

## State Management

- **React Query**: Server state management and caching
- **React Hooks**: Local state management
- **Custom Hooks**: Reusable state logic

## Styling

- **Tailwind CSS**: Utility-first CSS framework
- **Responsive Design**: Mobile-first approach
- **Custom Components**: Reusable styled components
- **Dark Mode Ready**: Prepared for dark mode implementation

## Docker Deployment

Build and run with Docker:

```bash
# Build the image
docker build -t equilibrio-frontend .

# Run the container
docker run -p 80:80 equilibrio-frontend
```

## Environment Variables

- `VITE_API_URL` - Backend API URL (default: /api)
- `VITE_APP_TITLE` - Application title

## Development

### Adding New Features

1. Create components in `src/components/`
2. Add types in `src/types/`
3. Implement API calls in `src/services/`
4. Create custom hooks in `src/hooks/`
5. Update the main App component

### Code Style

- Use TypeScript for all new code
- Follow React best practices
- Use functional components with hooks
- Implement proper error handling
- Add loading states for async operations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License

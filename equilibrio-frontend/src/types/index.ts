// Stock data types matching the backend models
export interface StockData {
  symbol: string;
  name: string;
  price: number;
  change: number;
  changePercent: number;
  volume: number;
  sector: string;
  industry: string;
  marketCap: number;
  rsi: number;
  stochRsi: number;
  historicRsiAvg: number;
  sma50: number;
  sma200: number;
  ema20: number;
  macd: number;
  macdSignal: number;
  macdHistogram: number;
  equilibriumLevel: number;
  priceToEquilibrium: number;
  trend: 'bullish' | 'bearish' | 'neutral';
  signal: 'buy' | 'sell' | 'hold';
  volumeProfile: 'high' | 'medium' | 'low';
  distanceFrom52WeekHigh: number;
  distanceFrom52WeekLow: number;
  lastUpdated: string;
}

export interface StockFilter {
  searchTerm: string;
  sectors: string[];
  rsiMin: number;
  rsiMax: number;
  priceMin: number;
  priceMax: number;
  volumeProfile: string[];
  signals: string[];
  trend: string[];
  equilibriumZone: string[];
}

export interface StockListRequest {
  // Filter fields (flattened to match backend)
  searchTerm: string;
  sectors: string[];
  rsiMin: number;
  rsiMax: number;
  priceMin: number;
  priceMax: number;
  volumeProfile: string[];
  signals: string[];
  trend: string[];
  equilibriumZone: string[];
  
  // Pagination and sorting
  sortField: string;
  sortOrder: 'asc' | 'desc';
  page: number;
  pageSize: number;
}

export interface StockListResponse {
  stocks: StockData[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface TechnicalIndicators {
  rsi: number;
  stochRsi: number;
  historicRsiAvg: number;
  sma50: number;
  sma200: number;
  ema20: number;
  macd: number;
  macdSignal: number;
  macdHistogram: number;
}

// API response types
export interface ApiResponse<T> {
  data: T;
  error?: string;
}

// Filter options
export interface FilterOptions {
  sectors: string[];
  volumeProfiles: string[];
  signals: string[];
  trends: string[];
  equilibriumZones: string[];
}

// Sort options
export interface SortOption {
  field: string;
  label: string;
  type: 'string' | 'number';
}

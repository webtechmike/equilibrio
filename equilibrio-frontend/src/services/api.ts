import axios, { AxiosResponse } from 'axios';
import { StockData, StockListRequest, StockListResponse, TechnicalIndicators, CandlestickData } from '../types';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for logging
api.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

// API service class
export class ApiService {
  // Get stocks with filtering and pagination
  static async getStocks(request: StockListRequest): Promise<StockListResponse> {
    const params = new URLSearchParams();
    
    // Add filter parameters (flattened structure)
    if (request.searchTerm) {
      params.append('searchTerm', request.searchTerm);
    }
    if (request.sectors.length > 0) {
      params.append('sectors', request.sectors.join(','));
    }
    if (request.rsiMin > 0) {
      params.append('rsiMin', request.rsiMin.toString());
    }
    if (request.rsiMax < 100) {
      params.append('rsiMax', request.rsiMax.toString());
    }
    if (request.priceMin > 0) {
      params.append('priceMin', request.priceMin.toString());
    }
    if (request.priceMax < 10000) {
      params.append('priceMax', request.priceMax.toString());
    }
    if (request.volumeProfile.length > 0) {
      params.append('volumeProfile', request.volumeProfile.join(','));
    }
    if (request.signals.length > 0) {
      params.append('signals', request.signals.join(','));
    }
    if (request.trend.length > 0) {
      params.append('trend', request.trend.join(','));
    }
    if (request.equilibriumZone.length > 0) {
      params.append('equilibriumZone', request.equilibriumZone.join(','));
    }
    
    // Add sorting and pagination parameters
    params.append('sortField', request.sortField);
    params.append('sortOrder', request.sortOrder);
    params.append('page', request.page.toString());
    params.append('pageSize', request.pageSize.toString());

    const response: AxiosResponse<StockListResponse> = await api.get(`/stocks?${params.toString()}`);
    return response.data;
  }

  // Get single stock by symbol
  static async getStock(symbol: string): Promise<StockData> {
    const response: AxiosResponse<StockData> = await api.get(`/stocks/${symbol}`);
    return response.data;
  }

  // Get stock chart data
  static async getStockChart(symbol: string, days: number = 90): Promise<CandlestickData[]> {
    const response: AxiosResponse<{ symbol: string; data: CandlestickData[] }> = await api.get(`/stocks/${symbol}/chart?days=${days}`);
    return response.data.data;
  }

  // Get available sectors
  static async getSectors(): Promise<string[]> {
    const response: AxiosResponse<{ sectors: string[] }> = await api.get('/sectors');
    return response.data.sectors;
  }

  // Calculate technical indicators
  static async calculateIndicators(symbol: string, period: number = 200): Promise<TechnicalIndicators> {
    const response: AxiosResponse<TechnicalIndicators> = await api.post('/indicators', {
      symbol,
      period,
    });
    return response.data;
  }

  // Refresh all data
  static async refreshData(): Promise<void> {
    await api.post('/refresh');
  }

  // Export stocks to CSV
  static async exportStocks(request: StockListRequest): Promise<Blob> {
    const params = new URLSearchParams();
    
    // Add filter parameters (flattened structure)
    if (request.searchTerm) {
      params.append('searchTerm', request.searchTerm);
    }
    if (request.sectors.length > 0) {
      params.append('sectors', request.sectors.join(','));
    }
    if (request.rsiMin > 0) {
      params.append('rsiMin', request.rsiMin.toString());
    }
    if (request.rsiMax < 100) {
      params.append('rsiMax', request.rsiMax.toString());
    }
    if (request.priceMin > 0) {
      params.append('priceMin', request.priceMin.toString());
    }
    if (request.priceMax < 10000) {
      params.append('priceMax', request.priceMax.toString());
    }
    if (request.volumeProfile.length > 0) {
      params.append('volumeProfile', request.volumeProfile.join(','));
    }
    if (request.signals.length > 0) {
      params.append('signals', request.signals.join(','));
    }
    if (request.trend.length > 0) {
      params.append('trend', request.trend.join(','));
    }
    if (request.equilibriumZone.length > 0) {
      params.append('equilibriumZone', request.equilibriumZone.join(','));
    }

    const response = await api.get(`/export?${params.toString()}`, {
      responseType: 'blob',
    });
    
    return response.data;
  }

  // Health check
  static async healthCheck(): Promise<{ status: string; service: string }> {
    const response: AxiosResponse<{ status: string; service: string }> = await api.get('/health');
    return response.data;
  }
}

export default api;

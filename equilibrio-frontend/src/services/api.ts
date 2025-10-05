import axios, { AxiosResponse } from 'axios';
import { StockData, StockListRequest, StockListResponse, TechnicalIndicators } from '../types';

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
    
    // Add filter parameters
    if (request.filter.searchTerm) {
      params.append('searchTerm', request.filter.searchTerm);
    }
    if (request.filter.sectors.length > 0) {
      params.append('sectors', request.filter.sectors.join(','));
    }
    if (request.filter.rsiMin > 0) {
      params.append('rsiMin', request.filter.rsiMin.toString());
    }
    if (request.filter.rsiMax < 100) {
      params.append('rsiMax', request.filter.rsiMax.toString());
    }
    if (request.filter.priceMin > 0) {
      params.append('priceMin', request.filter.priceMin.toString());
    }
    if (request.filter.priceMax < 10000) {
      params.append('priceMax', request.filter.priceMax.toString());
    }
    if (request.filter.volumeProfile.length > 0) {
      params.append('volumeProfile', request.filter.volumeProfile.join(','));
    }
    if (request.filter.signals.length > 0) {
      params.append('signals', request.filter.signals.join(','));
    }
    if (request.filter.trend.length > 0) {
      params.append('trend', request.filter.trend.join(','));
    }
    if (request.filter.equilibriumZone.length > 0) {
      params.append('equilibriumZone', request.filter.equilibriumZone.join(','));
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
    
    // Add filter parameters (same as getStocks)
    if (request.filter.searchTerm) {
      params.append('searchTerm', request.filter.searchTerm);
    }
    if (request.filter.sectors.length > 0) {
      params.append('sectors', request.filter.sectors.join(','));
    }
    if (request.filter.rsiMin > 0) {
      params.append('rsiMin', request.filter.rsiMin.toString());
    }
    if (request.filter.rsiMax < 100) {
      params.append('rsiMax', request.filter.rsiMax.toString());
    }
    if (request.filter.priceMin > 0) {
      params.append('priceMin', request.filter.priceMin.toString());
    }
    if (request.filter.priceMax < 10000) {
      params.append('priceMax', request.filter.priceMax.toString());
    }
    if (request.filter.volumeProfile.length > 0) {
      params.append('volumeProfile', request.filter.volumeProfile.join(','));
    }
    if (request.filter.signals.length > 0) {
      params.append('signals', request.filter.signals.join(','));
    }
    if (request.filter.trend.length > 0) {
      params.append('trend', request.filter.trend.join(','));
    }
    if (request.filter.equilibriumZone.length > 0) {
      params.append('equilibriumZone', request.filter.equilibriumZone.join(','));
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

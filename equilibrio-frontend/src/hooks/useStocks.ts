import { useState, useCallback, useEffect } from 'react';
import { useQuery, useQueryClient } from 'react-query';
import { StockData, StockListRequest, StockListResponse, StockFilter } from '../types';
import { ApiService } from '../services/api';
import { saveFilterConfig, loadFilterConfig, clearFilterConfig } from '../utils/filterStorage';
import { useDebounce } from './useDebounce';

export const useStocks = (request: StockListRequest) => {
  const queryClient = useQueryClient();

  // Debounce the search term with 500ms delay
  const debouncedSearchTerm = useDebounce(request.searchTerm, 500);

  // Create a modified request with debounced search term
  const debouncedRequest = {
    ...request,
    searchTerm: debouncedSearchTerm,
  };

  // Only enable the query if search term is empty or has at least 3 characters
  const shouldFetch = !debouncedSearchTerm || debouncedSearchTerm.length >= 3;

  const query = useQuery<StockListResponse, Error>(
    ['stocks', debouncedRequest],
    () => ApiService.getStocks(debouncedRequest),
    {
      enabled: shouldFetch,
      staleTime: 30000, // 30 seconds
      cacheTime: 300000, // 5 minutes
      refetchOnWindowFocus: false,
    }
  );

  const refreshData = useCallback(async () => {
    await ApiService.refreshData();
    queryClient.invalidateQueries(['stocks']);
  }, [queryClient]);

  return {
    ...query,
    refreshData,
  };
};

export const useStock = (symbol: string) => {
  return useQuery<StockData, Error>(
    ['stock', symbol],
    () => ApiService.getStock(symbol),
    {
      enabled: !!symbol,
      staleTime: 30000,
      cacheTime: 300000,
    }
  );
};

export const useSectors = () => {
  return useQuery<string[], Error>(
    ['sectors'],
    () => ApiService.getSectors(),
    {
      staleTime: 300000, // 5 minutes
      cacheTime: 600000, // 10 minutes
    }
  );
};

const DEFAULT_FILTERS: StockFilter = {
  searchTerm: '',
  sectors: [],
  rsiMin: 0,
  rsiMax: 100,
  priceMin: 0,
  priceMax: 10000,
  volumeProfile: [],
  signals: [],
  trend: [],
  equilibriumZone: [],
};

export const useStockFilters = () => {
  // Initialize from localStorage or use defaults
  const [filters, setFilters] = useState<StockFilter>(() => {
    const saved = loadFilterConfig();
    return saved || DEFAULT_FILTERS;
  });

  // Auto-save to localStorage whenever filters change
  useEffect(() => {
    saveFilterConfig(filters);
  }, [filters]);

  const updateFilter = useCallback((key: keyof StockFilter, value: any) => {
    setFilters(prev => ({
      ...prev,
      [key]: value,
    }));
  }, []);

  const resetFilters = useCallback(() => {
    setFilters(DEFAULT_FILTERS);
    clearFilterConfig();
  }, []);

  const loadFilters = useCallback((newFilters: StockFilter) => {
    setFilters(newFilters);
  }, []);

  return {
    filters,
    updateFilter,
    resetFilters,
    loadFilters,
  };
};

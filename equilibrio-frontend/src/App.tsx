import React, { useState, useCallback, useMemo, useRef, useEffect } from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { ThemeProvider } from './contexts/ThemeContext';
import StockHeader from './components/StockHeader';
import StockFilters from './components/StockFilters';
import StockTable from './components/StockTable';
import CandlestickChart from './components/CandlestickChart';
import StockPriceCard from './components/StockPriceCard';
import EquilibriumInfo from './components/EquilibriumInfo';
import { useStocks, useSectors, useStockFilters } from './hooks/useStocks';
import { useFavorites } from './hooks/useFavorites';
import { ApiService } from './services/api';
import { StockListRequest, CandlestickData, StockData } from './types';
import { mergeFavoritesWithStocks } from './utils/stockUtils';

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

const App: React.FC = () => {
  const [sortField, setSortField] = useState<string>('symbol');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');
  const [expandedRow, setExpandedRow] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(50);
  const [selectedStock, setSelectedStock] = useState<StockData | null>(null);
  const [chartData, setChartData] = useState<CandlestickData[]>([]);
  const [manuallyAddedStocks, setManuallyAddedStocks] = useState<StockData[]>([]);

  // Ref for scrolling to price card
  const priceCardRef = useRef<HTMLDivElement>(null);

  const { filters, updateFilter, resetFilters, loadFilters } = useStockFilters();
  const { data: sectors = [] } = useSectors();
  const { favorites, toggleFavorite, isFavorite } = useFavorites();

  const request: StockListRequest = useMemo(() => ({
    // Flattened filter fields
    searchTerm: filters.searchTerm,
    sectors: filters.sectors,
    rsiMin: filters.rsiMin,
    rsiMax: filters.rsiMax,
    priceMin: filters.priceMin,
    priceMax: filters.priceMax,
    volumeProfile: filters.volumeProfile,
    signals: filters.signals,
    trend: filters.trend,
    equilibriumZone: filters.equilibriumZone,
    
    // Pagination and sorting
    sortField,
    sortOrder: sortDirection,
    page,
    pageSize,
  }), [filters, sortField, sortDirection, page, pageSize]);

  const { data: stocksData, isLoading, error, refreshData } = useStocks(request);

  // Merge favorites with stocks data and manually added stocks
  const mergedStocks = useMemo(() => {
    const backendStocks = stocksData?.stocks || [];
    const allStocks = [...manuallyAddedStocks, ...backendStocks];
    
    // Remove duplicates based on symbol (manually added stocks take precedence)
    const uniqueStocks = allStocks.reduce((acc, stock) => {
      const existingIndex = acc.findIndex(s => s.symbol.toUpperCase() === stock.symbol.toUpperCase());
      if (existingIndex === -1) {
        acc.push(stock);
      } else {
        // Replace with manually added stock if it exists
        if (manuallyAddedStocks.some(ms => ms.symbol.toUpperCase() === stock.symbol.toUpperCase())) {
          acc[existingIndex] = stock;
        }
      }
      return acc;
    }, [] as StockData[]);
    
    return mergeFavoritesWithStocks(uniqueStocks, favorites);
  }, [stocksData?.stocks, favorites, manuallyAddedStocks]);

  const handleSort = useCallback((field: string) => {
    if (sortField === field) {
      setSortDirection(prev => prev === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortDirection('asc');
    }
  }, [sortField]);

  const handleSearchChange = useCallback((value: string) => {
    updateFilter('searchTerm', value);
    setPage(1); // Reset to first page when searching
  }, [updateFilter]);

  const handleRefresh = useCallback(async () => {
    await refreshData();
  }, [refreshData]);

  const handleExport = useCallback(async () => {
    try {
      const blob = await ApiService.exportStocks(request);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `stock_scan_${new Date().toISOString().split('T')[0]}.csv`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Export failed:', error);
      alert('Export failed. Please try again.');
    }
  }, [request]);

  const handleRowExpand = useCallback((symbol: string | null) => {
    setExpandedRow(symbol);
  }, []);

  const handleStockClick = useCallback(async (stock: StockData) => {
    setSelectedStock(stock);
    try {
      const data = await ApiService.getStockChart(stock.symbol, 90); // Default to 3 months
      setChartData(data);
    } catch (error) {
      console.error('Failed to fetch chart data:', error);
      setChartData([]);
    }
  }, []);

  // Scroll to price card when stock is selected
  useEffect(() => {
    if (selectedStock && priceCardRef.current) {
      // Small delay to ensure DOM is updated
      setTimeout(() => {
        priceCardRef.current?.scrollIntoView({ 
          behavior: 'smooth', 
          block: 'start',
          inline: 'nearest'
        });
      }, 100);
    }
  }, [selectedStock]);

  const handleTimeframeChange = useCallback(async (days: number) => {
    if (!selectedStock) return;
    try {
      const data = await ApiService.getStockChart(selectedStock.symbol, days);
      setChartData(data);
    } catch (error) {
      console.error('Failed to fetch chart data:', error);
      setChartData([]);
    }
  }, [selectedStock]);

  const handleCloseChart = useCallback(() => {
    setSelectedStock(null);
    setChartData([]);
  }, []);

  const handleClosePriceCard = useCallback(() => {
    setSelectedStock(null);
  }, []);

  const handleAddStock = useCallback(async (stock: StockData) => {
    // Add the stock to manually added stocks so it appears in the list
    setManuallyAddedStocks(prev => {
      const existingIndex = prev.findIndex(s => s.symbol.toUpperCase() === stock.symbol.toUpperCase());
      if (existingIndex === -1) {
        return [stock, ...prev]; // Add to the beginning
      }
      return prev; // Already exists
    });
    
    // When a new stock is added via search, display it immediately
    setSelectedStock(stock);
    try {
      const data = await ApiService.getStockChart(stock.symbol, 90);
      setChartData(data);
    } catch (error) {
      console.error('Failed to fetch chart data:', error);
      setChartData([]);
    }
  }, []);

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800 p-6 transition-colors">
        <div className="max-w-7xl mx-auto">
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6">
            <h2 className="text-lg font-semibold text-red-800 dark:text-red-200 mb-2">Error Loading Data</h2>
            <p className="text-red-700 dark:text-red-300">{error.message}</p>
            <button
              onClick={handleRefresh}
              className="mt-4 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition"
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800 p-6 transition-colors">
      <div className="max-w-7xl mx-auto">
        <StockHeader
          filters={filters}
          onSearchChange={handleSearchChange}
          onRefresh={handleRefresh}
          onExport={handleExport}
          onLoadPreset={loadFilters}
          onAddStock={handleAddStock}
          currentStocks={mergedStocks}
          loading={isLoading}
        />

        <StockFilters
          filters={filters}
          sectors={sectors}
          onFilterChange={updateFilter}
          onResetFilters={resetFilters}
        />

        {/* Price Summary Card - Prominent Display with scroll target */}
        <div ref={priceCardRef}>
          <StockPriceCard 
            stock={selectedStock} 
            onClose={handleClosePriceCard}
          />
        </div>

        {selectedStock && chartData.length > 0 && (
          <CandlestickChart
            symbol={selectedStock.symbol}
            companyName={selectedStock.name}
            data={chartData}
            onClose={handleCloseChart}
            onTimeframeChange={handleTimeframeChange}
          />
        )}

        <StockTable
          stocks={mergedStocks}
          loading={isLoading}
          sortField={sortField}
          sortDirection={sortDirection}
          onSort={handleSort}
          onRowExpand={handleRowExpand}
          expandedRow={expandedRow}
          onStockClick={handleStockClick}
          onToggleFavorite={toggleFavorite}
          isFavorite={isFavorite}
        />

        <EquilibriumInfo />
      </div>
    </div>
  );
};

const AppWithQueryClient: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <App />
      </ThemeProvider>
    </QueryClientProvider>
  );
};

export default AppWithQueryClient;

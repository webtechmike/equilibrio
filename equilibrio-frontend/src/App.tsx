import React, { useState, useCallback } from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import StockHeader from './components/StockHeader';
import StockFilters from './components/StockFilters';
import StockTable from './components/StockTable';
import EquilibriumInfo from './components/EquilibriumInfo';
import { useStocks, useSectors, useStockFilters } from './hooks/useStocks';
import { ApiService } from './services/api';
import { StockListRequest } from './types';

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

  const { filters, updateFilter, resetFilters } = useStockFilters();
  const { data: sectors = [] } = useSectors();

  const request: StockListRequest = {
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
  };

  const { data: stocksData, isLoading, error, refreshData } = useStocks(request);

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

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="bg-red-50 border border-red-200 rounded-lg p-6">
            <h2 className="text-lg font-semibold text-red-800 mb-2">Error Loading Data</h2>
            <p className="text-red-700">{error.message}</p>
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
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 p-6">
      <div className="max-w-7xl mx-auto">
        <StockHeader
          filters={filters}
          onSearchChange={handleSearchChange}
          onRefresh={handleRefresh}
          onExport={handleExport}
          loading={isLoading}
        />

        <StockFilters
          filters={filters}
          sectors={sectors}
          onFilterChange={updateFilter}
          onResetFilters={resetFilters}
        />

        <StockTable
          stocks={stocksData?.stocks || []}
          loading={isLoading}
          sortField={sortField}
          sortDirection={sortDirection}
          onSort={handleSort}
          onRowExpand={handleRowExpand}
          expandedRow={expandedRow}
        />

        <EquilibriumInfo />
      </div>
    </div>
  );
};

const AppWithQueryClient: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  );
};

export default AppWithQueryClient;

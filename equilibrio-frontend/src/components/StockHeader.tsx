import React, { useState, useEffect, useCallback } from 'react';
import { Search, RefreshCw, Download, TrendingUp, AlertCircle, Star, Trash2 } from 'lucide-react';
import { StockFilter, StockFilters, StockData } from '../types';
import ThemeToggle from './ui/ThemeToggle';
import FilterPresets from './FilterPresets';
import { ApiService } from '../services/api';
import { useFavorites } from '../hooks/useFavorites';

interface StockHeaderProps {
  filters: StockFilter;
  onSearchChange: (value: string) => void;
  onRefresh: () => void;
  onExport: () => void;
  onLoadPreset: (filters: StockFilters) => void;
  onAddStock?: (stock: StockData) => void;
  currentStocks: StockData[];
  loading: boolean;
}

const StockHeader: React.FC<StockHeaderProps> = ({
  filters,
  onSearchChange,
  onRefresh,
  onExport,
  onLoadPreset,
  onAddStock,
  currentStocks,
  loading,
}) => {
  const [showAddButton, setShowAddButton] = useState(false);
  const [isAdding, setIsAdding] = useState(false);
  const [addError, setAddError] = useState<string | null>(null);
  const [addSuccess, setAddSuccess] = useState<string | null>(null);
  const [showFavorites, setShowFavorites] = useState(false);
  
  const { favorites, clearFavorites, addFavorite } = useFavorites();

  // Close favorites dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Element;
      if (showFavorites && !target.closest('.favorites-dropdown')) {
        setShowFavorites(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [showFavorites]);

  // Check if search term is a potential ticker symbol (not in current list)
  useEffect(() => {
    const searchTerm = filters.searchTerm.trim().toUpperCase();
    
    if (searchTerm.length >= 1 && searchTerm.length <= 5) {
      // Check if it matches any existing stock in the current list
      const exists = currentStocks.some(
        stock => stock.symbol.toUpperCase() === searchTerm || 
                stock.name.toUpperCase().includes(searchTerm)
      );
      
      setShowAddButton(!exists && searchTerm.length > 0);
    } else {
      setShowAddButton(false);
    }
  }, [filters.searchTerm, currentStocks]);

  const handleAddStock = useCallback(async () => {
    const symbol = filters.searchTerm.trim().toUpperCase();
    
    setIsAdding(true);
    setAddError(null);
    setAddSuccess(null);

    try {
      const stock = await ApiService.searchStock(symbol);
      
      if (onAddStock) {
        onAddStock(stock);
        // Add to favorites automatically
        addFavorite(stock.symbol);
        setAddSuccess(`${stock.symbol} added to favorites!`);
        
        // Clear search after 2 seconds
        setTimeout(() => {
          onSearchChange('');
          setAddSuccess(null);
        }, 2000);
      }
    } catch (err: any) {
      setAddError(err.response?.data?.error || `Symbol "${symbol}" not found`);
      
      // Clear error after 3 seconds
      setTimeout(() => {
        setAddError(null);
      }, 3000);
    } finally {
      setIsAdding(false);
    }
  }, [filters.searchTerm, onAddStock, onSearchChange]);

  return (
    <div className="bg-white dark:bg-slate-800 rounded-lg shadow-lg p-6 mb-6 transition-colors">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-3xl font-bold text-slate-800 dark:text-slate-100">Equilibrio Scanner</h1>
            <p className="text-slate-600 dark:text-slate-400 mt-1 font-licorice">Equilibrium-based swing trading analysis</p>
        </div>
        <div className="flex gap-3">
          <ThemeToggle />
          <FilterPresets currentFilters={filters} onLoadPreset={onLoadPreset} />
          
          {/* Favorites Management */}
          <div className="relative favorites-dropdown">
            <button
              onClick={() => setShowFavorites(!showFavorites)}
              className="flex items-center gap-2 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition"
            >
              <Star className="w-4 h-4 fill-current" />
              Favorites ({favorites.length})
            </button>
            
            {showFavorites && (
              <div className="absolute right-0 top-full mt-2 w-80 bg-white dark:bg-slate-800 rounded-lg shadow-lg border border-slate-200 dark:border-slate-700 z-10">
                <div className="p-4">
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="font-semibold text-slate-800 dark:text-slate-200">Favorites</h3>
                    {favorites.length > 0 && (
                      <button
                        onClick={clearFavorites}
                        className="text-red-600 hover:text-red-700 text-sm flex items-center gap-1"
                      >
                        <Trash2 className="w-4 h-4" />
                        Clear All
                      </button>
                    )}
                  </div>
                  
                  {favorites.length === 0 ? (
                    <p className="text-slate-500 dark:text-slate-400 text-sm">No favorites yet. Click the ⭐ on any stock to add it.</p>
                  ) : (
                    <div className="space-y-2 max-h-48 overflow-y-auto">
                      {favorites.map((symbol, index) => (
                        <div key={symbol} className="flex items-center justify-between p-2 bg-slate-50 dark:bg-slate-700 rounded">
                          <span className="font-mono text-sm text-slate-800 dark:text-slate-200">{symbol}</span>
                          <span className="text-xs text-slate-500 dark:text-slate-400">#{index + 1}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
          
          <button
            onClick={onRefresh}
            disabled={loading}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition"
          >
            <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </button>
          <button
            onClick={onExport}
            className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition"
          >
            <Download className="w-4 h-4" />
            Export CSV
          </button>
        </div>
      </div>

      {/* Search with Smart Add */}
      <div className="space-y-2">
        <div className="relative flex gap-2">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-slate-400 dark:text-slate-500 w-5 h-5" />
            <input
              type="text"
              placeholder="Search by symbol or name (min 3 chars)... (e.g., AAPL, TSLA, BTC-USD)"
              value={filters.searchTerm}
              onChange={(e) => onSearchChange(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && showAddButton && !isAdding) {
                  e.preventDefault();
                  handleAddStock();
                }
              }}
              className="w-full pl-10 pr-4 py-3 border border-slate-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100 placeholder-slate-400 dark:placeholder-slate-500 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors"
            />
          </div>
          
          {showAddButton && (
            <button
              onClick={handleAddStock}
              disabled={isAdding}
              className="flex items-center gap-2 px-4 py-3 bg-emerald-600 hover:bg-emerald-700 disabled:bg-slate-400 text-white rounded-lg font-medium transition-colors disabled:cursor-not-allowed"
              title={`Add ${filters.searchTerm.toUpperCase()} to favorites`}
            >
              {isAdding ? (
                <div className="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent" />
              ) : (
                <Star className="w-5 h-5 fill-current" />
              )}
              <span className="hidden sm:inline">Add to Favorites</span>
            </button>
          )}
        </div>

        {/* Success Message */}
        {addSuccess && (
          <div className="flex items-center gap-2 p-3 bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 rounded-lg text-emerald-700 dark:text-emerald-300">
            <TrendingUp size={20} />
            <span>{addSuccess}</span>
          </div>
        )}

        {/* Error Message */}
        {addError && (
          <div className="flex items-center gap-2 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-700 dark:text-red-300">
            <AlertCircle size={20} />
            <span>{addError}</span>
          </div>
        )}

        {/* Helper Text */}
        {filters.searchTerm.length > 0 && filters.searchTerm.length < 3 && (
          <p className="text-sm text-amber-600 dark:text-amber-400 flex items-center gap-2">
            <AlertCircle size={16} />
            <span>Enter at least 3 characters to search</span>
          </p>
        )}
        
        {showAddButton && !addError && !addSuccess && (
          <p className="text-sm text-slate-600 dark:text-slate-400 flex items-center gap-2">
            <TrendingUp size={16} />
            <span>
              <strong>{filters.searchTerm.toUpperCase()}</strong> not in list. Click <Star className="inline w-4 h-4 fill-current text-yellow-500" /> or press <kbd className="px-1 py-0.5 text-xs bg-slate-200 dark:bg-slate-600 rounded">Enter</kbd> to add to favorites!
            </span>
          </p>
        )}
      </div>
    </div>
  );
};

export default StockHeader;

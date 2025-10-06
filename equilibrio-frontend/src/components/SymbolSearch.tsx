import React, { useState, useCallback } from 'react';
import { Search, TrendingUp, AlertCircle } from 'lucide-react';
import { ApiService } from '../services/api';
import { StockData } from '../types';

interface SymbolSearchProps {
  onStockFound: (stock: StockData) => void;
}

const SymbolSearch: React.FC<SymbolSearchProps> = ({ onStockFound }) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [isSearching, setIsSearching] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [recentSearches, setRecentSearches] = useState<string[]>([]);

  const handleSearch = useCallback(async (query: string) => {
    if (!query || query.length < 1) {
      setError('Please enter a ticker symbol');
      return;
    }

    const symbol = query.toUpperCase().trim();
    setIsSearching(true);
    setError(null);

    try {
      const stock = await ApiService.searchStock(symbol);
      
      // Add to recent searches (keep last 5)
      setRecentSearches(prev => {
        const updated = [symbol, ...prev.filter(s => s !== symbol)].slice(0, 5);
        localStorage.setItem('recentSearches', JSON.stringify(updated));
        return updated;
      });

      // Clear search and notify parent
      setSearchQuery('');
      onStockFound(stock);
    } catch (err: any) {
      setError(err.response?.data?.error || `Symbol "${symbol}" not found`);
    } finally {
      setIsSearching(false);
    }
  }, [onStockFound]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    handleSearch(searchQuery);
  };

  const handleRecentClick = (symbol: string) => {
    handleSearch(symbol);
  };

  const clearRecentSearches = () => {
    setRecentSearches([]);
    localStorage.removeItem('recentSearches');
  };

  // Load recent searches on mount
  React.useEffect(() => {
    const saved = localStorage.getItem('recentSearches');
    if (saved) {
      try {
        setRecentSearches(JSON.parse(saved));
      } catch (e) {
        // Ignore parse errors
      }
    }
  }, []);

  return (
    <div className="bg-white dark:bg-slate-800 rounded-lg shadow-md p-6 mb-6 transition-colors">
      <div className="flex items-center gap-3 mb-4">
        <TrendingUp className="text-blue-600 dark:text-blue-400" size={24} />
        <h2 className="text-xl font-bold text-slate-800 dark:text-slate-100">
          Search Any Ticker Symbol
        </h2>
      </div>

      <form onSubmit={handleSubmit} className="mb-4">
        <div className="flex gap-2">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-slate-400 dark:text-slate-500" size={20} />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => {
                setSearchQuery(e.target.value.toUpperCase());
                setError(null);
              }}
              placeholder="Enter ticker symbol (e.g., AAPL, TSLA, COIN)"
              className="w-full pl-10 pr-4 py-3 bg-slate-50 dark:bg-slate-700 border border-slate-200 dark:border-slate-600 rounded-lg text-slate-800 dark:text-slate-100 placeholder-slate-400 dark:placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-400 transition-colors"
              disabled={isSearching}
            />
          </div>
          <button
            type="submit"
            disabled={isSearching || !searchQuery}
            className="px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-300 dark:disabled:bg-slate-700 text-white rounded-lg font-medium transition-colors disabled:cursor-not-allowed flex items-center gap-2"
          >
            {isSearching ? (
              <>
                <div className="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent" />
                Searching...
              </>
            ) : (
              <>
                <Search size={20} />
                Search
              </>
            )}
          </button>
        </div>
      </form>

      {error && (
        <div className="flex items-center gap-2 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg text-red-700 dark:text-red-300 mb-4">
          <AlertCircle size={20} />
          <span>{error}</span>
        </div>
      )}

      {recentSearches.length > 0 && (
        <div className="mt-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-slate-600 dark:text-slate-400">
              Recent Searches
            </span>
            <button
              onClick={clearRecentSearches}
              className="text-xs text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-300 transition-colors"
            >
              Clear
            </button>
          </div>
          <div className="flex flex-wrap gap-2">
            {recentSearches.map((symbol) => (
              <button
                key={symbol}
                onClick={() => handleRecentClick(symbol)}
                className="px-3 py-1.5 bg-slate-100 dark:bg-slate-700 hover:bg-slate-200 dark:hover:bg-slate-600 text-slate-700 dark:text-slate-200 rounded-md text-sm font-medium transition-colors flex items-center gap-1"
              >
                {symbol}
                <TrendingUp size={14} />
              </button>
            ))}
          </div>
        </div>
      )}

      <div className="mt-4 pt-4 border-t border-slate-200 dark:border-slate-700">
        <p className="text-sm text-slate-600 dark:text-slate-400">
          💡 <span className="font-medium">Tip:</span> Search for any publicly traded stock symbol.
          The system will fetch live data and add it to your watchlist.
        </p>
      </div>
    </div>
  );
};

export default SymbolSearch;


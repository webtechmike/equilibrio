import React from 'react';
import { Search, RefreshCw, Download } from 'lucide-react';
import { StockFilter } from '../types';

interface StockHeaderProps {
  filters: StockFilter;
  onSearchChange: (value: string) => void;
  onRefresh: () => void;
  onExport: () => void;
  loading: boolean;
}

const StockHeader: React.FC<StockHeaderProps> = ({
  filters,
  onSearchChange,
  onRefresh,
  onExport,
  loading,
}) => {
  return (
    <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h1 className="text-3xl font-bold text-slate-800">Stock Scanner</h1>
          <p className="text-slate-600 mt-1">Equilibrium-based swing trading analysis</p>
        </div>
        <div className="flex gap-3">
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

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-slate-400 w-5 h-5" />
        <input
          type="text"
          placeholder="Search by symbol or name..."
          value={filters.searchTerm}
          onChange={(e) => onSearchChange(e.target.value)}
          className="w-full pl-10 pr-4 py-3 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>
    </div>
  );
};

export default StockHeader;

import React, { useCallback } from 'react';
import { Filter } from 'lucide-react';
import { StockFilter } from '../types';
import FilterButton from './ui/FilterButton';
import {
  EQUILIBRIUM_ZONES,
  SIGNALS,
  VOLUME_PROFILES,
  TRENDS,
} from '../constants/filterOptions';

interface StockFiltersProps {
  filters: StockFilter;
  sectors: string[];
  onFilterChange: (key: keyof StockFilter, value: any) => void;
  onResetFilters: () => void;
}

const StockFilters: React.FC<StockFiltersProps> = ({
  filters,
  sectors,
  onFilterChange,
  onResetFilters,
}) => {
  const handleArrayToggle = useCallback((key: keyof StockFilter, value: string) => {
    const currentArray = filters[key] as string[];
    const newArray = currentArray.includes(value)
      ? currentArray.filter(item => item !== value)
      : [...currentArray, value];
    onFilterChange(key, newArray);
  }, [filters, onFilterChange]);

  return (
    <div className="bg-white rounded-lg shadow-lg p-6 mb-6">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Filter className="w-5 h-5 text-slate-600" />
          <h2 className="text-lg font-semibold text-slate-800">Filters</h2>
        </div>
        <button
          onClick={onResetFilters}
          className="text-sm text-blue-600 hover:text-blue-800 font-medium"
        >
          Reset All
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* RSI Range */}
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            RSI Range: {filters.rsiMin} - {filters.rsiMax}
          </label>
          <div className="flex gap-2">
            <input
              type="number"
              value={filters.rsiMin}
              onChange={(e) => onFilterChange('rsiMin', Number(e.target.value))}
              className="w-20 px-2 py-1 border border-slate-300 rounded"
              min="0"
              max="100"
            />
            <input
              type="number"
              value={filters.rsiMax}
              onChange={(e) => onFilterChange('rsiMax', Number(e.target.value))}
              className="w-20 px-2 py-1 border border-slate-300 rounded"
              min="0"
              max="100"
            />
          </div>
        </div>

        {/* Price Range */}
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            Price Range
          </label>
          <div className="flex gap-2">
            <input
              type="number"
              value={filters.priceMin}
              onChange={(e) => onFilterChange('priceMin', Number(e.target.value))}
              className="w-24 px-2 py-1 border border-slate-300 rounded"
              placeholder="Min"
            />
            <input
              type="number"
              value={filters.priceMax}
              onChange={(e) => onFilterChange('priceMax', Number(e.target.value))}
              className="w-24 px-2 py-1 border border-slate-300 rounded"
              placeholder="Max"
            />
          </div>
        </div>

        {/* Equilibrium Zone */}
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            Equilibrium Zone
          </label>
          <div className="flex flex-col gap-1">
            {EQUILIBRIUM_ZONES.map(zone => (
              <label key={zone} className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={filters.equilibriumZone.includes(zone)}
                  onChange={() => handleArrayToggle('equilibriumZone', zone)}
                  className="rounded"
                />
                <span className="text-sm capitalize">{zone}</span>
              </label>
            ))}
          </div>
        </div>

        {/* Signal */}
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            Signal
          </label>
          <div className="flex flex-col gap-1">
            {SIGNALS.map(sig => (
              <label key={sig} className="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={filters.signals.includes(sig)}
                  onChange={() => handleArrayToggle('signals', sig)}
                  className="rounded"
                />
                <span className="text-sm capitalize">{sig}</span>
              </label>
            ))}
          </div>
        </div>
      </div>

      {/* Additional Filters Row */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
        {/* Volume Profile */}
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            Volume Profile
          </label>
          <div className="flex gap-2">
            {VOLUME_PROFILES.map(profile => (
              <FilterButton
                key={profile}
                label={profile}
                active={filters.volumeProfile.includes(profile)}
                onClick={() => handleArrayToggle('volumeProfile', profile)}
              />
            ))}
          </div>
        </div>

        {/* Trend */}
        <div>
          <label className="block text-sm font-medium text-slate-700 mb-2">
            Trend
          </label>
          <div className="flex gap-2">
            {TRENDS.map(trend => (
              <FilterButton
                key={trend}
                label={trend}
                active={filters.trend.includes(trend)}
                onClick={() => handleArrayToggle('trend', trend)}
              />
            ))}
          </div>
        </div>
      </div>

      {/* Sector Filter */}
      <div className="mt-4">
        <label className="block text-sm font-medium text-slate-700 mb-2">
          Sectors
        </label>
        <div className="flex flex-wrap gap-2">
          {sectors.map(sector => (
            <FilterButton
              key={sector}
              label={sector}
              active={filters.sectors.includes(sector)}
              onClick={() => handleArrayToggle('sectors', sector)}
            />
          ))}
        </div>
      </div>
    </div>
  );
};

export default React.memo(StockFilters);

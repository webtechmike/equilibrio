/**
 * Stock utility functions for calculations and classifications
 */

export const getEquilibriumColor = (value: number): string => {
  if (value < -5) return 'text-green-600 bg-green-50';
  if (value > 5) return 'text-red-600 bg-red-50';
  return 'text-yellow-600 bg-yellow-50';
};

export const getEquilibriumTextColor = (value: number): string => {
  if (value < -5) return 'text-green-600';
  if (value > 5) return 'text-red-600';
  return 'text-yellow-600';
};

export const getEquilibriumZone = (value: number): string => {
  if (value < -5) return 'Discount';
  if (value > 5) return 'Premium';
  return 'Equilibrium';
};

export const getRSIColor = (rsi: number): string => {
  if (rsi < 30) return 'text-green-600';
  if (rsi > 70) return 'text-red-600';
  return 'text-slate-700';
};

export const getTrendColor = (trend: string): string => {
  switch (trend) {
    case 'bullish':
      return 'bg-green-100 text-green-800';
    case 'bearish':
      return 'bg-red-100 text-red-800';
    default:
      return 'bg-slate-100 text-slate-800';
  }
};

export const getSignalColor = (signal: string): string => {
  switch (signal) {
    case 'buy':
      return 'bg-green-600 text-white';
    case 'sell':
      return 'bg-red-600 text-white';
    default:
      return 'bg-slate-300 text-slate-700';
  }
};

export const getVolumeProfileColor = (profile: string): string => {
  switch (profile) {
    case 'high':
      return 'text-green-600';
    case 'low':
      return 'text-red-600';
    default:
      return 'text-yellow-600';
  }
};

export const getChangeColor = (changePercent: number): string => {
  return changePercent >= 0 ? 'text-green-600' : 'text-red-600';
};

export const formatPrice = (price: number): string => {
  return `$${price.toFixed(2)}`;
};

export const formatPercent = (value: number): string => {
  return `${value.toFixed(2)}%`;
};

export const formatVolume = (volume: number): string => {
  return `${(volume / 1000000).toFixed(2)}M`;
};

export const formatMarketCap = (marketCap: number): string => {
  return `$${(marketCap / 1000000000).toFixed(2)}B`;
};


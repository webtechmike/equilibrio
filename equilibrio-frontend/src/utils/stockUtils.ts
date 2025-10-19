import { StockData } from '../types';

export const mergeFavoritesWithStocks = (
  stocks: StockData[],
  favorites: string[]
): StockData[] => {
  if (favorites.length === 0) {
    return stocks;
  }

  // Separate favorites and non-favorites
  const favoriteStocks: StockData[] = [];
  const nonFavoriteStocks: StockData[] = [];

  stocks.forEach(stock => {
    if (favorites.includes(stock.symbol.toUpperCase())) {
      favoriteStocks.push(stock);
    } else {
      nonFavoriteStocks.push(stock);
    }
  });

  // Sort favorites by their order in the favorites array
  favoriteStocks.sort((a, b) => {
    const aIndex = favorites.indexOf(a.symbol.toUpperCase());
    const bIndex = favorites.indexOf(b.symbol.toUpperCase());
    return aIndex - bIndex;
  });

  // Return favorites first, then non-favorites
  return [...favoriteStocks, ...nonFavoriteStocks];
};

export const getFavoriteStocks = (
  stocks: StockData[],
  favorites: string[]
): StockData[] => {
  return stocks.filter(stock => 
    favorites.includes(stock.symbol.toUpperCase())
  );
};

// Utility functions for stock display
export const getRSIColor = (rsi: number): string => {
  if (rsi < 30) return 'text-green-600 dark:text-green-400';
  if (rsi > 70) return 'text-red-600 dark:text-red-400';
  return 'text-slate-600 dark:text-slate-400';
};

export const getEquilibriumColor = (priceToEquilibrium: number): string => {
  if (priceToEquilibrium < -5) return 'text-green-600 dark:text-green-400';
  if (priceToEquilibrium > 5) return 'text-red-600 dark:text-red-400';
  return 'text-slate-600 dark:text-slate-400';
};

export const getEquilibriumZone = (priceToEquilibrium: number): string => {
  if (priceToEquilibrium < -5) return 'discount';
  if (priceToEquilibrium > 5) return 'premium';
  return 'equilibrium';
};

export const getTrendColor = (trend: string): string => {
  switch (trend) {
    case 'bullish': return 'text-green-600 dark:text-green-400';
    case 'bearish': return 'text-red-600 dark:text-red-400';
    default: return 'text-slate-600 dark:text-slate-400';
  }
};

export const getSignalColor = (signal: string): string => {
  switch (signal) {
    case 'buy': return 'text-green-600 dark:text-green-400';
    case 'sell': return 'text-red-600 dark:text-red-400';
    default: return 'text-slate-600 dark:text-slate-400';
  }
};

export const formatPrice = (price: number): string => {
  return `$${price.toFixed(2)}`;
};

export const formatPercent = (percent: number): string => {
  const sign = percent >= 0 ? '+' : '';
  return `${sign}${percent.toFixed(2)}%`;
};

export const getChangeColor = (change: number): string => {
  if (change > 0) return 'text-green-600 dark:text-green-400';
  if (change < 0) return 'text-red-600 dark:text-red-400';
  return 'text-slate-600 dark:text-slate-400';
};

export const getEquilibriumTextColor = (priceToEquilibrium: number): string => {
  return getEquilibriumColor(priceToEquilibrium);
};

export const getVolumeProfileColor = (volumeProfile: string): string => {
  switch (volumeProfile) {
    case 'high': return 'text-green-600 dark:text-green-400';
    case 'low': return 'text-red-600 dark:text-red-400';
    default: return 'text-slate-600 dark:text-slate-400';
  }
};

export const formatVolume = (volume: number): string => {
  if (volume >= 1e9) return `${(volume / 1e9).toFixed(1)}B`;
  if (volume >= 1e6) return `${(volume / 1e6).toFixed(1)}M`;
  if (volume >= 1e3) return `${(volume / 1e3).toFixed(1)}K`;
  return volume.toString();
};

export const formatMarketCap = (marketCap: number): string => {
  if (marketCap >= 1e12) return `$${(marketCap / 1e12).toFixed(1)}T`;
  if (marketCap >= 1e9) return `$${(marketCap / 1e9).toFixed(1)}B`;
  if (marketCap >= 1e6) return `$${(marketCap / 1e6).toFixed(1)}M`;
  return `$${marketCap.toFixed(0)}`;
};
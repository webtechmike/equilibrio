import { useState, useCallback, useEffect } from 'react';

const FAVORITES_KEY = 'equilibrio-favorites';

export const useFavorites = () => {
  const [favorites, setFavorites] = useState<string[]>(() => {
    try {
      const saved = localStorage.getItem(FAVORITES_KEY);
      return saved ? JSON.parse(saved) : [];
    } catch {
      return [];
    }
  });

  // Save to localStorage whenever favorites change
  useEffect(() => {
    try {
      localStorage.setItem(FAVORITES_KEY, JSON.stringify(favorites));
    } catch (error) {
      console.error('Failed to save favorites:', error);
    }
  }, [favorites]);

  const addFavorite = useCallback((symbol: string) => {
    setFavorites(prev => {
      const upperSymbol = symbol.toUpperCase();
      if (!prev.includes(upperSymbol)) {
        return [...prev, upperSymbol];
      }
      return prev;
    });
  }, []);

  const removeFavorite = useCallback((symbol: string) => {
    setFavorites(prev => prev.filter(s => s !== symbol.toUpperCase()));
  }, []);

  const toggleFavorite = useCallback((symbol: string) => {
    const upperSymbol = symbol.toUpperCase();
    setFavorites(prev => {
      if (prev.includes(upperSymbol)) {
        return prev.filter(s => s !== upperSymbol);
      } else {
        return [...prev, upperSymbol];
      }
    });
  }, []);

  const isFavorite = useCallback((symbol: string) => {
    return favorites.includes(symbol.toUpperCase());
  }, [favorites]);

  const clearFavorites = useCallback(() => {
    setFavorites([]);
  }, []);

  return {
    favorites,
    addFavorite,
    removeFavorite,
    toggleFavorite,
    isFavorite,
    clearFavorites,
  };
};

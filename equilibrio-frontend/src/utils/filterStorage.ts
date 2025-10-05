import { StockFilters } from '../types';

const FILTER_STORAGE_KEY = 'equilibrio_filter_config';
const FILTER_PRESETS_KEY = 'equilibrio_filter_presets';

export interface FilterPreset {
  id: string;
  name: string;
  description: string;
  filters: StockFilters;
  createdAt: string;
  updatedAt: string;
}

/**
 * Save current filter configuration to localStorage
 */
export const saveFilterConfig = (filters: StockFilters): void => {
  try {
    const serialized = JSON.stringify(filters);
    localStorage.setItem(FILTER_STORAGE_KEY, serialized);
  } catch (error) {
    console.error('Failed to save filter config to localStorage:', error);
  }
};

/**
 * Load filter configuration from localStorage
 */
export const loadFilterConfig = (): StockFilters | null => {
  try {
    const serialized = localStorage.getItem(FILTER_STORAGE_KEY);
    if (serialized === null) {
      return null;
    }
    return JSON.parse(serialized) as StockFilters;
  } catch (error) {
    console.error('Failed to load filter config from localStorage:', error);
    return null;
  }
};

/**
 * Clear filter configuration from localStorage
 */
export const clearFilterConfig = (): void => {
  try {
    localStorage.removeItem(FILTER_STORAGE_KEY);
  } catch (error) {
    console.error('Failed to clear filter config from localStorage:', error);
  }
};

/**
 * Save a filter preset
 */
export const saveFilterPreset = (name: string, description: string, filters: StockFilters): FilterPreset => {
  const presets = loadFilterPresets();
  
  const preset: FilterPreset = {
    id: generateId(),
    name,
    description,
    filters,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };
  
  presets.push(preset);
  
  try {
    localStorage.setItem(FILTER_PRESETS_KEY, JSON.stringify(presets));
  } catch (error) {
    console.error('Failed to save filter preset:', error);
  }
  
  return preset;
};

/**
 * Load all filter presets
 */
export const loadFilterPresets = (): FilterPreset[] => {
  try {
    const serialized = localStorage.getItem(FILTER_PRESETS_KEY);
    if (serialized === null) {
      return [];
    }
    return JSON.parse(serialized) as FilterPreset[];
  } catch (error) {
    console.error('Failed to load filter presets:', error);
    return [];
  }
};

/**
 * Delete a filter preset by ID
 */
export const deleteFilterPreset = (id: string): void => {
  const presets = loadFilterPresets();
  const filtered = presets.filter(p => p.id !== id);
  
  try {
    localStorage.setItem(FILTER_PRESETS_KEY, JSON.stringify(filtered));
  } catch (error) {
    console.error('Failed to delete filter preset:', error);
  }
};

/**
 * Update an existing filter preset
 */
export const updateFilterPreset = (id: string, name: string, description: string, filters: StockFilters): void => {
  const presets = loadFilterPresets();
  const index = presets.findIndex(p => p.id === id);
  
  if (index !== -1) {
    presets[index] = {
      ...presets[index],
      name,
      description,
      filters,
      updatedAt: new Date().toISOString(),
    };
    
    try {
      localStorage.setItem(FILTER_PRESETS_KEY, JSON.stringify(presets));
    } catch (error) {
      console.error('Failed to update filter preset:', error);
    }
  }
};

/**
 * Get a specific preset by ID
 */
export const getFilterPreset = (id: string): FilterPreset | null => {
  const presets = loadFilterPresets();
  return presets.find(p => p.id === id) || null;
};

/**
 * Generate a unique ID for presets
 */
function generateId(): string {
  return `preset_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
}


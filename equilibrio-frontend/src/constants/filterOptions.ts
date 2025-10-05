/**
 * Constants for filter options used throughout the application
 */

export const EQUILIBRIUM_ZONES = ['discount', 'equilibrium', 'premium'] as const;

export const SIGNALS = ['buy', 'hold', 'sell'] as const;

export const VOLUME_PROFILES = ['high', 'medium', 'low'] as const;

export const TRENDS = ['bullish', 'neutral', 'bearish'] as const;

export const EQUILIBRIUM_ZONE_LABELS: Record<string, string> = {
  discount: 'Discount',
  equilibrium: 'Equilibrium',
  premium: 'Premium',
};

export const SIGNAL_LABELS: Record<string, string> = {
  buy: 'Buy',
  hold: 'Hold',
  sell: 'Sell',
};

export const VOLUME_PROFILE_LABELS: Record<string, string> = {
  high: 'High',
  medium: 'Medium',
  low: 'Low',
};

export const TREND_LABELS: Record<string, string> = {
  bullish: 'Bullish',
  neutral: 'Neutral',
  bearish: 'Bearish',
};

export type EquilibriumZone = typeof EQUILIBRIUM_ZONES[number];
export type Signal = typeof SIGNALS[number];
export type VolumeProfile = typeof VOLUME_PROFILES[number];
export type Trend = typeof TRENDS[number];


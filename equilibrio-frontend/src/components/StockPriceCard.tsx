import React from 'react';
import { TrendingUp, TrendingDown } from 'lucide-react';
import { StockData } from '../types';

interface StockPriceCardProps {
  stock: StockData | null;
  onClose: () => void;
}

const StockPriceCard: React.FC<StockPriceCardProps> = ({ stock, onClose }) => {
  if (!stock) return null;

  const isPositive = stock.changePercent >= 0;
  const changeColor = isPositive ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400';
  const bgColor = isPositive ? 'bg-green-50 dark:bg-green-900/20' : 'bg-red-50 dark:bg-red-900/20';
  const TrendIcon = isPositive ? TrendingUp : TrendingDown;

  return (
    <div className="bg-white dark:bg-slate-800 rounded-lg shadow-lg p-6 mb-6 transition-colors border-l-4 border-blue-600">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          {/* Symbol and Name */}
          <div className="mb-4">
            <div className="flex items-center gap-3">
              <h2 className="text-3xl font-bold text-slate-800 dark:text-slate-100">
                {stock.symbol}
              </h2>
              <span className={`px-3 py-1 rounded-full text-xs font-semibold ${bgColor} ${changeColor}`}>
                {stock.signal.toUpperCase()}
              </span>
            </div>
            <p className="text-sm text-slate-600 dark:text-slate-400 mt-1">{stock.name}</p>
            <p className="text-xs text-slate-500 dark:text-slate-500 mt-0.5">
              {stock.sector} • {stock.industry}
            </p>
          </div>

          {/* Price Display */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-4">
            {/* Current Price */}
            <div>
              <p className="text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase mb-1">
                Current Price
              </p>
              <div className="flex items-baseline gap-2">
                <span className="text-4xl font-bold text-slate-900 dark:text-slate-100">
                  ${stock.price.toFixed(2)}
                </span>
              </div>
            </div>

            {/* Change */}
            <div>
              <p className="text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase mb-1">
                Today's Change
              </p>
              <div className={`flex items-center gap-2 ${changeColor}`}>
                <TrendIcon className="w-6 h-6" />
                <div>
                  <div className="text-2xl font-bold">
                    {isPositive ? '+' : ''}{stock.change.toFixed(2)}
                  </div>
                  <div className="text-lg font-semibold">
                    ({isPositive ? '+' : ''}{stock.changePercent.toFixed(2)}%)
                  </div>
                </div>
              </div>
            </div>

            {/* Volume */}
            <div>
              <p className="text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase mb-1">
                Volume
              </p>
              <div className="text-2xl font-bold text-slate-900 dark:text-slate-100">
                {(stock.volume / 1000000).toFixed(2)}M
              </div>
              <p className="text-xs text-slate-600 dark:text-slate-400 mt-1">
                Profile: {stock.volumeProfile}
              </p>
            </div>
          </div>

          {/* Key Stats Row */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t border-slate-200 dark:border-slate-700">
            <div>
              <p className="text-xs text-slate-600 dark:text-slate-400">RSI</p>
              <p className="text-lg font-semibold text-slate-900 dark:text-slate-100">
                {stock.rsi.toFixed(2)}
              </p>
            </div>
            <div>
              <p className="text-xs text-slate-600 dark:text-slate-400">52W High</p>
              <p className="text-lg font-semibold text-slate-900 dark:text-slate-100">
                ${stock.distanceFrom52WeekHigh.toFixed(2)}%
              </p>
            </div>
            <div>
              <p className="text-xs text-slate-600 dark:text-slate-400">52W Low</p>
              <p className="text-lg font-semibold text-slate-900 dark:text-slate-100">
                ${stock.distanceFrom52WeekLow.toFixed(2)}%
              </p>
            </div>
            <div>
              <p className="text-xs text-slate-600 dark:text-slate-400">Trend</p>
              <p className={`text-lg font-semibold ${
                stock.trend === 'bullish' ? 'text-green-600 dark:text-green-400' :
                stock.trend === 'bearish' ? 'text-red-600 dark:text-red-400' :
                'text-slate-600 dark:text-slate-400'
              }`}>
                {stock.trend.charAt(0).toUpperCase() + stock.trend.slice(1)}
              </p>
            </div>
          </div>
        </div>

        {/* Close Button */}
        <button
          onClick={onClose}
          className="ml-4 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 transition-colors"
          aria-label="Close price card"
        >
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  );
};

export default React.memo(StockPriceCard);


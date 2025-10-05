import React from 'react';
import { StockData } from '../../types';
import {
  getEquilibriumTextColor,
  getVolumeProfileColor,
  formatPrice,
  formatPercent,
  formatVolume,
  formatMarketCap,
} from '../../utils/stockUtils';
import DataRow from '../ui/DataRow';

interface StockDetailsProps {
  stock: StockData;
}

const StockDetails: React.FC<StockDetailsProps> = ({ stock }) => {
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {/* Technical Indicators */}
      <div>
        <h4 className="text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase mb-2">
          Technical Indicators
        </h4>
        <div className="space-y-1 text-sm text-slate-700 dark:text-slate-300">
          <DataRow label="Stoch RSI" value={stock.stochRsi.toFixed(1)} />
          <DataRow label="Historic RSI" value={stock.historicRsiAvg.toFixed(1)} />
          <DataRow 
            label="MACD" 
            value={stock.macd.toFixed(2)}
            valueClassName={stock.macdHistogram > 0 ? 'text-green-600' : 'text-red-600'}
          />
          <DataRow label="MACD Signal" value={stock.macdSignal.toFixed(2)} />
        </div>
      </div>

      {/* Moving Averages */}
      <div>
        <h4 className="text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase mb-2">
          Moving Averages
        </h4>
        <div className="space-y-1 text-sm text-slate-700 dark:text-slate-300">
          <DataRow label="SMA 50" value={formatPrice(stock.sma50)} />
          <DataRow label="SMA 200" value={formatPrice(stock.sma200)} />
          <DataRow label="EMA 20" value={formatPrice(stock.ema20)} />
          <DataRow 
            label="Volume Profile" 
            value={stock.volumeProfile.toUpperCase()}
            valueClassName={getVolumeProfileColor(stock.volumeProfile)}
          />
        </div>
      </div>

      {/* Equilibrium Analysis */}
      <div>
        <h4 className="text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase mb-2">
          Equilibrium Analysis
        </h4>
        <div className="space-y-1 text-sm text-slate-700 dark:text-slate-300">
          <DataRow label="Equilibrium" value={formatPrice(stock.equilibriumLevel)} />
          <DataRow 
            label="Distance" 
            value={formatPercent(stock.priceToEquilibrium)}
            valueClassName={getEquilibriumTextColor(stock.priceToEquilibrium)}
          />
          <DataRow 
            label="From 52W High" 
            value={formatPercent(stock.distanceFrom52WeekHigh)}
            valueClassName="text-red-600"
          />
          <DataRow 
            label="From 52W Low" 
            value={formatPercent(stock.distanceFrom52WeekLow)}
            valueClassName="text-green-600"
          />
        </div>
      </div>

      {/* Market Data */}
      <div>
        <h4 className="text-xs font-semibold text-slate-600 dark:text-slate-400 uppercase mb-2">
          Market Data
        </h4>
        <div className="space-y-1 text-sm text-slate-700 dark:text-slate-300">
          <DataRow label="Volume" value={formatVolume(stock.volume)} />
          <DataRow label="Market Cap" value={formatMarketCap(stock.marketCap)} />
          <DataRow 
            label="Industry" 
            value={stock.industry}
            valueClassName="text-xs"
          />
        </div>
      </div>
    </div>
  );
};

export default React.memo(StockDetails);


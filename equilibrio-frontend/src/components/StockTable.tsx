import React from 'react';
import { TrendingUp, TrendingDown, ChevronDown, ChevronUp } from 'lucide-react';
import { StockData } from '../types';

interface StockTableProps {
  stocks: StockData[];
  loading: boolean;
  sortField: string;
  sortDirection: 'asc' | 'desc';
  onSort: (field: string) => void;
  onRowExpand: (symbol: string | null) => void;
  expandedRow: string | null;
}

const StockTable: React.FC<StockTableProps> = ({
  stocks,
  loading,
  sortField,
  sortDirection,
  onSort,
  onRowExpand,
  expandedRow,
}) => {
  const getEquilibriumColor = (value: number) => {
    if (value < -5) return 'text-green-600 bg-green-50';
    if (value > 5) return 'text-red-600 bg-red-50';
    return 'text-yellow-600 bg-yellow-50';
  };

  const getEquilibriumZone = (value: number) => {
    if (value < -5) return 'Discount';
    if (value > 5) return 'Premium';
    return 'Equilibrium';
  };

  const renderSortIcon = (field: string) => {
    if (sortField !== field) return null;
    return sortDirection === 'asc' ? '↑' : '↓';
  };

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-lg overflow-hidden">
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span className="ml-2 text-slate-600">Loading stocks...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden">
      <div className="px-6 py-4 bg-slate-50 border-b border-slate-200">
        <p className="text-sm text-slate-600">
          Showing <span className="font-semibold text-slate-800">{stocks.length}</span> stocks
        </p>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-slate-100 border-b border-slate-200">
            <tr>
              <th 
                className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200" 
                onClick={() => onSort('symbol')}
              >
                Symbol {renderSortIcon('symbol')}
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">
                Name
              </th>
              <th 
                className="px-4 py-3 text-right text-xs font-semibold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200" 
                onClick={() => onSort('price')}
              >
                Price {renderSortIcon('price')}
              </th>
              <th 
                className="px-4 py-3 text-right text-xs font-semibold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200" 
                onClick={() => onSort('changePercent')}
              >
                Change {renderSortIcon('changePercent')}
              </th>
              <th 
                className="px-4 py-3 text-right text-xs font-semibold text-slate-600 uppercase tracking-wider cursor-pointer hover:bg-slate-200" 
                onClick={() => onSort('rsi')}
              >
                RSI {renderSortIcon('rsi')}
              </th>
              <th className="px-4 py-3 text-center text-xs font-semibold text-slate-600 uppercase tracking-wider">
                Equilibrium
              </th>
              <th className="px-4 py-3 text-center text-xs font-semibold text-slate-600 uppercase tracking-wider">
                Trend
              </th>
              <th className="px-4 py-3 text-center text-xs font-semibold text-slate-600 uppercase tracking-wider">
                Signal
              </th>
              <th className="px-4 py-3 text-left text-xs font-semibold text-slate-600 uppercase tracking-wider">
                Sector
              </th>
              <th className="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200">
            {stocks.map((stock) => (
              <React.Fragment key={stock.symbol}>
                <tr className="hover:bg-slate-50 transition">
                  <td className="px-4 py-3 whitespace-nowrap">
                    <span className="font-semibold text-slate-800">{stock.symbol}</span>
                  </td>
                  <td className="px-4 py-3">
                    <div className="text-sm text-slate-700">{stock.name}</div>
                    <div className="text-xs text-slate-500">{stock.industry}</div>
                  </td>
                  <td className="px-4 py-3 text-right whitespace-nowrap">
                    <span className="font-medium text-slate-800">${stock.price.toFixed(2)}</span>
                  </td>
                  <td className="px-4 py-3 text-right whitespace-nowrap">
                    <div className={`flex items-center justify-end gap-1 ${stock.changePercent >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {stock.changePercent >= 0 ? <TrendingUp className="w-4 h-4" /> : <TrendingDown className="w-4 h-4" />}
                      <span className="font-medium">{stock.changePercent.toFixed(2)}%</span>
                    </div>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <span className={`font-medium ${
                      stock.rsi < 30 ? 'text-green-600' : 
                      stock.rsi > 70 ? 'text-red-600' : 
                      'text-slate-700'
                    }`}>
                      {stock.rsi.toFixed(1)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-center">
                    <span className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${getEquilibriumColor(stock.priceToEquilibrium)}`}>
                      {getEquilibriumZone(stock.priceToEquilibrium)}
                      <div className="text-xs mt-0.5">{stock.priceToEquilibrium.toFixed(1)}%</div>
                    </span>
                  </td>
                  <td className="px-4 py-3 text-center">
                    <span className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${
                      stock.trend === 'bullish' ? 'bg-green-100 text-green-800' :
                      stock.trend === 'bearish' ? 'bg-red-100 text-red-800' :
                      'bg-slate-100 text-slate-800'
                    }`}>
                      {stock.trend}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-center">
                    <span className={`inline-block px-2 py-1 rounded-full text-xs font-bold uppercase ${
                      stock.signal === 'buy' ? 'bg-green-600 text-white' :
                      stock.signal === 'sell' ? 'bg-red-600 text-white' :
                      'bg-slate-300 text-slate-700'
                    }`}>
                      {stock.signal}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <span className="text-sm text-slate-700">{stock.sector}</span>
                  </td>
                  <td className="px-4 py-3">
                    <button
                      onClick={() => onRowExpand(expandedRow === stock.symbol ? null : stock.symbol)}
                      className="text-blue-600 hover:text-blue-800"
                    >
                      {expandedRow === stock.symbol ? <ChevronUp className="w-5 h-5" /> : <ChevronDown className="w-5 h-5" />}
                    </button>
                  </td>
                </tr>
                
                {/* Expanded Details Row */}
                {expandedRow === stock.symbol && (
                  <tr className="bg-slate-50">
                    <td colSpan={10} className="px-4 py-4">
                      <StockDetails stock={stock} />
                    </td>
                  </tr>
                )}
              </React.Fragment>
            ))}
          </tbody>
        </table>
      </div>

      {stocks.length === 0 && (
        <div className="text-center py-12">
          <p className="text-slate-500">No stocks match your current filters</p>
        </div>
      )}
    </div>
  );
};

// StockDetails component for expanded row content
const StockDetails: React.FC<{ stock: StockData }> = ({ stock }) => {
  const getEquilibriumColor = (value: number) => {
    if (value < -5) return 'text-green-600';
    if (value > 5) return 'text-red-600';
    return 'text-yellow-600';
  };

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div>
        <h4 className="text-xs font-semibold text-slate-600 uppercase mb-2">Technical Indicators</h4>
        <div className="space-y-1 text-sm">
          <div className="flex justify-between">
            <span className="text-slate-600">Stoch RSI:</span>
            <span className="font-medium">{stock.stochRsi.toFixed(1)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">Historic RSI:</span>
            <span className="font-medium">{stock.historicRsiAvg.toFixed(1)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">MACD:</span>
            <span className={`font-medium ${stock.macdHistogram > 0 ? 'text-green-600' : 'text-red-600'}`}>
              {stock.macd.toFixed(2)}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">MACD Signal:</span>
            <span className="font-medium">{stock.macdSignal.toFixed(2)}</span>
          </div>
        </div>
      </div>

      <div>
        <h4 className="text-xs font-semibold text-slate-600 uppercase mb-2">Moving Averages</h4>
        <div className="space-y-1 text-sm">
          <div className="flex justify-between">
            <span className="text-slate-600">SMA 50:</span>
            <span className="font-medium">${stock.sma50.toFixed(2)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">SMA 200:</span>
            <span className="font-medium">${stock.sma200.toFixed(2)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">EMA 20:</span>
            <span className="font-medium">${stock.ema20.toFixed(2)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">Volume Profile:</span>
            <span className={`font-medium uppercase ${
              stock.volumeProfile === 'high' ? 'text-green-600' :
              stock.volumeProfile === 'low' ? 'text-red-600' :
              'text-yellow-600'
            }`}>
              {stock.volumeProfile}
            </span>
          </div>
        </div>
      </div>

      <div>
        <h4 className="text-xs font-semibold text-slate-600 uppercase mb-2">Equilibrium Analysis</h4>
        <div className="space-y-1 text-sm">
          <div className="flex justify-between">
            <span className="text-slate-600">Equilibrium:</span>
            <span className="font-medium">${stock.equilibriumLevel.toFixed(2)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">Distance:</span>
            <span className={`font-medium ${getEquilibriumColor(stock.priceToEquilibrium)}`}>
              {stock.priceToEquilibrium.toFixed(1)}%
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">From 52W High:</span>
            <span className="font-medium text-red-600">{stock.distanceFrom52WeekHigh.toFixed(1)}%</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">From 52W Low:</span>
            <span className="font-medium text-green-600">{stock.distanceFrom52WeekLow.toFixed(1)}%</span>
          </div>
        </div>
      </div>

      <div>
        <h4 className="text-xs font-semibold text-slate-600 uppercase mb-2">Market Data</h4>
        <div className="space-y-1 text-sm">
          <div className="flex justify-between">
            <span className="text-slate-600">Volume:</span>
            <span className="font-medium">{(stock.volume / 1000000).toFixed(2)}M</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">Market Cap:</span>
            <span className="font-medium">${(stock.marketCap / 1000000000).toFixed(2)}B</span>
          </div>
          <div className="flex justify-between">
            <span className="text-slate-600">Industry:</span>
            <span className="font-medium text-xs">{stock.industry}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default StockTable;

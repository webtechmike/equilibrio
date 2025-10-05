import React from 'react';
import { TrendingUp, TrendingDown, ChevronDown, ChevronUp } from 'lucide-react';
import { StockData } from '../../types';
import { 
  getRSIColor, 
  getEquilibriumColor, 
  getEquilibriumZone,
  getTrendColor,
  getSignalColor,
  formatPrice,
  formatPercent,
  getChangeColor
} from '../../utils/stockUtils';
import Badge from '../ui/Badge';

interface StockRowProps {
  stock: StockData;
  isExpanded: boolean;
  onToggleExpand: () => void;
  onStockClick: () => void;
}

const StockRow: React.FC<StockRowProps> = ({ stock, isExpanded, onToggleExpand, onStockClick }) => {
  return (
    <tr className="hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors cursor-pointer" onClick={onStockClick}>
      <td className="px-4 py-3 whitespace-nowrap">
        <span className="font-semibold text-slate-800 dark:text-slate-200">{stock.symbol}</span>
      </td>
      
      <td className="px-4 py-3">
        <div className="text-sm text-slate-700 dark:text-slate-300">{stock.name}</div>
        <div className="text-xs text-slate-500 dark:text-slate-500">{stock.industry}</div>
      </td>
      
      <td className="px-4 py-3 text-right whitespace-nowrap">
        <span className="font-medium text-slate-800 dark:text-slate-200">{formatPrice(stock.price)}</span>
      </td>
      
      <td className="px-4 py-3 text-right whitespace-nowrap">
        <div className={`flex items-center justify-end gap-1 ${getChangeColor(stock.changePercent)}`}>
          {stock.changePercent >= 0 ? (
            <TrendingUp className="w-4 h-4" />
          ) : (
            <TrendingDown className="w-4 h-4" />
          )}
          <span className="font-medium">{formatPercent(stock.changePercent)}</span>
        </div>
      </td>
      
      <td className="px-4 py-3 text-right">
        <span className={`font-medium ${getRSIColor(stock.rsi)}`}>
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
        <Badge className={getTrendColor(stock.trend)}>
          {stock.trend}
        </Badge>
      </td>
      
      <td className="px-4 py-3 text-center">
        <Badge className={`${getSignalColor(stock.signal)} uppercase font-bold`}>
          {stock.signal}
        </Badge>
      </td>
      
      <td className="px-4 py-3">
        <span className="text-sm text-slate-700 dark:text-slate-300">{stock.sector}</span>
      </td>
      
      <td className="px-4 py-3">
        <button
          onClick={(e) => {
            e.stopPropagation();
            onToggleExpand();
          }}
          className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 transition-colors"
        >
          {isExpanded ? (
            <ChevronUp className="w-5 h-5" />
          ) : (
            <ChevronDown className="w-5 h-5" />
          )}
        </button>
      </td>
    </tr>
  );
};

export default React.memo(StockRow);


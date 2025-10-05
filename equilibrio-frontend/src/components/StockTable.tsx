import React from 'react';
import { StockData } from '../types';
import TableHeader from './table/TableHeader';
import StockRow from './table/StockRow';
import StockDetails from './table/StockDetails';

interface StockTableProps {
  stocks: StockData[];
  loading: boolean;
  sortField: string;
  sortDirection: 'asc' | 'desc';
  onSort: (field: string) => void;
  onRowExpand: (symbol: string | null) => void;
  expandedRow: string | null;
  onStockClick: (stock: StockData) => void;
}

const StockTable: React.FC<StockTableProps> = ({
  stocks,
  loading,
  sortField,
  sortDirection,
  onSort,
  onRowExpand,
  expandedRow,
  onStockClick,
}) => {

  if (loading) {
    return (
      <div className="bg-white dark:bg-slate-800 rounded-lg shadow-lg overflow-hidden transition-colors">
        <div className="flex items-center justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 dark:border-blue-400"></div>
          <span className="ml-2 text-slate-600 dark:text-slate-400">Loading stocks...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-slate-800 rounded-lg shadow-lg overflow-hidden transition-colors">
      <div className="px-6 py-4 bg-slate-50 dark:bg-slate-700/50 border-b border-slate-200 dark:border-slate-700">
        <p className="text-sm text-slate-600 dark:text-slate-400">
          Showing <span className="font-semibold text-slate-800 dark:text-slate-200">{stocks.length}</span> stocks
        </p>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-slate-100 dark:bg-slate-700/50 border-b border-slate-200 dark:border-slate-700">
            <tr>
              <TableHeader 
                label="Symbol" 
                field="symbol" 
                sortable 
                currentSortField={sortField}
                sortDirection={sortDirection}
                onSort={onSort}
              />
              <TableHeader label="Name" />
              <TableHeader 
                label="Price" 
                field="price" 
                sortable 
                currentSortField={sortField}
                sortDirection={sortDirection}
                align="right"
                onSort={onSort}
              />
              <TableHeader 
                label="Change" 
                field="changePercent" 
                sortable 
                currentSortField={sortField}
                sortDirection={sortDirection}
                align="right"
                onSort={onSort}
              />
              <TableHeader 
                label="RSI" 
                field="rsi" 
                sortable 
                currentSortField={sortField}
                sortDirection={sortDirection}
                align="right"
                onSort={onSort}
              />
              <TableHeader label="Equilibrium" align="center" />
              <TableHeader label="Trend" align="center" />
              <TableHeader label="Signal" align="center" />
              <TableHeader label="Sector" />
              <th className="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-200 dark:divide-slate-700">
            {stocks.map((stock) => (
              <React.Fragment key={stock.symbol}>
                <StockRow 
                  stock={stock}
                  isExpanded={expandedRow === stock.symbol}
                  onToggleExpand={() => onRowExpand(expandedRow === stock.symbol ? null : stock.symbol)}
                  onStockClick={() => onStockClick(stock)}
                />
                
                {/* Expanded Details Row */}
                {expandedRow === stock.symbol && (
                  <tr className="bg-slate-50 dark:bg-slate-700/30">
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
          <p className="text-slate-500 dark:text-slate-400">No stocks match your current filters</p>
        </div>
      )}
    </div>
  );
};

export default React.memo(StockTable);

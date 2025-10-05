import React from 'react';

const EquilibriumInfo: React.FC = () => {
  return (
    <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
      <h3 className="text-sm font-semibold text-blue-900 mb-2">📊 About Equilibrium Trading</h3>
      <p className="text-sm text-blue-800">
        <strong>Equilibrium</strong> represents the 50% retracement level between a stock's 52-week high and low. 
        The <strong>Discount Zone</strong> (below -5%) is where smart money typically accumulates positions, 
        while the <strong>Premium Zone</strong> (above +5%) is where distribution often occurs. 
        Combine equilibrium analysis with RSI, MACD, and trend indicators for optimal swing trade entries.
      </p>
    </div>
  );
};

export default EquilibriumInfo;

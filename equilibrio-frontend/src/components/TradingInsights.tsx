import React from 'react';
import { StockData } from '../types';

interface TradingInsightsProps {
  stock: StockData;
}

const TradingInsights: React.FC<TradingInsightsProps> = ({ stock }) => {

  const renderInsights = () => {
    const insights = [];

    // Strong Buy Setup
    if (stock.signal === 'buy' && stock.priceToEquilibrium < -10) {
      insights.push(
        <p key="strong-buy">
          💡 <strong>Strong Buy Setup:</strong> Price is in discount zone ({stock.priceToEquilibrium.toFixed(1)}% below equilibrium) with oversold RSI ({stock.rsi.toFixed(1)}). Consider entry for swing trade.
        </p>
      );
    }

    // Potential Exit
    if (stock.signal === 'sell' && stock.priceToEquilibrium > 10) {
      insights.push(
        <p key="potential-exit">
          ⚠️ <strong>Potential Exit:</strong> Price is in premium zone ({stock.priceToEquilibrium.toFixed(1)}% above equilibrium) with overbought RSI ({stock.rsi.toFixed(1)}). Consider taking profits.
        </p>
      );
    }

    // At Equilibrium
    if (stock.signal === 'hold' && Math.abs(stock.priceToEquilibrium) < 5) {
      insights.push(
        <p key="at-equilibrium">
          ⚖️ <strong>At Equilibrium:</strong> Price is near the 50% retracement level. Wait for clear directional move before entering.
        </p>
      );
    }

    // Strong Uptrend
    if (stock.trend === 'bullish' && stock.price > stock.sma50 && stock.sma50 > stock.sma200) {
      insights.push(
        <p key="strong-uptrend" className="mt-2">
          📈 <strong>Strong Uptrend:</strong> Price above SMA50 (${stock.sma50.toFixed(2)}) and SMA50 above SMA200 (${stock.sma200.toFixed(2)}). Trend is intact.
        </p>
      );
    }

    // Strong Downtrend
    if (stock.trend === 'bearish' && stock.price < stock.sma50 && stock.sma50 < stock.sma200) {
      insights.push(
        <p key="strong-downtrend" className="mt-2">
          📉 <strong>Strong Downtrend:</strong> Price below SMA50 and SMA50 below SMA200. Consider short setups or avoid longs.
        </p>
      );
    }

    // MACD Bullish
    if (stock.macdHistogram > 0 && stock.macd > stock.macdSignal) {
      insights.push(
        <p key="macd-bullish" className="mt-2">
          ✅ <strong>MACD Bullish:</strong> MACD ({stock.macd.toFixed(2)}) above signal line ({stock.macdSignal.toFixed(2)}). Momentum is positive.
        </p>
      );
    }

    // MACD Bearish
    if (stock.macdHistogram < 0 && stock.macd < stock.macdSignal) {
      insights.push(
        <p key="macd-bearish" className="mt-2">
          ❌ <strong>MACD Bearish:</strong> MACD below signal line. Momentum is negative.
        </p>
      );
    }

    return insights;
  };

  const insights = renderInsights();

  if (insights.length === 0) {
    return null;
  }

  return (
    <div className="mt-4 p-3 bg-blue-50 rounded-lg border border-blue-200">
      <h4 className="text-sm font-semibold text-blue-900 mb-2">Trading Insights</h4>
      <div className="text-sm text-blue-800">
        {insights}
      </div>
    </div>
  );
};

export default TradingInsights;

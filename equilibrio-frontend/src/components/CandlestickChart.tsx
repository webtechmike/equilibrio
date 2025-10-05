import React, { useEffect, useRef, useState } from 'react';
import { createChart, ColorType, IChartApi, ISeriesApi } from 'lightweight-charts';
import { X } from 'lucide-react';

interface CandlestickData {
  time: string;
  open: number;
  high: number;
  low: number;
  close: number;
}

interface CandlestickChartProps {
  symbol: string;
  companyName: string;
  data: CandlestickData[];
  onClose: () => void;
  onTimeframeChange: (days: number) => void;
}

type Timeframe = {
  label: string;
  days: number;
};

const CandlestickChart: React.FC<CandlestickChartProps> = ({
  symbol,
  companyName,
  data,
  onClose,
  onTimeframeChange,
}) => {
  const chartContainerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);
  const seriesRef = useRef<ISeriesApi<'Candlestick'> | null>(null);
  const [selectedTimeframe, setSelectedTimeframe] = useState<number>(90);

  const timeframes: Timeframe[] = [
    { label: '1W', days: 7 },
    { label: '1M', days: 30 },
    { label: '3M', days: 90 },
    { label: '6M', days: 180 },
    { label: '1Y', days: 365 },
  ];

  const handleTimeframeChange = (days: number) => {
    setSelectedTimeframe(days);
    onTimeframeChange(days);
  };

  useEffect(() => {
    if (!chartContainerRef.current || data.length === 0) return;

    // Create chart with dark theme
    const chart = createChart(chartContainerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: '#1e293b' }, // slate-800
        textColor: '#cbd5e1', // slate-300
      },
      grid: {
        vertLines: { color: '#334155' }, // slate-700
        horzLines: { color: '#334155' }, // slate-700
      },
      width: chartContainerRef.current.clientWidth,
      height: 400,
      timeScale: {
        borderColor: '#475569', // slate-600
        timeVisible: true,
      },
      rightPriceScale: {
        borderColor: '#475569', // slate-600
      },
    });

    // Create candlestick series with green/red colors
    const candlestickSeries = chart.addCandlestickSeries({
      upColor: '#22c55e', // green-500
      downColor: '#ef4444', // red-500
      borderUpColor: '#16a34a', // green-600
      borderDownColor: '#dc2626', // red-600
      wickUpColor: '#22c55e', // green-500
      wickDownColor: '#ef4444', // red-500
    });

    // Set data
    candlestickSeries.setData(data);

    // Fit content to view
    chart.timeScale().fitContent();

    chartRef.current = chart;
    seriesRef.current = candlestickSeries;

    // Handle resize
    const handleResize = () => {
      if (chartContainerRef.current && chartRef.current) {
        chartRef.current.applyOptions({
          width: chartContainerRef.current.clientWidth,
        });
      }
    };

    window.addEventListener('resize', handleResize);

    // Cleanup
    return () => {
      window.removeEventListener('resize', handleResize);
      if (chartRef.current) {
        chartRef.current.remove();
      }
    };
  }, [data]);

  if (data.length === 0) {
    return (
      <div className="bg-slate-800 dark:bg-slate-900 rounded-lg shadow-lg p-6 mb-6 transition-colors">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-xl font-bold text-slate-100">
              {symbol} - {companyName}
            </h2>
            <p className="text-slate-400 text-sm">Daily Candlestick Chart</p>
          </div>
          <button
            onClick={onClose}
            className="p-2 rounded-lg hover:bg-slate-700 transition-colors text-slate-400 hover:text-slate-200"
          >
            <X className="w-5 h-5" />
          </button>
        </div>
        <div className="flex items-center justify-center h-64 text-slate-400">
          Loading chart data...
        </div>
      </div>
    );
  }

  return (
    <div className="bg-slate-800 dark:bg-slate-900 rounded-lg shadow-lg p-6 mb-6 transition-colors">
      <div className="flex items-center justify-between mb-4">
        <div className="flex-1">
          <h2 className="text-xl font-bold text-slate-100">
            {symbol} - {companyName}
          </h2>
          <p className="text-slate-400 text-sm">Daily Candlestick Chart</p>
        </div>
        
        {/* Timeframe Controls */}
        <div className="flex items-center gap-2 mx-4">
          {timeframes.map((tf) => (
            <button
              key={tf.days}
              onClick={() => handleTimeframeChange(tf.days)}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                selectedTimeframe === tf.days
                  ? 'bg-blue-600 text-white'
                  : 'bg-slate-700 text-slate-300 hover:bg-slate-600'
              }`}
            >
              {tf.label}
            </button>
          ))}
        </div>

        <button
          onClick={onClose}
          className="p-2 rounded-lg hover:bg-slate-700 transition-colors text-slate-400 hover:text-slate-200"
          aria-label="Close chart"
        >
          <X className="w-5 h-5" />
        </button>
      </div>
      <div ref={chartContainerRef} className="w-full" />
    </div>
  );
};

export default React.memo(CandlestickChart);


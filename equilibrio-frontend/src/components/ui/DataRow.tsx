import React from 'react';

interface DataRowProps {
  label: string;
  value: string | React.ReactNode;
  valueClassName?: string;
}

const DataRow: React.FC<DataRowProps> = ({ label, value, valueClassName = '' }) => {
  return (
    <div className="flex justify-between">
      <span className="text-slate-600 dark:text-slate-400">{label}:</span>
      <span className={`font-medium ${valueClassName || 'text-slate-900 dark:text-slate-100'}`}>{value}</span>
    </div>
  );
};

export default DataRow;


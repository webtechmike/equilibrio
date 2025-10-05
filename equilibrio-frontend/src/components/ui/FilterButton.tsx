import React from 'react';

interface FilterButtonProps {
  label: string;
  active: boolean;
  onClick: () => void;
}

const FilterButton: React.FC<FilterButtonProps> = ({ label, active, onClick }) => {
  return (
    <button
      onClick={onClick}
      className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
        active
          ? 'bg-blue-600 dark:bg-blue-500 text-white'
          : 'bg-slate-200 dark:bg-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-300 dark:hover:bg-slate-600'
      }`}
    >
      {label}
    </button>
  );
};

export default FilterButton;


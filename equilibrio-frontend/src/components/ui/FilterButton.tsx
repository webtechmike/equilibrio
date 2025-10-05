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
      className={`px-3 py-1 rounded-full text-sm font-medium transition ${
        active
          ? 'bg-blue-600 text-white'
          : 'bg-slate-200 text-slate-700 hover:bg-slate-300'
      }`}
    >
      {label}
    </button>
  );
};

export default FilterButton;


import React from 'react';

interface TableHeaderProps {
  label: string;
  field?: string;
  sortable?: boolean;
  currentSortField?: string;
  sortDirection?: 'asc' | 'desc';
  align?: 'left' | 'center' | 'right';
  onSort?: (field: string) => void;
}

const TableHeader: React.FC<TableHeaderProps> = ({
  label,
  field,
  sortable = false,
  currentSortField,
  sortDirection,
  align = 'left',
  onSort,
}) => {
  const alignmentClasses = {
    left: 'text-left',
    center: 'text-center',
    right: 'text-right',
  };

  const handleClick = () => {
    if (sortable && field && onSort) {
      onSort(field);
    }
  };

  const showSortIcon = sortable && field === currentSortField;
  const sortIcon = sortDirection === 'asc' ? '↑' : '↓';

  return (
    <th
      className={`px-4 py-3 ${alignmentClasses[align]} text-xs font-semibold text-slate-600 uppercase tracking-wider ${
        sortable ? 'cursor-pointer hover:bg-slate-200' : ''
      }`}
      onClick={handleClick}
    >
      {label} {showSortIcon && sortIcon}
    </th>
  );
};

export default TableHeader;


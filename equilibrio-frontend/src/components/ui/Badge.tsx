import React from 'react';

interface BadgeProps {
  children: React.ReactNode;
  variant?: 'default' | 'success' | 'danger' | 'warning' | 'info';
  size?: 'sm' | 'md';
  className?: string;
}

const Badge: React.FC<BadgeProps> = ({ 
  children, 
  variant = 'default', 
  size = 'sm',
  className = '' 
}) => {
  const baseClasses = 'inline-block rounded-full font-medium';
  
  const sizeClasses = {
    sm: 'px-2 py-1 text-xs',
    md: 'px-3 py-1.5 text-sm',
  };
  
  const variantClasses = {
    default: 'bg-slate-100 dark:bg-slate-700 text-slate-800 dark:text-slate-200',
    success: 'bg-green-100 text-green-800',
    danger: 'bg-red-100 text-red-800',
    warning: 'bg-yellow-50 text-yellow-600',
    info: 'bg-blue-100 text-blue-800',
  };
  
  // If custom className is provided, don't apply variant classes (allows full customization)
  const colorClasses = className ? '' : variantClasses[variant];
  
  return (
    <span className={`${baseClasses} ${sizeClasses[size]} ${colorClasses} ${className}`}>
      {children}
    </span>
  );
};

export default Badge;


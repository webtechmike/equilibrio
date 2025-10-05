import React, { useState, useEffect } from 'react';
import { Save, FolderOpen, Trash2, X } from 'lucide-react';
import {
  loadFilterPresets,
  saveFilterPreset,
  deleteFilterPreset,
  FilterPreset,
} from '../utils/filterStorage';
import { StockFilters } from '../types';

interface FilterPresetsProps {
  currentFilters: StockFilters;
  onLoadPreset: (filters: StockFilters) => void;
}

const FilterPresets: React.FC<FilterPresetsProps> = ({ currentFilters, onLoadPreset }) => {
  const [presets, setPresets] = useState<FilterPreset[]>([]);
  const [showSaveDialog, setShowSaveDialog] = useState(false);
  const [showLoadDialog, setShowLoadDialog] = useState(false);
  const [presetName, setPresetName] = useState('');
  const [presetDescription, setPresetDescription] = useState('');

  useEffect(() => {
    loadPresets();
  }, []);

  const loadPresets = () => {
    const loaded = loadFilterPresets();
    setPresets(loaded);
  };

  const handleSave = () => {
    if (presetName.trim()) {
      saveFilterPreset(presetName, presetDescription, currentFilters);
      setPresetName('');
      setPresetDescription('');
      setShowSaveDialog(false);
      loadPresets();
    }
  };

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this preset?')) {
      deleteFilterPreset(id);
      loadPresets();
    }
  };

  const handleLoad = (preset: FilterPreset) => {
    onLoadPreset(preset.filters);
    setShowLoadDialog(false);
  };

  return (
    <div className="flex items-center gap-2">
      {/* Save Preset Button */}
      <button
        onClick={() => setShowSaveDialog(true)}
        className="flex items-center gap-2 px-3 py-2 bg-green-600 dark:bg-green-500 text-white rounded-lg hover:bg-green-700 dark:hover:bg-green-600 transition-colors text-sm"
        title="Save current filter configuration"
      >
        <Save className="w-4 h-4" />
        <span>Save Filters</span>
      </button>

      {/* Load Preset Button */}
      <button
        onClick={() => setShowLoadDialog(true)}
        className="flex items-center gap-2 px-3 py-2 bg-blue-600 dark:bg-blue-500 text-white rounded-lg hover:bg-blue-700 dark:hover:bg-blue-600 transition-colors text-sm"
        title="Load saved filter configuration"
      >
        <FolderOpen className="w-4 h-4" />
        <span>Load Filters</span>
      </button>

      {/* Save Dialog */}
      {showSaveDialog && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-slate-800 rounded-lg shadow-xl p-6 max-w-md w-full mx-4">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-bold text-slate-800 dark:text-slate-100">
                Save Filter Preset
              </h3>
              <button
                onClick={() => setShowSaveDialog(false)}
                className="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
                  Preset Name *
                </label>
                <input
                  type="text"
                  value={presetName}
                  onChange={(e) => setPresetName(e.target.value)}
                  placeholder="e.g., Oversold Tech Stocks"
                  className="w-full px-3 py-2 border border-slate-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100 focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
                  Description (optional)
                </label>
                <textarea
                  value={presetDescription}
                  onChange={(e) => setPresetDescription(e.target.value)}
                  placeholder="Brief description of this filter configuration..."
                  rows={3}
                  className="w-full px-3 py-2 border border-slate-300 dark:border-slate-600 rounded-lg bg-white dark:bg-slate-700 text-slate-900 dark:text-slate-100 focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div className="flex justify-end gap-2 mt-6">
                <button
                  onClick={() => setShowSaveDialog(false)}
                  className="px-4 py-2 text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleSave}
                  disabled={!presetName.trim()}
                  className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Save Preset
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Load Dialog */}
      {showLoadDialog && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-slate-800 rounded-lg shadow-xl p-6 max-w-2xl w-full mx-4 max-h-[80vh] overflow-y-auto">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-bold text-slate-800 dark:text-slate-100">
                Load Filter Preset
              </h3>
              <button
                onClick={() => setShowLoadDialog(false)}
                className="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            {presets.length === 0 ? (
              <div className="text-center py-8 text-slate-500 dark:text-slate-400">
                No saved presets yet. Save your current filters to create a preset.
              </div>
            ) : (
              <div className="space-y-3">
                {presets.map((preset) => (
                  <div
                    key={preset.id}
                    className="border border-slate-200 dark:border-slate-700 rounded-lg p-4 hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors"
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <h4 className="font-semibold text-slate-800 dark:text-slate-100 mb-1">
                          {preset.name}
                        </h4>
                        {preset.description && (
                          <p className="text-sm text-slate-600 dark:text-slate-400 mb-2">
                            {preset.description}
                          </p>
                        )}
                        <p className="text-xs text-slate-500 dark:text-slate-500">
                          Created: {new Date(preset.createdAt).toLocaleDateString()}
                        </p>
                      </div>
                      <div className="flex gap-2 ml-4">
                        <button
                          onClick={() => handleLoad(preset)}
                          className="px-3 py-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm"
                        >
                          Load
                        </button>
                        <button
                          onClick={() => handleDelete(preset.id)}
                          className="p-1.5 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default FilterPresets;


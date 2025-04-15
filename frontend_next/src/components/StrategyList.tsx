'use client';

import React from 'react';
import { useStrategies } from '@/hooks/useStrategies';
import { formatDate } from '@/lib/utils';

export const StrategyList: React.FC = () => {
  const { getAll, loading, error } = useStrategies();
  const strategiesQuery = getAll();
  
  if (loading || strategiesQuery.isLoading) {
    return <div>Loading strategies...</div>;
  }

  if (error || strategiesQuery.error) {
    const errorMessage = strategiesQuery.error instanceof Error 
      ? strategiesQuery.error.message 
      : error?.message || 'Unknown error';
    return <div>Error loading strategies: {errorMessage}</div>;
  }

  const strategies = strategiesQuery.data || [];

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-bold">Strategies</h2>
      {strategies.length === 0 ? (
        <p className="text-gray-500">No strategies found.</p>
      ) : (
        <ul className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {strategies.map((strategy) => (
            <li key={strategy.id} className="border rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow">
              <h3 className="text-xl font-medium mb-2">{strategy.name}</h3>
              <p className="text-gray-700 mb-4">{strategy.description}</p>
              <div className="flex justify-between items-center text-sm">
                <span 
                  className={`px-2 py-1 rounded-full ${
                    strategy.isEnabled 
                      ? 'bg-green-100 text-green-800' 
                      : 'bg-gray-100 text-gray-800'
                  }`}
                >
                  {strategy.isEnabled ? 'Enabled' : 'Disabled'}
                </span>
                <span className="text-gray-500">
                  Updated: {formatDate(strategy.updatedAt)}
                </span>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default StrategyList; 
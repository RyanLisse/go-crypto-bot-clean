import React, { useEffect, useState } from 'react';
import { useRepository } from '../hooks/useRepository';
import { Strategy } from '../api/repository/strategy.repository';

export const StrategyList: React.FC = () => {
  const { getAll, loading, error } = useRepository<Strategy>('strategy');
  const [strategies, setStrategies] = useState<Strategy[]>([]);

  useEffect(() => {
    const fetchStrategies = async () => {
      try {
        const data = await getAll();
        setStrategies(data);
      } catch (err) {
        console.error('Failed to fetch strategies:', err);
      }
    };

    fetchStrategies();
  }, [getAll]);

  if (loading) {
    return <div>Loading strategies...</div>;
  }

  if (error) {
    return <div>Error loading strategies: {error.message}</div>;
  }

  return (
    <div>
      <h2>Strategies</h2>
      {strategies.length === 0 ? (
        <p>No strategies found.</p>
      ) : (
        <ul>
          {strategies.map((strategy) => (
            <li key={strategy.id}>
              <h3>{strategy.name}</h3>
              <p>{strategy.description}</p>
              <p>Status: {strategy.isEnabled ? 'Enabled' : 'Disabled'}</p>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default StrategyList;

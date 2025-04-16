import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './MexcApiTest.css';

const API_BASE_URL = 'http://localhost:8080/api/v1/mexc';

const MexcApiTest = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [selectedEndpoint, setSelectedEndpoint] = useState('ticker/BTCUSDT');
  const [responseData, setResponseData] = useState(null);
  const [symbol, setSymbol] = useState('BTCUSDT');
  const [interval, setInterval] = useState('1h');

  const endpoints = [
    { value: 'account', label: 'Account Information' },
    { value: `ticker/${symbol}`, label: 'Ticker' },
    { value: `orderbook/${symbol}`, label: 'Order Book' },
    { value: `klines/${symbol}/${interval}`, label: 'Klines (Candlestick)' },
    { value: 'exchange-info', label: 'Exchange Information' },
    { value: `symbol/${symbol}`, label: 'Symbol Information' },
    { value: 'new-listings', label: 'New Listings' },
  ];

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    setResponseData(null);

    try {
      const response = await axios.get(`${API_BASE_URL}/${selectedEndpoint}`);
      setResponseData(response.data);
    } catch (err) {
      setError(err.message || 'An error occurred');
      console.error('Error fetching data:', err);
    } finally {
      setLoading(false);
    }
  };

  // Update selected endpoint when symbol or interval changes
  useEffect(() => {
    if (selectedEndpoint.startsWith('ticker/')) {
      setSelectedEndpoint(`ticker/${symbol}`);
    } else if (selectedEndpoint.startsWith('orderbook/')) {
      setSelectedEndpoint(`orderbook/${symbol}`);
    } else if (selectedEndpoint.startsWith('klines/')) {
      setSelectedEndpoint(`klines/${symbol}/${interval}`);
    } else if (selectedEndpoint.startsWith('symbol/')) {
      setSelectedEndpoint(`symbol/${symbol}`);
    }
  }, [symbol, interval, selectedEndpoint]);

  return (
    <div className="mexc-api-test">
      <h1>MEXC API Test</h1>
      
      <div className="control-panel">
        <div className="form-group">
          <label htmlFor="endpoint">Endpoint:</label>
          <select 
            id="endpoint" 
            value={selectedEndpoint} 
            onChange={(e) => setSelectedEndpoint(e.target.value)}
          >
            {endpoints.map((endpoint) => (
              <option key={endpoint.value} value={endpoint.value}>
                {endpoint.label}
              </option>
            ))}
          </select>
        </div>
        
        {(selectedEndpoint.includes('ticker/') || 
          selectedEndpoint.includes('orderbook/') || 
          selectedEndpoint.includes('klines/') || 
          selectedEndpoint.includes('symbol/')) && (
          <div className="form-group">
            <label htmlFor="symbol">Symbol:</label>
            <input 
              id="symbol" 
              type="text" 
              value={symbol} 
              onChange={(e) => setSymbol(e.target.value)}
              placeholder="e.g., BTCUSDT"
            />
          </div>
        )}
        
        {selectedEndpoint.includes('klines/') && (
          <div className="form-group">
            <label htmlFor="interval">Interval:</label>
            <select 
              id="interval" 
              value={interval} 
              onChange={(e) => setInterval(e.target.value)}
            >
              <option value="1m">1 minute</option>
              <option value="5m">5 minutes</option>
              <option value="15m">15 minutes</option>
              <option value="30m">30 minutes</option>
              <option value="1h">1 hour</option>
              <option value="4h">4 hours</option>
              <option value="1d">1 day</option>
            </select>
          </div>
        )}
        
        <button 
          onClick={fetchData} 
          disabled={loading}
          className="fetch-button"
        >
          {loading ? 'Loading...' : 'Fetch Data'}
        </button>
      </div>
      
      {error && (
        <div className="error-message">
          <h3>Error:</h3>
          <p>{error}</p>
        </div>
      )}
      
      {responseData && (
        <div className="response-container">
          <h3>Response:</h3>
          <pre>{JSON.stringify(responseData, null, 2)}</pre>
        </div>
      )}
    </div>
  );
};

export default MexcApiTest;

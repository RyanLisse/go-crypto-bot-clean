import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import TradingPairSelector from '../TradingPairSelector';

describe('TradingPairSelector', () => {
  const availablePairs = [
    'BTC/USD',
    'ETH/USD',
    'SOL/USD',
    'XRP/USD',
    'ADA/USD'
  ];
  
  const mockOnChange = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders with the selected trading pair highlighted', () => {
    render(
      <TradingPairSelector 
        availablePairs={availablePairs} 
        selectedPair="BTC/USD" 
        onChange={mockOnChange}
      />
    );
    
    expect(screen.getByText('Trading Pairs')).toBeInTheDocument();
    
    // Check that all pairs are rendered
    availablePairs.forEach(pair => {
      expect(screen.getByText(pair)).toBeInTheDocument();
    });
    
    // The selected pair should have a different styling
    const selectedElement = screen.getByText('BTC/USD').closest('button');
    expect(selectedElement).toHaveClass('bg-brutal-active-bg');
  });

  it('calls onChange when a different pair is selected', () => {
    render(
      <TradingPairSelector 
        availablePairs={availablePairs} 
        selectedPair="BTC/USD" 
        onChange={mockOnChange}
      />
    );
    
    // Click on a different pair
    fireEvent.click(screen.getByText('ETH/USD'));
    
    // Verify the onChange handler was called with the correct value
    expect(mockOnChange).toHaveBeenCalledWith('ETH/USD');
  });

  it('does not call onChange when the currently selected pair is clicked', () => {
    render(
      <TradingPairSelector 
        availablePairs={availablePairs} 
        selectedPair="BTC/USD" 
        onChange={mockOnChange}
      />
    );
    
    // Click on the already selected pair
    fireEvent.click(screen.getByText('BTC/USD'));
    
    // Verify onChange was not called
    expect(mockOnChange).not.toHaveBeenCalled();
  });

  it('applies search filtering when search input is used', () => {
    render(
      <TradingPairSelector 
        availablePairs={availablePairs} 
        selectedPair="BTC/USD" 
        onChange={mockOnChange}
      />
    );
    
    // Get the search input and type in it
    const searchInput = screen.getByPlaceholderText('Search pairs...');
    fireEvent.change(searchInput, { target: { value: 'eth' } });
    
    // Only ETH/USD should be visible now
    expect(screen.getByText('ETH/USD')).toBeInTheDocument();
    expect(screen.queryByText('BTC/USD')).not.toBeInTheDocument();
    expect(screen.queryByText('SOL/USD')).not.toBeInTheDocument();
    
    // Clear the search
    fireEvent.change(searchInput, { target: { value: '' } });
    
    // All pairs should be visible again
    availablePairs.forEach(pair => {
      expect(screen.getByText(pair)).toBeInTheDocument();
    });
  });

  it('handles empty availablePairs array gracefully', () => {
    render(
      <TradingPairSelector 
        availablePairs={[]} 
        selectedPair="" 
        onChange={mockOnChange}
      />
    );
    
    expect(screen.getByText('Trading Pairs')).toBeInTheDocument();
    expect(screen.getByText('No trading pairs available')).toBeInTheDocument();
  });

  it('displays "loading" state when specified', () => {
    render(
      <TradingPairSelector 
        availablePairs={[]} 
        selectedPair="" 
        onChange={mockOnChange}
        isLoading={true}
      />
    );
    
    expect(screen.getByText('Loading trading pairs...')).toBeInTheDocument();
  });
}); 
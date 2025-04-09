import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Portfolio from './index';

describe('Portfolio Component', () => {
  it('renders the portfolio title', () => {
    render(<Portfolio />);
    
    expect(screen.getByText('Portfolio')).toBeInTheDocument();
  });
  
  it('displays portfolio summary section', () => {
    render(<Portfolio />);
    
    expect(screen.getByText('Portfolio Summary')).toBeInTheDocument();
  });
});

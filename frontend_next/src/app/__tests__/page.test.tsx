// Import custom JSDOM environment
import '../../test/bun-custom-jsdom';

// Test environment verification - make sure this is at the top
console.log('Test environment check:', {
  hasDocument: typeof document !== 'undefined',
  hasWindow: typeof window !== 'undefined',
  documentBody: typeof document !== 'undefined' ? document.body : 'undefined'
});

import React from 'react';
import { render, screen } from '@testing-library/react';
import Home from '../page';

// Log test-library methods
console.log('Testing library methods:', {
  render: typeof render,
  screen: typeof screen
});

describe('Home Page', () => {
  it('renders the main heading', () => {
    console.log('Before render - document exists:', typeof document !== 'undefined');
    render(<Home />);
    console.log('After render - document exists:', typeof document !== 'undefined');
    expect(screen.getByRole('heading', { name: /go crypto bot/i })).toBeInTheDocument();
  });
}); 
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import App from './App';

// Mock the components used in App.tsx
vi.mock('@/components/Layout', () => ({
  default: () => <div data-testid="layout">Layout Component</div>,
}));

vi.mock('@/pages/Dashboard', () => ({
  default: () => <div>Dashboard Component</div>,
}));

vi.mock('@/pages/Portfolio', () => ({
  default: () => <div>Portfolio Component</div>,
}));

vi.mock('@/pages/Trading', () => ({
  default: () => <div>Trading Component</div>,
}));

vi.mock('@/pages/Backtesting', () => ({
  default: () => <div>Backtesting Component</div>,
}));

vi.mock('@/pages/Settings', () => ({
  default: () => <div>Settings Component</div>,
}));

vi.mock('@/pages/NotFound', () => ({
  default: () => <div>NotFound Component</div>,
}));

vi.mock('@/components/ProtectedRoute', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

describe('App Component', () => {
  it('renders the layout component', () => {
    render(
      <BrowserRouter>
        <App />
      </BrowserRouter>
    );
    
    expect(screen.getByTestId('layout')).toBeInTheDocument();
  });
});

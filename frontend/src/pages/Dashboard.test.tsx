import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Dashboard from './Dashboard';

describe('Dashboard Page', () => {
  const queryClient = new QueryClient();

  it('renders portfolio summary section', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Dashboard />
      </QueryClientProvider>
    );
    expect(screen.getByText(/portfolio/i)).toBeInTheDocument();
  });

  it('renders active trades list', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Dashboard />
      </QueryClientProvider>
    );
    expect(screen.getByText(/active trades/i)).toBeInTheDocument();
  });

  it('renders market overview section', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Dashboard />
      </QueryClientProvider>
    );
    expect(screen.getByText(/market overview/i)).toBeInTheDocument();
  });

  it('renders performance charts', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Dashboard />
      </QueryClientProvider>
    );
    expect(screen.getByText(/performance/i)).toBeInTheDocument();
  });

  it('renders system status indicators', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Dashboard />
      </QueryClientProvider>
    );
    expect(screen.getByText(/system status/i)).toBeInTheDocument();
  });

  it('shows real-time update elements', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Dashboard />
      </QueryClientProvider>
    );
    expect(screen.getByText(/real-time/i)).toBeInTheDocument();
  });
});
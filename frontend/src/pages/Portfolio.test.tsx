import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Portfolio from './Portfolio';

describe('Portfolio Page', () => {
  const queryClient = new QueryClient();

  it('renders holdings list', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Portfolio />
      </QueryClientProvider>
    );
    expect(screen.getByText(/holdings/i)).toBeInTheDocument();
  });

  it('renders historical performance charts', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Portfolio />
      </QueryClientProvider>
    );
    expect(screen.getByText(/performance/i)).toBeInTheDocument();
  });

  it('renders position details', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Portfolio />
      </QueryClientProvider>
    );
    expect(screen.getByText(/position/i)).toBeInTheDocument();
  });

  it('renders export options', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <Portfolio />
      </QueryClientProvider>
    );
    expect(screen.getByText(/export/i)).toBeInTheDocument();
  });
});
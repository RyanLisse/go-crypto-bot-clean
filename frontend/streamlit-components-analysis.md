# Missing Frontend Components Analysis for Go Project

This document identifies the key components from the Python Streamlit application that are currently missing in the Go project's React frontend implementation. The focus is on ensuring feature parity between the two projects, particularly for the recently implemented date-based filtering functionality.

## Table of Contents

1. [Overview](#overview)
2. [Missing Components](#missing-components)
3. [Implementation Priorities](#implementation-priorities)
4. [API Integration Requirements](#api-integration-requirements)

## Overview

After analyzing the existing React frontend implementation in the `new_frontend` directory, we've identified several key components from the Python Streamlit application that need to be implemented to achieve feature parity. The current React frontend has a good foundation with dashboard, portfolio, and basic new coins functionality, but lacks several advanced features present in the Streamlit application.

## Missing Components

### 1. Date-Based New Coin Filtering
- **Current Status**: 
  - The existing NewCoins page uses mock data and lacks date-based filtering
  - The backend API integration for date filtering is not implemented
  - No date picker components for selecting date ranges

- **Required Features**:
  - Date range selection with date pickers (start date and end date)
  - API integration with the recently implemented date-based filtering endpoints
  - Session state management to remember selected date ranges
  - Navigation buttons to easily move to next/previous day

- **Implementation Notes**:
  - Need to update the API client to support date range parameters
  - Should implement proper error handling for API requests
  - Must include loading states during data fetching

### 2. Advanced New Coin Management
- **Current Status**:
  - Basic display of new coins exists but lacks management features
  - No archive/restore functionality for coins
  - No symbol filtering capability

- **Required Features**:
  - Symbol text filtering for finding specific coins
  - Option to include/exclude archived coins
  - Archive/restore functionality for coin management
  - Proper display of coin listing dates and times

- **Implementation Notes**:
  - Need to implement API endpoints for archiving/restoring coins
  - Should include confirmation dialogs for important actions
  - Must maintain consistent state after archive/restore operations

### 3. Transaction History Improvements
- **Current Status**:
  - Basic transaction display exists but lacks filtering options
  - No date-based filtering for transactions

- **Required Features**:
  - Date range filtering for transactions
  - Symbol filtering for transactions
  - Sorting options for different columns
  - Detailed transaction information display

- **Implementation Notes**:
  - Need to update API client to support transaction filtering
  - Should implement pagination for large result sets
  - Must ensure consistent timestamp formatting

### 4. Log Events Viewer Enhancements
- **Current Status**:
  - No dedicated log events viewer in the current implementation
  - Missing filtering capabilities for logs

- **Required Features**:
  - Date range filtering for log events
  - Event type filtering via dropdown
  - Text search in log details
  - Categorized display of different log types
  - Error categorization and visualization

- **Implementation Notes**:
  - Need to implement API endpoints for log retrieval with filters
  - Should include export functionality for logs
  - Must handle different log event types appropriately

### 5. Performance Dashboard Improvements
- **Current Status**:
  - Basic performance metrics exist but lack detailed analysis
  - Missing historical session selection

- **Required Features**:
  - Session selection for historical performance data
  - Detailed trade profit analysis
  - Error pattern analysis and visualization
  - Time-series performance metrics

- **Implementation Notes**:
  - Need to implement API endpoints for historical performance data
  - Should support interactive charts for time series data
  - Must handle performance data efficiently

## Implementation Priorities

Based on the recent backend implementation of date-based filtering for new coins, we recommend the following implementation priorities:

1. **High Priority**:
   - Date-based new coin filtering (to match the recently implemented backend functionality)
   - Symbol filtering for new coins
   - Archive/restore functionality for new coins

2. **Medium Priority**:
   - Transaction history improvements
   - Log events viewer enhancements

3. **Lower Priority**:
   - Performance dashboard improvements
   - Additional visualization components

## API Integration Requirements

To support the missing components, particularly the date-based filtering functionality, the following API integrations need to be implemented:

### 1. New Coin API Endpoints
- **Current Implementation**:
  ```typescript
  getNewCoins: async (): Promise<any[]> => {
    try {
      const response = await fetch(`${API_BASE_URL}/newcoins`);
      // ...
    }
  }
  ```

- **Required Updates**:
  ```typescript
  // Get new coins with date filtering
  getNewCoinsByDate: async (date: string): Promise<any[]> => {
    try {
      const response = await fetch(`${API_BASE_URL}/newcoins/by-date?date=${date}`);
      // ...
    }
  },
  
  // Get new coins with date range filtering
  getNewCoinsByDateRange: async (startDate: string, endDate: string): Promise<any[]> => {
    try {
      const response = await fetch(`${API_BASE_URL}/newcoins/by-date-range?startDate=${startDate}&endDate=${endDate}`);
      // ...
    }
  },
  
  // Archive/restore coin
  updateCoinStatus: async (symbol: string, archived: boolean): Promise<any> => {
    try {
      const response = await fetch(`${API_BASE_URL}/newcoins/${symbol}/status`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ archived })
      });
      // ...
    }
  }
  ```

### 2. React Query Integration
- Implement React Query hooks for the new API endpoints:
  ```typescript
  // Query keys
  export const newCoinKeys = {
    all: ['newCoins'] as const,
    byDate: (date: string) => [...newCoinKeys.all, 'byDate', date] as const,
    byDateRange: (startDate: string, endDate: string) => [...newCoinKeys.all, 'byDateRange', startDate, endDate] as const,
  };
  
  // Get new coins by date
  export const useNewCoinsByDateQuery = (date: string) => {
    return useQuery({
      queryKey: newCoinKeys.byDate(date),
      queryFn: () => api.getNewCoinsByDate(date),
      staleTime: 30000,
    });
  };
  
  // Get new coins by date range
  export const useNewCoinsByDateRangeQuery = (startDate: string, endDate: string) => {
    return useQuery({
      queryKey: newCoinKeys.byDateRange(startDate, endDate),
      queryFn: () => api.getNewCoinsByDateRange(startDate, endDate),
      staleTime: 30000,
    });
  };
  ```

By implementing these API integrations and the missing UI components, the Go project's frontend will achieve feature parity with the Python Streamlit application, particularly for the recently implemented date-based filtering functionality.

import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

// Define the base API
export const api = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1',
    prepareHeaders: (headers) => {
      // Get token from localStorage
      const token = localStorage.getItem('token');
      
      // If we have a token, add it to the headers
      if (token) {
        headers.set('authorization', `Bearer ${token}`);
      }
      
      return headers;
    },
  }),
  tagTypes: ['Portfolio', 'Trade', 'Backtest', 'Settings'],
  endpoints: (builder) => ({
    // Define endpoints here
    getStatus: builder.query<{ status: string }, void>({
      query: () => '/status',
    }),
  }),
});

// Export hooks for usage in components
export const { useGetStatusQuery } = api;

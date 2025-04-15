import '@testing-library/jest-dom';
import React from 'react';
import { vi } from 'vitest';

// Mock next/image for Vitest
vi.mock('next/image', () => {
  return {
    __esModule: true,
    default: (props: React.ImgHTMLAttributes<HTMLImageElement>) => {
      // eslint-disable-next-line @next/next/no-img-element
      return React.createElement('img', props);
    },
  };
}); 
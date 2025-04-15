import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Enable React strict mode for improved error detection
  reactStrictMode: true,
  
  // Environment variables that will be available in the browser
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || '',
    NEXT_PUBLIC_WS_URL: process.env.NEXT_PUBLIC_WS_URL || '',
  },
  
  // Configure images to allow external domains if needed
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
    ],
  },
  
  // Configure redirects for any needed path forwarding
  async redirects() {
    return [
      {
        source: '/home',
        destination: '/dashboard',
        permanent: true,
      },
    ];
  },
  
  // Add trailing slashes to all URLs
  trailingSlash: false,
  
  // Configure various compiler options
  compiler: {
    // Remove console.log in production
    removeConsole: process.env.NODE_ENV === 'production' 
      ? { exclude: ['error', 'warn'] } 
      : false,
  },
  
  // Transpile specific modules that require it
  transpilePackages: [],
  
  // Override the default webpack config if needed
  webpack: (config) => {
    // Add custom webpack configurations here if needed
    return config;
  },
};

export default nextConfig;

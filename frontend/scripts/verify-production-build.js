#!/usr/bin/env node

/**
 * Production Build Verification Script
 * 
 * This script verifies that the frontend is ready for production deployment on Netlify.
 * It checks for common issues and provides recommendations for fixing them.
 */

import fs from 'fs';
import path from 'path';
import { execSync } from 'child_process';
import { fileURLToPath } from 'url';

// ANSI color codes for better readability
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m',
};

// Helper functions
const log = {
  info: (msg) => console.log(`${colors.blue}ℹ️ ${msg}${colors.reset}`),
  success: (msg) => console.log(`${colors.green}✅ ${msg}${colors.reset}`),
  warning: (msg) => console.log(`${colors.yellow}⚠️ ${msg}${colors.reset}`),
  error: (msg) => console.log(`${colors.red}❌ ${msg}${colors.reset}`),
  section: (msg) => console.log(`\n${colors.bright}${colors.cyan}== ${msg} ==${colors.reset}\n`)
};

// Get the project root directory
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const projectRoot = path.resolve(__dirname, '..');

// Check if a file exists
function fileExists(filePath) {
  return fs.existsSync(path.join(projectRoot, filePath));
}

// Read a file and return its contents
function readFile(filePath) {
  try {
    return fs.readFileSync(path.join(projectRoot, filePath), 'utf8');
  } catch (error) {
    return null;
  }
}

// Check if a file contains a specific string
function fileContains(filePath, searchString) {
  const content = readFile(filePath);
  return content && content.includes(searchString);
}

// Main verification function
async function verifyProductionBuild() {
  let issues = 0;
  let warnings = 0;
  
  log.section('Verifying Frontend Production Build for Netlify');
  
  // Check for required files
  log.section('Checking Required Files');
  
  const requiredFiles = [
    { path: 'netlify.toml', description: 'Netlify configuration' },
    { path: 'package.json', description: 'Package configuration' },
    { path: 'vite.config.ts', description: 'Vite configuration' },
    { path: '.env.production', description: 'Production environment variables' },
  ];
  
  for (const file of requiredFiles) {
    if (fileExists(file.path)) {
      log.success(`Found ${file.path} (${file.description})`);
    } else {
      log.error(`Missing ${file.path} (${file.description})`);
      issues++;
    }
  }
  
  // Check Netlify configuration
  log.section('Checking Netlify Configuration');
  
  const netlifyToml = readFile('netlify.toml');
  if (netlifyToml) {
    if (netlifyToml.includes('publish = "dist"')) {
      log.success('Netlify publish directory correctly set to "dist"');
    } else {
      log.error('Netlify publish directory not correctly configured');
      issues++;
    }
    
    if (netlifyToml.includes('[[redirects]]') && netlifyToml.includes('from = "/*"') && netlifyToml.includes('to = "/index.html"')) {
      log.success('SPA redirect rule properly configured');
    } else {
      log.error('SPA redirect rule missing or incorrectly configured');
      issues++;
    }
    
    if (netlifyToml.includes('[[headers]]') && netlifyToml.includes('X-Frame-Options')) {
      log.success('Security headers configured');
    } else {
      log.warning('Security headers not fully configured');
      warnings++;
    }
  }
  
  // Check environment variables
  log.section('Checking Environment Variables');
  
  const envProduction = readFile('.env.production');
  if (envProduction) {
    if (envProduction.includes('${API_URL}') && envProduction.includes('${WS_URL}')) {
      log.success('Production environment variables using Netlify variable substitution');
    } else if (envProduction.includes('api.your-production-domain.com')) {
      log.error('Production API URLs still using placeholder values');
      issues++;
    }
  }
  
  // Check build configuration
  log.section('Checking Build Configuration');
  
  const viteConfig = readFile('vite.config.ts');
  if (viteConfig) {
    if (viteConfig.includes('sourcemap: process.env.NODE_ENV !== \'production\'')) {
      log.success('Source maps disabled in production for better performance and security');
    } else if (viteConfig.includes('sourcemap: true')) {
      log.warning('Source maps enabled in production builds - consider disabling for security');
      warnings++;
    }
    
    if (viteConfig.includes('manualChunks')) {
      log.success('Code splitting configured for better performance');
    } else {
      log.warning('Code splitting not configured - consider adding for better performance');
      warnings++;
    }
  }
  
  // Check for production optimizations
  log.section('Checking Production Optimizations');
  
  try {
    const packageJson = JSON.parse(readFile('package.json'));
    const hasDependencyCheck = packageJson.scripts && 
      (packageJson.scripts['check:deps'] || packageJson.scripts['audit'] || packageJson.scripts['audit:fix']);
    
    if (hasDependencyCheck) {
      log.success('Dependency audit scripts found');
    } else {
      log.warning('No dependency audit scripts found - consider adding npm audit');
      warnings++;
    }
  } catch (error) {
    log.error('Error parsing package.json');
    issues++;
  }
  
  // Check for documentation
  log.section('Checking Deployment Documentation');
  
  if (fileExists('docs/netlify-deployment-guide.md')) {
    log.success('Netlify deployment guide found');
  } else {
    log.warning('No Netlify deployment guide found - consider adding documentation');
    warnings++;
  }
  
  // Summary
  log.section('Verification Summary');
  
  if (issues === 0 && warnings === 0) {
    log.success('All checks passed! The frontend is ready for production deployment on Netlify.');
  } else {
    if (issues > 0) {
      log.error(`Found ${issues} issue${issues !== 1 ? 's' : ''} that need to be fixed before deployment.`);
    }
    if (warnings > 0) {
      log.warning(`Found ${warnings} warning${warnings !== 1 ? 's' : ''} that should be addressed for optimal deployment.`);
    }
  }
  
  // Next steps
  log.section('Next Steps');
  
  console.log(`
${colors.bright}To deploy to Netlify:${colors.reset}

1. Ensure all issues above are fixed
2. Set up environment variables in Netlify:
   - API_URL: Your production API endpoint
   - WS_URL: Your production WebSocket endpoint
3. Connect your repository to Netlify
4. Configure build settings:
   - Base directory: new_frontend
   - Build command: npm run build
   - Publish directory: dist

${colors.bright}For detailed instructions:${colors.reset}
See the deployment guide at docs/netlify-deployment-guide.md
`);
}

// Run the verification
verifyProductionBuild().catch(error => {
  console.error('Error during verification:', error);
  process.exit(1);
});

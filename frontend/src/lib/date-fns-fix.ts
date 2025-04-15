/**
 * This file provides a workaround for date-fns import issues with Vite
 * It re-exports the commonly used date-fns functions to avoid the .mjs import errors
 */

import { format, formatDistance, formatRelative, subDays, addDays, isValid, parse, parseISO } from 'date-fns';
import { enUS } from 'date-fns/locale';

// Re-export the functions
export {
  format,
  formatDistance,
  formatRelative,
  subDays,
  addDays,
  isValid,
  parse,
  parseISO,
  enUS
};

// Add any other date-fns functions you need here

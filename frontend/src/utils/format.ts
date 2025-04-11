/**
 * Formats a number with appropriate decimal places and thousands separators
 * @param value The number to format
 * @param decimals Optional number of decimal places (default: 2)
 * @returns Formatted string representation of the number
 */
export const formatNumber = (value: number, decimals: number = 2): string => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(value);
}; 
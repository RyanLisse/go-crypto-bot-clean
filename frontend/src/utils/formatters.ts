/**
 * Format a number with appropriate decimal places and thousands separators
 * @param num - The number to format
 * @param decimals - Optional number of decimal places (default: 2)
 * @returns Formatted number string
 */
export const formatNumber = (num: number, decimals: number = 2): string => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(num);
}; 
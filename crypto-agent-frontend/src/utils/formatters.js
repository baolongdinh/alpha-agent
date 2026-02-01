/**
 * Format large numbers with K, M, B suffixes
 * @param {number} num - Number to format
 * @returns {string} Formatted number
 */
export function formatNumber(num) {
  if (!num || num === 0) return '0'
  
  const absNum = Math.abs(num)
  
  if (absNum >= 1e9) {
    return `${(num / 1e9).toFixed(2)}B`
  }
  if (absNum >= 1e6) {
    return `${(num / 1e6).toFixed(2)}M`
  }
  if (absNum >= 1e3) {
    return `${(num / 1e3).toFixed(2)}K`
  }
  
  return num.toLocaleString('en-US', { maximumFractionDigits: 0 })
}

/**
 * Format currency with full precision
 * @param {number} num - Number to format
 * @param {number} decimals - Decimal places (default: 2)
 * @returns {string} Formatted currency
 */
export function formatCurrency(num) {
  if (!num || num === 0) return '$0.00'
  
  // Use high precision for small numbers
  let decimals = 2
  if (num < 0.01) decimals = 8
  else if (num < 1) decimals = 6
  
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(num)
}

/**
 * Format crypto price with dynamic precision
 * @param {number} price - Price to format
 * @returns {string} Formatted price
 */
export function formatPrice(price) {
  if (!price || price === 0) return '0.00'
  
  // For small prices (like memecoins), show more decimals
  if (price < 0.01) return price.toFixed(8)
  if (price < 1) return price.toFixed(6)
  if (price < 10) return price.toFixed(4)
  
  // Standard formatting for larger prices
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(price)
}

/**
 * Format percentage with color
 * @param {number} value - Percentage value
 * @returns {Object} { text, color, icon }
 */
export function formatPercentage(value) {
  if (value === null || value === undefined) {
    return { text: 'N/A', color: 'text-gray-400', icon: '' }
  }
  
  const formatted = `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`
  const color = value >= 0 ? 'text-green-400' : 'text-red-400'
  const icon = value >= 0 ? '↑' : '↓'
  
  return { text: formatted, color, icon }
}

/**
 * Format timestamp to readable date
 * @param {string|Date} timestamp - Timestamp to format
 * @returns {string} Formatted date
 */
export function formatDate(timestamp) {
  return new Date(timestamp).toLocaleString('vi-VN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/**
 * Format time ago (e.g., "2 minutes ago")
 * @param {string|Date} timestamp - Timestamp
 * @returns {string} Relative time
 */
export function timeAgo(timestamp) {
  const seconds = Math.floor((new Date() - new Date(timestamp)) / 1000)
  
  const intervals = {
    năm: 31536000,
    tháng: 2592000,
    tuần: 604800,
    ngày: 86400,
    giờ: 3600,
    phút: 60,
    giây: 1,
  }
  
  for (const [name, value] of Object.entries(intervals)) {
    const interval = Math.floor(seconds / value)
    if (interval >= 1) {
      return `${interval} ${name} trước`
    }
  }
  
  return 'vừa xong'
}

/**
 * Truncate text with ellipsis
 * @param {string} text - Text to truncate
 * @param {number} maxLength - Maximum length
 * @returns {string} Truncated text
 */
export function truncate(text, maxLength = 50) {
  if (!text || text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

/**
 * Format rank with ordinal suffix
 * @param {number} rank - Rank number
 * @returns {string} Formatted rank (e.g., "1st", "2nd")
 */
export function formatRank(rank) {
  if (!rank) return 'N/A'
  
  const suffix = ['th', 'st', 'nd', 'rd']
  const v = rank % 100
  return rank + (suffix[(v - 20) % 10] || suffix[v] || suffix[0])
}

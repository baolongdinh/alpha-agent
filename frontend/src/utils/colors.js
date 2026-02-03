/**
 * Get color class based on trust score
 * @param {number} score - Trust score (0-100)
 * @returns {Object} { bg, text, border, ring }
 */
export function getScoreColor(score) {
  if (score >= 80) {
    return {
      bg: 'bg-green-500/20',
      text: 'text-green-400',
      border: 'border-green-500/30',
      ring: 'ring-green-500/50',
      gradient: 'from-green-500 to-emerald-500',
    }
  }
  
  if (score >= 60) {
    return {
      bg: 'bg-blue-500/20',
      text: 'text-blue-400',
      border: 'border-blue-500/30',
      ring: 'ring-blue-500/50',
      gradient: 'from-blue-500 to-cyan-500',
    }
  }
  
  if (score >= 40) {
    return {
      bg: 'bg-yellow-500/20',
      text: 'text-yellow-400',
      border: 'border-yellow-500/30',
      ring: 'ring-yellow-500/50',
      gradient: 'from-yellow-500 to-orange-500',
    }
  }
  
  return {
    bg: 'bg-red-500/20',
    text: 'text-red-400',
    border: 'border-red-500/30',
    ring: 'ring-red-500/50',
    gradient: 'from-red-500 to-pink-500',
  }
}

/**
 * Get color class based on percentage change
 * @param {number} value - Percentage change
 * @returns {string} Tailwind color class
 */
export function getChangeColor(value) {
  if (value > 0) return 'text-green-400'
  if (value < 0) return 'text-red-400'
  return 'text-gray-400'
}

/**
 * Get background gradient based on score
 * @param {number} score - Trust score (0-100)
 * @returns {string} Tailwind gradient classes
 */
export function getScoreGradient(score) {
  if (score >= 80) {
    return 'bg-gradient-to-r from-green-500 to-emerald-400'
  }
  if (score >= 60) {
    return 'bg-gradient-to-r from-blue-500 to-cyan-400'
  }
  if (score >= 40) {
    return 'bg-gradient-to-r from-yellow-500 to-orange-400'
  }
  return 'bg-gradient-to-r from-red-500 to-pink-400'
}

/**
 * Get grade badge color
 * @param {string} grade - Grade (S, A, B, C, D, F)
 * @returns {Object} Color classes
 */
export function getGradeColor(grade) {
  const colors = {
    S: {
      bg: 'bg-purple-500/20',
      text: 'text-purple-400',
      border: 'border-purple-500/30',
    },
    A: {
      bg: 'bg-green-500/20',
      text: 'text-green-400',
      border: 'border-green-500/30',
    },
    B: {
      bg: 'bg-blue-500/20',
      text: 'text-blue-400',
      border: 'border-blue-500/30',
    },
    C: {
      bg: 'bg-yellow-500/20',
      text: 'text-yellow-400',
      border: 'border-yellow-500/30',
    },
    D: {
      bg: 'bg-orange-500/20',
      text: 'text-orange-400',
      border: 'border-orange-500/30',
    },
    F: {
      bg: 'bg-red-500/20',
      text: 'text-red-400',
      border: 'border-red-500/30',
    },
  }
  
  return colors[grade] || colors.F
}

/**
 * Get confidence level color
 * @param {number} confidence - Confidence score (0-100)
 * @returns {string} Color class
 */
export function getConfidenceColor(confidence) {
  if (confidence >= 80) return 'text-green-400'
  if (confidence >= 60) return 'text-yellow-400'
  return 'text-red-400'
}

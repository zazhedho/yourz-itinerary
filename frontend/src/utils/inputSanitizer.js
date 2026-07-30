/**
 * Live input sanitizer for text fields:
 * 1. Removes leading spaces (cannot start a field with space)
 * 2. Collapses multiple consecutive spaces into a single space ("   " -> " ")
 */
export const sanitizeLiveInput = (value) => {
  if (typeof value !== 'string') return value
  return value.replace(/^\s+/, '').replace(/\s{2,}/g, ' ')
}

export const sanitizeNoWhitespaceInput = (value) => {
  if (typeof value !== 'string') return value
  return value.replace(/\s/g, '')
}

/**
 * On Blur input sanitizer:
 * Removes trailing spaces when the input loses focus ("Bali  " -> "Bali")
 */
export const sanitizeBlurInput = (value) => {
  if (typeof value !== 'string') return value
  return value.trim()
}

import { describe, expect, it } from 'vitest'
import { sanitizeBlurInput, sanitizeLiveInput, sanitizeNoWhitespaceInput } from './inputSanitizer'

describe('inputSanitizer', () => {
  it('removes leading spaces and collapses consecutive spaces on live input', () => {
    expect(sanitizeLiveInput('   Hello')).toBe('Hello')
    expect(sanitizeLiveInput('Bali   Indah')).toBe('Bali Indah')
    expect(sanitizeLiveInput('  Double   Space  ')).toBe('Double Space ')
  })

  it('trims leading and trailing spaces on blur', () => {
    expect(sanitizeBlurInput('  Bali Indah  ')).toBe('Bali Indah')
  })

  it('removes all whitespace from fields that do not allow it', () => {
    expect(sanitizeNoWhitespaceInput(' zaqi @example.com ')).toBe('zaqi@example.com')
  })
})

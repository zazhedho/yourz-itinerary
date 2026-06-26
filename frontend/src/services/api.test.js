import { describe, expect, it } from 'vitest'

import { getErrorMessage, getResponseData } from './api'

describe('api helpers', () => {
  it('extracts response data from backend envelope', () => {
    expect(getResponseData({ data: { data: { id: 'trip-1' } } })).toEqual({ id: 'trip-1' })
  })

  it('joins validation errors from backend envelope', () => {
    const error = {
      response: {
        data: {
          error: [{ message: 'email is required' }, { message: 'password is required' }],
        },
      },
    }

    expect(getErrorMessage(error, 'fallback')).toBe('email is required, password is required')
  })
})

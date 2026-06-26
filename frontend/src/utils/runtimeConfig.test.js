import { afterEach, describe, expect, it } from 'vitest'

import { getRuntimeConfigValue } from './runtimeConfig'

describe('runtime config', () => {
  afterEach(() => {
    delete window.ENV_CONFIG
  })

  it('prefers window ENV_CONFIG over Vite env', () => {
    window.ENV_CONFIG = { GOOGLE_CLIENT_ID: 'runtime-client-id' }

    expect(getRuntimeConfigValue('GOOGLE_CLIENT_ID', 'VITE_GOOGLE_CLIENT_ID')).toBe('runtime-client-id')
  })

  it('falls back to Vite env when runtime config is empty', () => {
    expect(getRuntimeConfigValue('MISSING_RUNTIME_KEY', 'MODE')).toBe('test')
  })
})

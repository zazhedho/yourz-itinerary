import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import Login from './Login'

const mocks = vi.hoisted(() => ({
  googleClientId: 'google-client-id',
  googleLogin: vi.fn(),
}))

vi.mock('../../components/common/GoogleIdentityButton', () => ({
  default: ({ label, onCredential }) => <button onClick={() => onCredential('google-token')}>{label}</button>,
}))

vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => ({
    error: '',
    googleLogin: mocks.googleLogin,
    login: vi.fn(),
  }),
}))

vi.mock('../../hooks/useRegisterStatus', () => ({
  default: () => ({ enabled: true }),
}))

vi.mock('../../utils/runtimeConfig', () => ({
  getGoogleClientId: () => mocks.googleClientId,
}))

const renderLogin = () => render(<Login />, { wrapper: MemoryRouter })

describe('Login', () => {
  beforeEach(() => {
    mocks.googleClientId = 'google-client-id'
    mocks.googleLogin.mockReset()
    mocks.googleLogin.mockResolvedValue(true)
  })

  it('allows login with Google credential', async () => {
    const user = userEvent.setup()
    renderLogin()

    await user.click(screen.getByRole('button', { name: /lanjutkan dengan google/i }))

    expect(mocks.googleLogin).toHaveBeenCalledWith('google-token')
  })

  it('hides Google login when Google client id is not configured', () => {
    mocks.googleClientId = ''
    renderLogin()

    expect(screen.queryByRole('button', { name: /lanjutkan dengan google/i })).not.toBeInTheDocument()
  })
})

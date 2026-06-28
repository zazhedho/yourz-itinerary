import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import authService from '../../services/authService'
import Register from './Register'

const mocks = vi.hoisted(() => ({
  googleClientId: 'google-client-id',
  googleLogin: vi.fn(),
}))

vi.mock('../../components/common/GoogleIdentityButton', () => ({
  default: ({ label, onCredential }) => <button onClick={() => onCredential('google-token')}>{label}</button>,
}))

vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => ({
    googleLogin: mocks.googleLogin,
  }),
}))

vi.mock('../../utils/runtimeConfig', () => ({
  getGoogleClientId: () => mocks.googleClientId,
}))

vi.mock('../../services/authService', () => ({
  default: {
    getRegisterStatus: vi.fn(),
    sendRegisterOTP: vi.fn(),
    register: vi.fn(),
  },
}))

const renderRegister = () => render(<Register />, { wrapper: MemoryRouter })

describe('Register', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    sessionStorage.clear()
    mocks.googleClientId = 'google-client-id'
    mocks.googleLogin.mockResolvedValue(true)
    authService.getRegisterStatus.mockResolvedValue({
      data: { data: { enabled: true, otp_enabled: true, otp_cooldown: 60 } },
    })
    authService.sendRegisterOTP.mockResolvedValue({ data: { data: {} } })
    authService.register.mockResolvedValue({ data: { data: {} } })
  })

  it('shows OTP only after requesting it when register OTP is enabled', async () => {
    const user = userEvent.setup()
    renderRegister()

    await screen.findByLabelText(/nama/i)
    expect(screen.queryByLabelText(/kode otp/i)).not.toBeInTheDocument()

    await user.type(screen.getByLabelText(/nama/i), 'Zaqi')
    await user.type(screen.getByLabelText(/email/i), 'zaqi@example.com')
    await user.type(screen.getByLabelText(/nomor hp/i), '628123456789')
    await user.type(screen.getByLabelText(/^password$/i), 'Password123!')
    await user.type(screen.getByLabelText(/konfirmasi password/i), 'Password123!')
    await user.click(screen.getByRole('button', { name: /^daftar$/i }))

    await waitFor(() => expect(authService.sendRegisterOTP).toHaveBeenCalledWith({
      email: 'zaqi@example.com',
      phone: '628123456789',
    }))
    expect(authService.register).not.toHaveBeenCalled()
    expect(await screen.findByLabelText(/kode otp/i)).toBeInTheDocument()
  })

  it('keeps OTP input numeric and limited to six digits', async () => {
    const user = userEvent.setup()
    sessionStorage.setItem('register_otp_step', 'true')
    sessionStorage.setItem('register_form', JSON.stringify({
      name: 'Zaqi',
      email: 'zaqi@example.com',
      phone: '628123456789',
      password: 'Password123!',
      confirm_password: 'Password123!',
      otp_code: '',
    }))

    renderRegister()

    const otpInput = await screen.findByLabelText(/kode otp/i)
    await user.type(otpInput, '12ab345678')

    expect(otpInput).toHaveValue('123456')
  })

  it('blocks OTP request when password requirements are not met', async () => {
    const user = userEvent.setup()
    renderRegister()

    await screen.findByLabelText(/nama/i)
    await user.type(screen.getByLabelText(/nama/i), 'Zaqi')
    await user.type(screen.getByLabelText(/email/i), 'zaqi@example.com')
    await user.type(screen.getByLabelText(/nomor hp/i), '628123456789')
    await user.type(screen.getByLabelText(/^password$/i), 'password123')
    await user.type(screen.getByLabelText(/konfirmasi password/i), 'password123')
    await user.click(screen.getByRole('button', { name: /^daftar$/i }))

    expect(await screen.findByText(/password belum memenuhi/i)).toBeInTheDocument()
    expect(authService.sendRegisterOTP).not.toHaveBeenCalled()
  })

  it('allows registration with Google credential', async () => {
    const user = userEvent.setup()
    renderRegister()

    await screen.findByLabelText(/nama/i)
    await user.click(screen.getByRole('button', { name: /daftar dengan google/i }))

    expect(mocks.googleLogin).toHaveBeenCalledWith('google-token')
  })

  it('hides Google registration when Google client id is not configured', async () => {
    mocks.googleClientId = ''
    renderRegister()

    await screen.findByLabelText(/nama/i)
    expect(screen.queryByRole('button', { name: /daftar dengan google/i })).not.toBeInTheDocument()
  })

  it('moves to OTP step when backend requires otp_code despite stale status', async () => {
    const user = userEvent.setup()
    authService.getRegisterStatus.mockResolvedValue({
      data: { data: { enabled: true, otp_enabled: false, otp_cooldown: 60 } },
    })
    authService.register.mockRejectedValue({
      response: {
        data: {
          error: { message: 'otp_code is required' },
        },
      },
    })

    renderRegister()

    await screen.findByLabelText(/nama/i)
    await user.type(screen.getByLabelText(/nama/i), 'Zaqi')
    await user.type(screen.getByLabelText(/email/i), 'zaqi@example.com')
    await user.type(screen.getByLabelText(/nomor hp/i), '628123456789')
    await user.type(screen.getByLabelText(/^password$/i), 'Password123!')
    await user.type(screen.getByLabelText(/konfirmasi password/i), 'Password123!')
    await user.click(screen.getByRole('button', { name: /^daftar$/i }))

    await waitFor(() => expect(authService.sendRegisterOTP).toHaveBeenCalledWith({
      email: 'zaqi@example.com',
      phone: '628123456789',
    }))
    expect(await screen.findByLabelText(/kode otp/i)).toBeInTheDocument()
  })
})

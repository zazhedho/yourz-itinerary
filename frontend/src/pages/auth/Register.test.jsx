import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import authService from '../../services/authService'
import Register from './Register'

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
    await user.click(screen.getByRole('button', { name: /daftar/i }))

    await waitFor(() => expect(authService.sendRegisterOTP).toHaveBeenCalledWith({
      email: 'zaqi@example.com',
      phone: '628123456789',
    }))
    expect(authService.register).not.toHaveBeenCalled()
    expect(await screen.findByLabelText(/kode otp/i)).toBeInTheDocument()
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
    await user.click(screen.getByRole('button', { name: /daftar/i }))

    expect(await screen.findByText(/password belum memenuhi/i)).toBeInTheDocument()
    expect(authService.sendRegisterOTP).not.toHaveBeenCalled()
  })
})

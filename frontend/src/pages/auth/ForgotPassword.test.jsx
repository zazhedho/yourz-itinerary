import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import authService from '../../services/authService'
import ForgotPassword from './ForgotPassword'

vi.mock('../../services/authService', () => ({
  default: {
    forgotPassword: vi.fn(),
  },
}))

describe('ForgotPassword', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    sessionStorage.clear()
    authService.forgotPassword.mockResolvedValue({ data: { data: { cooldown: 60 } } })
  })

  it('shows and persists the resend countdown after requesting a reset', async () => {
    const user = userEvent.setup()
    const { unmount } = render(<ForgotPassword />, { wrapper: MemoryRouter })

    await user.type(screen.getByLabelText(/email/i), 'user@example.com')
    await user.click(screen.getByRole('button', { name: /kirim instruksi/i }))

    expect(await screen.findByRole('button', { name: /kirim ulang dalam 60s/i })).toBeDisabled()

    unmount()
    render(<ForgotPassword />, { wrapper: MemoryRouter })
    expect(screen.getByRole('button', { name: /kirim ulang dalam 60s/i })).toBeDisabled()
  })
})

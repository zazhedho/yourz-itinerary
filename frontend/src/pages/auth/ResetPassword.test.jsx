import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import authService from '../../services/authService'
import ResetPassword from './ResetPassword'

vi.mock('../../services/authService', () => ({
  default: {
    resetPassword: vi.fn(),
  },
}))

const renderResetPassword = () =>
  render(
    <MemoryRouter initialEntries={['/reset-password?token=reset-token']}>
      <Routes>
        <Route path="/reset-password" element={<ResetPassword />} />
        <Route path="/login" element={<div>Login page</div>} />
      </Routes>
    </MemoryRouter>,
  )

describe('ResetPassword', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    authService.resetPassword.mockResolvedValue({})
  })

  it('submits matching valid password and returns to login', async () => {
    const user = userEvent.setup()
    renderResetPassword()

    await user.type(screen.getByLabelText(/^password baru$/i), 'Password123!')
    await user.type(screen.getByLabelText(/konfirmasi password/i), 'Password123!')
    await user.click(screen.getByRole('button', { name: /simpan password/i }))

    await waitFor(() =>
      expect(authService.resetPassword).toHaveBeenCalledWith({
        token: 'reset-token',
        new_password: 'Password123!',
      }),
    )
    expect(await screen.findByText('Login page')).toBeInTheDocument()
  })
})

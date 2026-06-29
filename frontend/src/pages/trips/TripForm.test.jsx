import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { AuthContext } from '../../contexts/auth-context'
import tripService from '../../services/tripService'
import TripForm from './TripForm'

vi.mock('../../services/tripService', () => ({
  default: {
    create: vi.fn(),
    update: vi.fn(),
    getById: vi.fn(),
  },
}))

const authValue = {
  user: { id: 'user-1' },
  loading: false,
  isAuthenticated: true,
}

const renderTripForm = () => render(
  <AuthContext.Provider value={authValue}>
    <TripForm />
  </AuthContext.Provider>,
  { wrapper: MemoryRouter },
)

describe('TripForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    tripService.create.mockResolvedValue({ data: { data: { id: 'trip-1' } } })
  })

  it('prevents end date from being earlier than start date', async () => {
    const user = userEvent.setup()
    renderTripForm()

    const startDate = screen.getByLabelText(/mulai/i)
    const endDate = screen.getByLabelText(/selesai/i)

    await user.type(screen.getByLabelText(/judul trip/i), 'Bali Trip')
    await user.type(startDate, '2026-07-10')

    expect(endDate).toHaveAttribute('min', '2026-07-10')

    await user.type(endDate, '2026-07-09')
    await user.click(screen.getByRole('button', { name: /simpan trip/i }))

    expect(tripService.create).not.toHaveBeenCalled()
  })

  it('moves end date forward when start date passes existing end date', async () => {
    const user = userEvent.setup()
    renderTripForm()

    const startDate = screen.getByLabelText(/mulai/i)
    const endDate = screen.getByLabelText(/selesai/i)

    await user.type(endDate, '2026-07-09')
    await user.type(startDate, '2026-07-10')

    await waitFor(() => expect(endDate).toHaveValue('2026-07-10'))
  })
})

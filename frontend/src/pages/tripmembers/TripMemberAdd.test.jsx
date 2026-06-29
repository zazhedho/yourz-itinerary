import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import useTripAccess from '../../hooks/useTripAccess'
import TripMemberAdd from './TripMemberAdd'

vi.mock('../../hooks/useTripAccess')

const renderPage = () => render(
  <MemoryRouter initialEntries={['/trips/trip-1/members/add']}>
    <Routes>
      <Route path="/trips/:tripId/members/add" element={<TripMemberAdd />} />
    </Routes>
  </MemoryRouter>,
)

describe('TripMemberAdd', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('blocks member form when current user cannot manage members', () => {
    useTripAccess.mockReturnValue({
      allowed: false,
      error: '',
      loading: false,
    })

    renderPage()

    expect(screen.getByText('Akses ditolak')).toBeInTheDocument()
    expect(screen.queryByLabelText(/email member/i)).not.toBeInTheDocument()
  })

  it('shows member form when current user can manage members', () => {
    useTripAccess.mockReturnValue({
      allowed: true,
      error: '',
      loading: false,
    })

    renderPage()

    expect(screen.getByLabelText(/email member/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /tambah member/i })).toBeInTheDocument()
  })

  it('explains viewer and editor roles before submit', () => {
    useTripAccess.mockReturnValue({
      allowed: true,
      error: '',
      loading: false,
    })

    renderPage()

    expect(screen.getByText(/viewer bisa melihat itinerary tanpa mengubah data/i)).toBeInTheDocument()
    expect(screen.getByText(/editor bisa mengubah trip, hari, dan aktivitas/i)).toBeInTheDocument()
  })
})

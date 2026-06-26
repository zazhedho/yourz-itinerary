import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'

import AppShell from './AppShell'

vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => ({ user: { id: 'user-1', name: 'User' } }),
}))

const renderShell = (initialPath) => render(
  <MemoryRouter initialEntries={[initialPath]}>
    <Routes>
      <Route element={<AppShell />}>
        <Route path="/trips" element={<div>Trips page</div>} />
        <Route path="/trips/new" element={<div>New trip page</div>} />
        <Route path="/map-picker" element={<div>Map page</div>} />
      </Route>
    </Routes>
  </MemoryRouter>,
)

describe('AppShell', () => {
  it('keeps Trips and Map as separate navigation targets', () => {
    renderShell('/trips')

    const trips = screen.getByRole('link', { name: /trips/i })
    const map = screen.getByRole('link', { name: /map/i })

    expect(trips).toHaveAttribute('href', '/trips')
    expect(map).toHaveAttribute('href', '/map-picker')
    expect(trips).toHaveClass('active')
    expect(map).not.toHaveClass('active')
  })

  it('marks only Map active on the map page', () => {
    renderShell('/map-picker')

    expect(screen.getByRole('link', { name: /trips/i })).not.toHaveClass('active')
    expect(screen.getByRole('link', { name: /map/i })).toHaveClass('active')
  })
})

import { render, screen, within } from '@testing-library/react'
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
        <Route path="/account" element={<div>Account page</div>} />
      </Route>
    </Routes>
  </MemoryRouter>,
)

describe('AppShell', () => {
  it('keeps Trips and Account as separate navigation targets', () => {
    renderShell('/trips')

    const nav = within(screen.getByRole('navigation', { name: /primary/i }))
    const trips = nav.getByRole('link', { name: /trips/i })
    const account = nav.getByRole('link', { name: /akun/i })

    expect(trips).toHaveAttribute('href', '/trips')
    expect(account).toHaveAttribute('href', '/account')
    expect(trips).toHaveClass('active')
    expect(account).not.toHaveClass('active')
  })

  it('marks only Account active on the account page', () => {
    renderShell('/account')

    const nav = within(screen.getByRole('navigation', { name: /primary/i }))
    expect(nav.getByRole('link', { name: /trips/i })).not.toHaveClass('active')
    expect(nav.getByRole('link', { name: /akun/i })).toHaveClass('active')
  })
})

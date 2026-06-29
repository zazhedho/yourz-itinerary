import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import TripMapOverview from './TripMapOverview'

vi.mock('@react-google-maps/api', () => ({
  GoogleMap: ({ children }) => <div data-testid="google-map">{children}</div>,
  useLoadScript: () => ({ isLoaded: false, loadError: null }),
}))

const days = [
  {
    id: 'day-1',
    day_number: 1,
    title: 'Day one',
    items: [{ id: 'item-1', title: 'Monas', latitude: -6.1754, longitude: 106.8272 }],
  },
  {
    id: 'day-2',
    day_number: 2,
    title: 'Day two',
    items: [{ id: 'item-2', title: 'Museum', latitude: -6.2, longitude: 106.8 }],
  },
]

describe('TripMapOverview', () => {
  it('shows pinned item list and filters by day on mobile-friendly controls', async () => {
    const user = userEvent.setup()

    render(<TripMapOverview days={days} />)

    expect(screen.getByText('Monas')).toBeInTheDocument()
    expect(screen.getByText('Museum')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: /filter day 2/i }))

    expect(screen.queryByText('Monas')).not.toBeInTheDocument()
    expect(screen.getByText('Museum')).toBeInTheDocument()
  })
})

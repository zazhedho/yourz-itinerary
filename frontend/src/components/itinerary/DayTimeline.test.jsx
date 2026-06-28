import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it } from 'vitest'

import DayTimeline from './DayTimeline'

const days = [
  {
    id: 'day-1',
    day_number: 1,
    title: 'Jakarta',
    date: '2026-06-26',
    items: [{ id: 'item-1', title: 'Monas', start_time: '09:00', cost_estimate: 0 }],
  },
]

const renderTimeline = () => render(<DayTimeline days={days} />, { wrapper: MemoryRouter })

describe('DayTimeline', () => {
  it('opens day actions without collapsing the timeline content', async () => {
    const user = userEvent.setup()
    renderTimeline()

    await user.click(screen.getByRole('button', { name: /tampilkan opsi/i }))

    expect(screen.getByRole('link', { name: /edit hari/i })).toBeInTheDocument()
    expect(screen.getByText('Monas')).toBeInTheDocument()
  })

  it('hides mutation actions for viewers', () => {
    render(<DayTimeline canEdit={false} days={days} />, { wrapper: MemoryRouter })

    expect(screen.queryByRole('button', { name: /tampilkan opsi/i })).not.toBeInTheDocument()
    expect(screen.queryByRole('link', { name: /edit aktivitas/i })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /hapus aktivitas/i })).not.toBeInTheDocument()
    expect(screen.getByText('Monas')).toBeInTheDocument()
  })
})

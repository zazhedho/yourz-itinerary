import { fireEvent, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'

import DayTimeline from './DayTimeline'

const days = [
  {
    id: 'day-1',
    day_number: 1,
    title: 'Jakarta',
    date: '2026-06-26',
    items: [{
      id: 'item-1',
      title: 'Monas',
      start_time: '09:00',
      cost_estimate: 0,
      description: 'Bawa payung dan tiket online.',
      created_by: 'user-1',
      updated_by: 'user-2',
      created_at: '2026-06-26T09:00:00+07:00',
      updated_at: '2026-06-26T10:30:00+07:00',
    }],
  },
]

const renderTimeline = () => render(<DayTimeline days={days} />, { wrapper: MemoryRouter })

const weatherDays = [{
  ...days[0],
  items: [{ ...days[0].items[0], latitude: -6.1, longitude: 106.8 }],
}]

const availableWeather = {
  'day-1': {
    data: {
      status: 'available',
      items: [{
        item_id: 'item-1',
        status: 'available',
        condition_code: 'PARTLY_CLOUDY',
        condition_description: 'Berawan sebagian',
        min_temperature_c: 24,
        max_temperature_c: 31,
        precipitation_probability: 40,
        forecast_date: '2026-06-26',
        humidity_percent: 78,
        wind_speed_kph: 12.4,
      }],
    },
  },
}

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

  it('shows compact item audit metadata with member names', () => {
    render(
      <DayTimeline
        days={days}
        memberNameByUserId={{
          'user-1': 'Zaki',
          'user-2': 'Nadia',
        }}
      />,
      { wrapper: MemoryRouter },
    )

    expect(screen.getByText(/dibuat zaki/i)).toBeInTheDocument()
    expect(screen.getByText(/diubah nadia/i)).toBeInTheDocument()
  })

  it('shows item notes in day timeline', () => {
    renderTimeline()

    expect(screen.getByText('Bawa payung dan tiket online.')).toBeInTheDocument()
  })

  it('requests weather once when a coordinate-bearing day is first expanded', async () => {
    const user = userEvent.setup()
    const onDayExpanded = vi.fn()
    render(<DayTimeline days={weatherDays} onDayExpanded={onDayExpanded} />, { wrapper: MemoryRouter })

    await user.click(screen.getByRole('button', { name: /toggle collapse/i }))
    await user.click(screen.getByRole('button', { name: /toggle collapse/i }))
    await user.click(screen.getByRole('button', { name: /toggle collapse/i }))

    expect(onDayExpanded).toHaveBeenCalledTimes(1)
  })

  it.each([
    ['out_of_range', 'Prakiraan cuaca hanya tersedia untuk 10 hari ke depan.'],
    ['past_date', 'Prakiraan cuaca untuk tanggal ini sudah tidak tersedia.'],
    ['limit_reached', 'Prakiraan cuaca sedang tidak tersedia'],
    ['provider_unavailable', 'Coba lagi'],
  ])('shows %s day state', (status, copy) => {
    render(<DayTimeline days={weatherDays} weatherByDay={{ 'day-1': { data: { status } } }} />, { wrapper: MemoryRouter })
    fireEvent.click(screen.getByRole('button', { name: /toggle collapse/i }))
    expect(screen.getByText(copy)).toBeInTheDocument()
  })

  it('shows compact weather for an available item', async () => {
    const user = userEvent.setup()
    render(<DayTimeline days={weatherDays} weatherByDay={availableWeather} />, { wrapper: MemoryRouter })
    await user.click(screen.getByRole('button', { name: /toggle collapse/i }))
    expect(screen.getByRole('button', { name: /lihat detail cuaca/i }).closest('.item-time-column')).toBeInTheDocument()
    expect(screen.getByText('Berawan sebagian')).toBeInTheDocument()
    expect(screen.getByText('24–31°C')).toBeInTheDocument()
    expect(screen.getByText('Peluang hujan 40%')).toBeInTheDocument()
  })
})

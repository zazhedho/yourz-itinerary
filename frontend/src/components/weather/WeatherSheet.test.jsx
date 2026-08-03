import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import WeatherSheet from './WeatherSheet'

const weather = {
  status: 'available',
  forecast_date: '2026-08-08',
  condition_code: 'PARTLY_CLOUDY',
  condition_description: 'Berawan sebagian',
  min_temperature_c: 24.1,
  max_temperature_c: 31.8,
  feels_like_min_c: 25,
  feels_like_max_c: 35.2,
  precipitation_probability: 40,
  humidity_percent: 78,
  wind_speed_kph: 12.4,
}

describe('WeatherSheet', () => {
  it('renders weather details with dialog semantics', () => {
    render(<WeatherSheet onClose={vi.fn()} weather={weather} />)

    expect(screen.getByRole('dialog')).toBeInTheDocument()
    expect(screen.getByText('Berawan sebagian')).toBeInTheDocument()
    expect(screen.getByText('24–32°C')).toBeInTheDocument()
    expect(screen.getByText('Peluang hujan 40%')).toBeInTheDocument()
    expect(screen.getByText('Suhu terasa')).toBeInTheDocument()
    expect(screen.getByText('Kelembapan udara')).toBeInTheDocument()
    expect(screen.getByText('78%')).toBeInTheDocument()
    expect(screen.getByText('Source: Includes weather data from Google')).toBeInTheDocument()
    expect(screen.getByText('Google Maps')).toHaveAttribute('translate', 'no')
  })

  it('closes on Escape and returns focus to opener', async () => {
    const user = userEvent.setup()
    const onClose = vi.fn()
    const opener = document.createElement('button')
    document.body.appendChild(opener)
    opener.focus()
    render(<WeatherSheet onClose={onClose} returnFocus={opener} weather={weather} />)

    await user.keyboard('{Escape}')
    expect(onClose).toHaveBeenCalledOnce()
    opener.remove()
  })
})

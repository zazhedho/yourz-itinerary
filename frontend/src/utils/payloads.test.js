import { describe, expect, it } from 'vitest'

import { buildItineraryItemPayload, buildTripPayload, emptyToUndefined, normalizeClockTime } from './payloads'

describe('payload builders', () => {
  it('removes empty optional trip fields', () => {
    expect(
      buildTripPayload({
        title: 'Bali',
        destination: '',
        start_date: '',
        end_date: '',
        timezone: 'Asia/Jakarta',
        currency_code: 'IDR',
      }),
    ).toEqual({
      title: 'Bali',
      timezone: 'Asia/Jakarta',
      currency_code: 'IDR',
    })
  })

  it('converts item coordinates and empty strings correctly', () => {
    expect(
      buildItineraryItemPayload({
        title: 'Pantai',
        description: '',
        location_name: 'Kuta',
        latitude: '-8.718',
        longitude: '115.168',
        start_time: '',
        end_time: '',
        cost_estimate: '120000',
      }),
    ).toEqual({
      title: 'Pantai',
      location_name: 'Kuta',
      latitude: -8.718,
      longitude: 115.168,
      cost_estimate: 120000,
    })
  })

  it('normalizes API clock times before submitting item payloads', () => {
    expect(normalizeClockTime('09:15:00')).toBe('09:15')
    expect(normalizeClockTime('09:15')).toBe('09:15')
    expect(normalizeClockTime('invalid')).toBe('invalid')

    expect(
      buildItineraryItemPayload({
        title: 'Sarapan',
        latitude: '',
        longitude: '',
        start_time: '09:15:00',
        end_time: '10:45:00',
        cost_estimate: '',
      }),
    ).toEqual({
      title: 'Sarapan',
      start_time: '09:15',
      end_time: '10:45',
      cost_estimate: 0,
    })
  })

  it('turns empty strings into undefined', () => {
    expect(emptyToUndefined('')).toBeUndefined()
    expect(emptyToUndefined('x')).toBe('x')
  })
})

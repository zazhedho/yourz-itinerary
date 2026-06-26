import { describe, expect, it } from 'vitest'

import { buildItineraryItemPayload, buildTripPayload, emptyToUndefined } from './payloads'

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

  it('turns empty strings into undefined', () => {
    expect(emptyToUndefined('')).toBeUndefined()
    expect(emptyToUndefined('x')).toBe('x')
  })
})

import { describe, expect, it } from 'vitest'

import { placeToItineraryLocation } from './googlePlaces'

describe('google place helpers', () => {
  it('maps selected place into itinerary location fields', () => {
    const patch = placeToItineraryLocation({
      name: 'Monas',
      formatted_address: 'Gambir, Jakarta Pusat',
      address_components: [
        { long_name: 'Gambir', types: ['administrative_area_level_3'] },
        { long_name: 'Jakarta Pusat', types: ['administrative_area_level_2'] },
      ],
      geometry: {
        location: {
          lat: () => -6.1753924,
          lng: () => 106.8271528,
        },
      },
    })

    expect(patch).toEqual({
      location_name: 'Monas, Jakarta Pusat',
      latitude: '-6.1753924',
      longitude: '106.8271528',
    })
  })

  it('falls back to formatted address when place has no name', () => {
    expect(
      placeToItineraryLocation({
        formatted_address: 'Kuta, Bali',
        geometry: {
          location: {
            lat: () => -8.718,
            lng: () => 115.168,
          },
        },
      }),
    ).toEqual({
      location_name: 'Kuta, Bali',
      latitude: '-8.718',
      longitude: '115.168',
    })
  })

  it('maps Places API New fields into itinerary location fields', () => {
    const patch = placeToItineraryLocation({
      displayName: 'Museum Nasional Indonesia',
      formattedAddress: 'Jl. Medan Merdeka Barat, Jakarta',
      addressComponents: [
        { longText: 'Gambir', types: ['administrative_area_level_3'] },
        { longText: 'Jakarta Pusat', types: ['administrative_area_level_2'] },
      ],
      location: {
        lat: () => -6.1764021,
        lng: () => 106.8215901,
      },
    })

    expect(patch).toEqual({
      location_name: 'Museum Nasional Indonesia, Jakarta Pusat',
      latitude: '-6.1764021',
      longitude: '106.8215901',
    })
  })

  it('does not duplicate city when display name already contains it', () => {
    const patch = placeToItineraryLocation({
      displayName: 'Braga, Bandung',
      addressComponents: [{ longText: 'Bandung', types: ['locality'] }],
      location: {
        lat: () => -6.9175,
        lng: () => 107.6191,
      },
    })

    expect(patch.location_name).toBe('Braga, Bandung')
  })
})

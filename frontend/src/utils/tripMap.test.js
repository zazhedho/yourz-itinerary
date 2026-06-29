import { describe, expect, it } from 'vitest'

import { getPinnedItems, getTripMapCenter } from './tripMap'

const days = [
  {
    id: 'day-1',
    day_number: 1,
    title: 'Day one',
    items: [
      { id: 'item-1', title: 'Monas', latitude: -6.1754, longitude: 106.8272 },
      { id: 'item-2', title: 'No pin' },
    ],
  },
  {
    id: 'day-2',
    day_number: 2,
    title: 'Day two',
    items: [
      { id: 'item-3', title: 'Museum', latitude: '-6.2', longitude: '106.8' },
    ],
  },
]

describe('tripMap', () => {
  it('returns pinned items with day metadata and skips items without coordinates', () => {
    expect(getPinnedItems(days)).toEqual([
      expect.objectContaining({ id: 'item-1', dayId: 'day-1', dayNumber: 1, lat: -6.1754, lng: 106.8272 }),
      expect.objectContaining({ id: 'item-3', dayId: 'day-2', dayNumber: 2, lat: -6.2, lng: 106.8 }),
    ])
  })

  it('filters pinned items by day', () => {
    expect(getPinnedItems(days, 'day-2')).toEqual([
      expect.objectContaining({ id: 'item-3' }),
    ])
  })

  it('uses first pinned item as map center', () => {
    expect(getTripMapCenter(getPinnedItems(days))).toEqual({ lat: -6.1754, lng: 106.8272 })
  })
})

const defaultCenter = { lat: -6.2, lng: 106.816666 }

export const getPinnedItems = (days = [], dayId = 'all') =>
  days
    .filter((day) => dayId === 'all' || day.id === dayId)
    .flatMap((day) => (day.items || []).map((item) => {
      const lat = Number(item.latitude)
      const lng = Number(item.longitude)
      if (!Number.isFinite(lat) || !Number.isFinite(lng)) return null
      return {
        ...item,
        dayId: day.id,
        dayNumber: day.day_number,
        dayTitle: day.title,
        lat,
        lng,
      }
    }))
    .filter(Boolean)

export const getTripMapCenter = (items = []) => {
  const first = items[0]
  if (!first) return defaultCenter
  return { lat: first.lat, lng: first.lng }
}

export const placeToItineraryLocation = (place) => {
  const location = place?.location || place?.geometry?.location
  const locationName = place?.displayName || place?.name || place?.formattedAddress || place?.formatted_address || ''

  if (!location) {
    return { location_name: locationName }
  }

  return {
    location_name: locationName,
    latitude: String(location.lat()),
    longitude: String(location.lng()),
  }
}

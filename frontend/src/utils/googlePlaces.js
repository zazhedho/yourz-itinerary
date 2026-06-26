const areaTypes = ['locality', 'administrative_area_level_2']

const getComponentName = (component) => component?.longText || component?.long_name || component?.shortText || component?.short_name || ''

const getPlaceAreaName = (place) => {
  const components = place?.addressComponents || place?.address_components || []
  const area = components.find((component) => areaTypes.some((type) => component?.types?.includes(type)))
  return getComponentName(area)
}

const appendAreaName = (name, areaName) => {
  if (!name || !areaName) return name
  if (name.toLowerCase().includes(areaName.toLowerCase())) return name
  return `${name}, ${areaName}`
}

export const placeToItineraryLocation = (place) => {
  const location = place?.location || place?.geometry?.location
  const rawLocationName = place?.displayName || place?.name || place?.formattedAddress || place?.formatted_address || ''
  const locationName = appendAreaName(rawLocationName, getPlaceAreaName(place))

  if (!location) {
    return { location_name: locationName }
  }

  return {
    location_name: locationName,
    latitude: String(location.lat()),
    longitude: String(location.lng()),
  }
}

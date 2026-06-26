import { useEffect, useRef } from 'react'

const GooglePlaceAutocomplete = ({ isLoaded, name, placeholder, disabled, value, onChange, onPlaceSelect }) => {
  const inputRef = useRef(null)
  const autocompleteRef = useRef(null)
  const onPlaceSelectRef = useRef(onPlaceSelect)

  useEffect(() => {
    onPlaceSelectRef.current = onPlaceSelect
  }, [onPlaceSelect])

  useEffect(() => {
    if (!isLoaded || !inputRef.current) return

    const { Autocomplete } = window.google.maps.places
    if (!Autocomplete) return

    autocompleteRef.current = new Autocomplete(inputRef.current, {
      fields: ['name', 'formatted_address', 'geometry', 'address_components'],
    })

    const listener = autocompleteRef.current.addListener('place_changed', () => {
      const place = autocompleteRef.current.getPlace()
      if (!place || !place.geometry) return

      // Extract City/Regency name for a concise format (e.g. "Braga, Bandung")
      const cityComponent = place.address_components?.find((c) =>
        c.types.includes('locality') || c.types.includes('administrative_area_level_2')
      )
      const city = cityComponent ? cityComponent.short_name : ''
      
      let displayName = place.name || ''
      if (city && !displayName.toLowerCase().includes(city.toLowerCase())) {
        displayName = `${displayName}, ${city}`
      }
      
      // Fallback if somehow name is missing
      if (!displayName) displayName = place.formatted_address || ''

      onPlaceSelectRef.current({
        displayName,
        formattedAddress: place.formatted_address,
        location: {
          lat: () => place.geometry.location.lat(),
          lng: () => place.geometry.location.lng(),
        },
      })
    })

    return () => {
      if (window.google.maps.event) {
        window.google.maps.event.removeListener(listener)
      }
    }
  }, [isLoaded])

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault()
    }
  }

  return (
    <input
      ref={inputRef}
      name={name}
      type="text"
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      onKeyDown={handleKeyDown}
      disabled={disabled}
      className="google-place-autocomplete-input"
      style={{ width: '100%' }}
    />
  )
}

export default GooglePlaceAutocomplete

import { useEffect, useRef, useState } from 'react'

const GooglePlaceAutocomplete = ({ disabled, isLoaded, name, onChange, onPlaceSelect, placeholder, value }) => {
  const inputRef = useRef(null)
  const onPlaceSelectRef = useRef(onPlaceSelect)
  const [error, setError] = useState('')

  useEffect(() => {
    onPlaceSelectRef.current = onPlaceSelect
  }, [onPlaceSelect])

  useEffect(() => {
    if (!isLoaded || !inputRef.current) return undefined

    let active = true
    let autocomplete = null
    let listener = null

    const initAutocomplete = async () => {
      try {
        // Fallback or explicit check if places library is available
        if (!window.google?.maps?.places?.Autocomplete) {
          await window.google.maps.importLibrary('places')
        }
        
        if (!active || !inputRef.current) return

        autocomplete = new window.google.maps.places.Autocomplete(inputRef.current, {
          fields: ['name', 'geometry', 'address_components', 'formatted_address'],
        })

        // Prevent the autocomplete from closing the keyboard and becoming full screen by keeping it contained in the native input dropdown.
        listener = autocomplete.addListener('place_changed', () => {
          const place = autocomplete.getPlace()
          
          let cityName = ''
          if (place.address_components) {
            const city = place.address_components.find(c => 
              c.types.includes('locality') || c.types.includes('administrative_area_level_2')
            )
            if (city) cityName = city.long_name
          }
          
          const displayName = place.name || ''
          const locationName = cityName && !displayName.includes(cityName) 
            ? `${displayName}, ${cityName}`
            : displayName

          // Call parent handler
          onPlaceSelectRef.current({
            displayName: locationName,
            formattedAddress: place.formatted_address,
            location: place.geometry?.location
          })
        })
      } catch (err) {
        if (active) setError('Pencarian lokasi gagal dimuat.')
      }
    }

    initAutocomplete()

    return () => {
      active = false
      if (listener) {
        window.google.maps.event.removeListener(listener)
      }
      // Cleanup the google maps autocomplete bindings on the input element if possible
      if (inputRef.current) {
        const pacContainer = document.querySelector('.pac-container')
        if (pacContainer) pacContainer.remove()
      }
    }
  }, [isLoaded])

  return (
    <>
      <input
        ref={inputRef}
        name={name}
        type="text"
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        disabled={disabled}
      />
      {error ? <span className="field-note error">{error}</span> : null}
    </>
  )
}

export default GooglePlaceAutocomplete

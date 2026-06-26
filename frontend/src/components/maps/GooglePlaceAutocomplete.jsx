import { useEffect, useRef, useState } from 'react'

const GooglePlaceAutocomplete = ({ disabled, isLoaded, name, onChange, onPlaceSelect, placeholder, value }) => {
  const containerRef = useRef(null)
  const onPlaceSelectRef = useRef(onPlaceSelect)
  const [error, setError] = useState('')

  useEffect(() => {
    onPlaceSelectRef.current = onPlaceSelect
  }, [onPlaceSelect])

  useEffect(() => {
    if (!isLoaded || !containerRef.current) return undefined

    let active = true
    let autocompleteElement = null

    const initAutocomplete = async () => {
      try {
        const { PlaceAutocompleteElement } = await window.google.maps.importLibrary('places')
        if (!active || !containerRef.current) return

        autocompleteElement = new PlaceAutocompleteElement()
        autocompleteElement.placeholder = 'Cari tempat di Google Maps'
        autocompleteElement.className = 'google-place-autocomplete'
        autocompleteElement.addEventListener('gmp-select', async ({ placePrediction }) => {
          const place = placePrediction.toPlace()
          await place.fetchFields({ fields: ['displayName', 'formattedAddress', 'location'] })
          onPlaceSelectRef.current(place)
        })

        containerRef.current.replaceChildren(autocompleteElement)
      } catch {
        if (active) setError('Pencarian lokasi gagal dimuat.')
      }
    }

    initAutocomplete()

    return () => {
      active = false
      autocompleteElement?.remove()
    }
  }, [isLoaded])

  return (
    <>
      {isLoaded ? <div ref={containerRef} className="google-place-autocomplete-host" /> : null}
      <input
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

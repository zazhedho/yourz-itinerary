import { useEffect, useRef, useState } from 'react'

import { placeToItineraryLocation } from '../../utils/googlePlaces'

const minSearchLength = 3

const getPredictionLabel = (prediction) => prediction?.text?.toString?.() || prediction?.mainText?.toString?.() || ''

const GooglePlaceAutocomplete = ({ disabled, isLoaded, name, onChange, onPlaceSelect, placeholder, value }) => {
  const onPlaceSelectRef = useRef(onPlaceSelect)
  const requestRef = useRef(0)
  const sessionTokenRef = useRef(null)
  const [hasUserTyped, setHasUserTyped] = useState(false)
  const [placesApi, setPlacesApi] = useState(null)
  const [selectedValue, setSelectedValue] = useState('')
  const [suggestions, setSuggestions] = useState([])
  const [error, setError] = useState('')

  useEffect(() => {
    onPlaceSelectRef.current = onPlaceSelect
  }, [onPlaceSelect])

  useEffect(() => {
    if (!isLoaded) return undefined

    let active = true

    const loadPlacesApi = async () => {
      try {
        const api = await window.google.maps.importLibrary('places')
        if (active) setPlacesApi(api)
      } catch {
        if (active) setError('Pencarian lokasi gagal dimuat.')
      }
    }

    loadPlacesApi()

    return () => {
      active = false
    }
  }, [isLoaded])

  useEffect(() => {
    if (!hasUserTyped || !placesApi || disabled || value.trim().length < minSearchLength) {
      return undefined
    }
    if (value === selectedValue) {
      return undefined
    }

    const requestId = requestRef.current + 1
    requestRef.current = requestId

    const timeoutId = window.setTimeout(async () => {
      try {
        sessionTokenRef.current = sessionTokenRef.current || new placesApi.AutocompleteSessionToken()
        const response = await placesApi.AutocompleteSuggestion.fetchAutocompleteSuggestions({
          input: value,
          sessionToken: sessionTokenRef.current,
        })

        if (requestRef.current !== requestId) return
        setSuggestions(response.suggestions || [])
      } catch {
        if (requestRef.current === requestId) {
          setSuggestions([])
          setError('Pencarian lokasi gagal.')
        }
      }
    }, 250)

    return () => window.clearTimeout(timeoutId)
  }, [disabled, hasUserTyped, placesApi, selectedValue, value])

  const selectSuggestion = async (suggestion) => {
    const prediction = suggestion?.placePrediction
    if (!prediction) return

    setSuggestions([])
    try {
      const place = prediction.toPlace()
      await place.fetchFields({ fields: ['displayName', 'formattedAddress', 'location', 'addressComponents'] })
      sessionTokenRef.current = null
      setSelectedValue(placeToItineraryLocation(place).location_name)
      setHasUserTyped(false)
      onPlaceSelectRef.current(place)
    } catch {
      setError('Gagal mengambil detail lokasi.')
    }
  }

  const handleChange = (event) => {
    setHasUserTyped(true)
    setSelectedValue('')
    if (event.target.value.trim().length < minSearchLength) setSuggestions([])
    onChange(event)
  }

  const visibleSuggestions =
    hasUserTyped && !disabled && value.trim().length >= minSearchLength && value !== selectedValue ? suggestions : []

  const handleKeyDown = (event) => {
    if (event.key !== 'Enter' || !visibleSuggestions.length) return
    event.preventDefault()
    selectSuggestion(visibleSuggestions[0])
  }

  return (
    <div className="place-search-field">
      <input
        name={name}
        type="text"
        placeholder={placeholder}
        value={value}
        onChange={handleChange}
        onKeyDown={handleKeyDown}
        disabled={disabled}
        autoComplete="off"
      />
      {visibleSuggestions.length > 0 && (
        <div className="place-suggestions" role="listbox">
          {visibleSuggestions.map((suggestion) => {
            const prediction = suggestion.placePrediction
            const label = getPredictionLabel(prediction)
            if (!label) return null

            return (
              <button
                key={prediction.placeId}
                className="place-suggestion"
                onMouseDown={(event) => event.preventDefault()}
                onClick={() => selectSuggestion(suggestion)}
                type="button"
              >
                {label}
              </button>
            )
          })}
        </div>
      )}
      {error ? <span className="field-note error">{error}</span> : null}
    </div>
  )
}

export default GooglePlaceAutocomplete

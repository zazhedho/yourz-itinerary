import { GoogleMap } from '@react-google-maps/api'
import { useEffect, useMemo, useRef, useState } from 'react'

import { getGoogleMapsMapId } from '../../utils/runtimeConfig'

const defaultCenter = { lat: -6.2, lng: 106.816666 }
const mapContainerStyle = { width: '100%', height: '100%' }

const CoordinatePicker = ({ latitude, longitude, onPick }) => {
  const markerRef = useRef(null)
  const [map, setMap] = useState(null)
  const mapId = getGoogleMapsMapId()
  const lat = Number(latitude)
  const lng = Number(longitude)
  const hasPosition =
    latitude !== '' && longitude !== '' && latitude !== null && longitude !== null && Number.isFinite(lat) && Number.isFinite(lng)
  const position = useMemo(() => (hasPosition ? { lat, lng } : defaultCenter), [hasPosition, lat, lng])

  useEffect(() => {
    if (!map) return undefined

    let active = true

    const syncMarker = async () => {
      if (!hasPosition) {
        if (markerRef.current) markerRef.current.map = null
        markerRef.current = null
        return
      }

      const { AdvancedMarkerElement } = await window.google.maps.importLibrary('marker')
      if (!active) return

      if (!markerRef.current) {
        markerRef.current = new AdvancedMarkerElement({
          map,
          position,
        })
        return
      }

      markerRef.current.position = position
      markerRef.current.map = map
    }

    syncMarker()

    return () => {
      active = false
    }
  }, [hasPosition, map, position])

  const handleUnmount = () => {
    if (markerRef.current) markerRef.current.map = null
    markerRef.current = null
    setMap(null)
  }

  return (
    <div className="map-panel compact">
      <div className="map-canvas" style={{ position: 'relative', width: '100%', height: '100%', minHeight: '220px' }}>
        <GoogleMap
          mapContainerStyle={mapContainerStyle}
          center={position}
          zoom={hasPosition ? 15 : 12}
          onLoad={setMap}
          onUnmount={handleUnmount}
          onClick={(e) => {
            onPick({
              latitude: Number(e.latLng.lat().toFixed(7)),
              longitude: Number(e.latLng.lng().toFixed(7)),
            })
          }}
          options={{
            disableDefaultUI: true,
            zoomControl: true,
            clickableIcons: false,
            mapId,
          }}
        />
      </div>
    </div>
  )
}

export default CoordinatePicker

import { GoogleMap, useLoadScript } from '@react-google-maps/api'
import { MapPin } from 'lucide-react'
import { useEffect, useMemo, useRef, useState } from 'react'

import { getPinnedItems, getTripMapCenter } from '../../utils/tripMap'
import { getGoogleMapsApiKey, getGoogleMapsMapId } from '../../utils/runtimeConfig'

const libraries = ['places']
const mapContainerStyle = { width: '100%', height: '100%' }

const TripMapOverviewContent = ({ apiKey, days = [] }) => {
  const markerRefs = useRef([])
  const [dayId, setDayId] = useState('all')
  const [map, setMap] = useState(null)
  const [selectedItemId, setSelectedItemId] = useState('')
  const mapId = getGoogleMapsMapId()
  const items = useMemo(() => getPinnedItems(days, dayId), [days, dayId])
  const center = useMemo(() => getTripMapCenter(items), [items])
  const selectedItem = items.find((item) => item.id === selectedItemId) || items[0]
  const { isLoaded, loadError } = useLoadScript({
    googleMapsApiKey: apiKey,
    libraries,
  })

  useEffect(() => {
    if (!map || !isLoaded) return undefined

    let active = true

    const syncMarkers = async () => {
      markerRefs.current.forEach((marker) => {
        marker.map = null
      })
      markerRefs.current = []

      const { AdvancedMarkerElement } = await window.google.maps.importLibrary('marker')
      if (!active) return

      markerRefs.current = items.map((item) => {
        const marker = new AdvancedMarkerElement({
          map,
          position: { lat: item.lat, lng: item.lng },
          title: item.title,
        })
        marker.addListener('click', () => setSelectedItemId(item.id))
        return marker
      })
    }

    syncMarkers()
    map.panTo(center)

    return () => {
      active = false
      markerRefs.current.forEach((marker) => {
        marker.map = null
      })
      markerRefs.current = []
    }
  }, [center, isLoaded, items, map])

  return (
    <section className="content-section trip-map-section">
      <div className="section-heading">
        <div>
          <p className="eyebrow">Map</p>
          <h2>Pin lokasi</h2>
        </div>
      </div>
      <div className="map-day-filter" aria-label="Filter pin berdasarkan hari" role="group">
        <button aria-label="Filter semua hari" className={dayId === 'all' ? 'active' : ''} onClick={() => setDayId('all')} type="button">
          Semua
        </button>
        {days.map((day) => (
          <button aria-label={`Filter Day ${day.day_number}`} className={dayId === day.id ? 'active' : ''} key={day.id} onClick={() => setDayId(day.id)} type="button">
            Day {day.day_number}
          </button>
        ))}
      </div>
      <div className="trip-map-layout">
        <div className="trip-map-canvas">
          {loadError ? (
            <div className="map-placeholder">Peta gagal dimuat. Cek koneksi dan Google Maps config.</div>
          ) : isLoaded ? (
            <GoogleMap
              center={center}
              mapContainerStyle={mapContainerStyle}
              onLoad={setMap}
              onUnmount={() => setMap(null)}
              options={{
                clickableIcons: false,
                disableDefaultUI: true,
                mapId,
                zoomControl: true,
              }}
              zoom={items.length ? 13 : 11}
            />
          ) : (
            <div className="map-placeholder">Memuat peta...</div>
          )}
        </div>
        <div className="map-pin-list">
          {items.length ? (
            items.map((item) => (
              <button
                className={selectedItem?.id === item.id ? 'active' : ''}
                key={item.id}
                onClick={() => setSelectedItemId(item.id)}
                type="button"
              >
                <MapPin size={15} />
                <span>
                  <strong>{item.title}</strong>
                  <small>Day {item.dayNumber}{item.location_name ? ` • ${item.location_name}` : ''}</small>
                </span>
              </button>
            ))
          ) : (
            <div className="map-empty">Belum ada aktivitas dengan koordinat.</div>
          )}
        </div>
      </div>
    </section>
  )
}

const TripMapOverview = ({ days = [] }) => {
  const apiKey = getGoogleMapsApiKey()
  if (!apiKey) return null
  return <TripMapOverviewContent apiKey={apiKey} days={days} />
}

export default TripMapOverview

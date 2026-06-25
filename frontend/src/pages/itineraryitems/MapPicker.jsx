import { MapContainer, Marker, TileLayer, useMapEvents } from 'react-leaflet'
import { useState } from 'react'

const ClickMarker = ({ onPick, position }) => {
  useMapEvents({
    click(event) {
      onPick(event.latlng)
    },
  })

  return position ? <Marker position={position} /> : null
}

const MapPicker = () => {
  const [position, setPosition] = useState(null)

  return (
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">Map pin</p>
          <h2>Pilih koordinat lokasi</h2>
        </div>
      </div>
      <div className="map-panel">
        <MapContainer center={[-6.2, 106.816666]} className="map-canvas" zoom={12}>
          <TileLayer
            attribution="&copy; OpenStreetMap contributors"
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <ClickMarker onPick={setPosition} position={position} />
        </MapContainer>
      </div>
      <div className="empty-card">
        {position
          ? `Latitude ${position.lat.toFixed(7)}, Longitude ${position.lng.toFixed(7)}`
          : 'Tap map untuk memilih pin. Koordinat bisa dipindahkan ke form aktivitas.'}
      </div>
    </section>
  )
}

export default MapPicker

import { CircleMarker, MapContainer, TileLayer, useMapEvents } from 'react-leaflet'

const defaultCenter = [-6.2, 106.816666]

const ClickLayer = ({ onPick }) => {
  useMapEvents({
    click(event) {
      onPick({
        latitude: Number(event.latlng.lat.toFixed(7)),
        longitude: Number(event.latlng.lng.toFixed(7)),
      })
    },
  })

  return null
}

const CoordinatePicker = ({ latitude, longitude, onPick }) => {
  const hasPosition = latitude !== '' && longitude !== '' && latitude !== null && longitude !== null
  const position = hasPosition ? [Number(latitude), Number(longitude)] : defaultCenter

  return (
    <div className="map-panel compact">
      <MapContainer center={position} className="map-canvas" zoom={hasPosition ? 15 : 12}>
        <TileLayer attribution="&copy; OpenStreetMap contributors" url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
        <ClickLayer onPick={onPick} />
        {hasPosition && (
          <CircleMarker center={position} fillColor="#ff385c" fillOpacity={0.92} pathOptions={{ color: '#ffffff' }} radius={10} />
        )}
      </MapContainer>
    </div>
  )
}

export default CoordinatePicker

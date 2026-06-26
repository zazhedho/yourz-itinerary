import { useState } from 'react'

import CoordinatePicker from '../../components/maps/CoordinatePicker'

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
      <CoordinatePicker latitude={position?.latitude || ''} longitude={position?.longitude || ''} onPick={setPosition} />
      <div className="empty-card">
        {position
          ? `Latitude ${position.latitude.toFixed(7)}, Longitude ${position.longitude.toFixed(7)}`
          : 'Tap map untuk memilih pin. Koordinat bisa dipindahkan ke form aktivitas.'}
      </div>
    </section>
  )
}

export default MapPicker

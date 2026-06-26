import { useLoadScript } from '@react-google-maps/api'
import { MapPin } from 'lucide-react'
import { lazy, Suspense, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import GooglePlaceAutocomplete from '../../components/maps/GooglePlaceAutocomplete'
import { getErrorMessage } from '../../services/api'
import itineraryItemService from '../../services/itineraryItemService'
import { placeToItineraryLocation } from '../../utils/googlePlaces'
import { buildItineraryItemPayload, normalizeClockTime } from '../../utils/payloads'
import { getGoogleMapsApiKey } from '../../utils/runtimeConfig'

const CoordinatePicker = lazy(() => import('../../components/maps/CoordinatePicker'))

const libraries = ['places']

const ItineraryItemForm = () => {
  const { dayId, itemId } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const existingItem = location.state?.item
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({
    title: existingItem?.title || '',
    description: existingItem?.description || '',
    location_name: existingItem?.location_name || '',
    latitude: existingItem?.latitude ?? '',
    longitude: existingItem?.longitude ?? '',
    start_time: normalizeClockTime(existingItem?.start_time || ''),
    end_time: normalizeClockTime(existingItem?.end_time || ''),
    cost_estimate: existingItem?.cost_estimate || 0,
  })

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  const { isLoaded, loadError } = useLoadScript({
    googleMapsApiKey: getGoogleMapsApiKey(),
    libraries,
  })

  const handlePlaceSelect = (place) => {
    setForm((current) => ({
      ...current,
      ...placeToItineraryLocation(place),
    }))
  }

  const handleCoordinatePick = ({ latitude, longitude }) => {
    setForm((current) => ({
      ...current,
      latitude: String(latitude),
      longitude: String(longitude),
    }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      const payload = buildItineraryItemPayload(form)
      if (itemId) await itineraryItemService.update(itemId, payload)
      else await itineraryItemService.create(dayId, payload)
      navigate(-1)
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal menyimpan item'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">Itinerary item</p>
          <h2>{itemId ? 'Edit aktivitas' : 'Tambah aktivitas'}</h2>
        </div>
      </div>
      <form className="form-card" onSubmit={handleSubmit}>
        <ErrorBanner message={error} />
        <label>
          Judul aktivitas
          <input name="title" value={form.title} onChange={handleChange} placeholder="Contoh: Makan Siang di Restoran Lokal" required />
        </label>
        <label>
          Lokasi
          {loadError ? <span className="field-note error">Google Maps gagal dimuat. Cek API key, billing, dan Places API.</span> : null}
          <GooglePlaceAutocomplete
            isLoaded={isLoaded}
            name="location_name"
            placeholder={loadError ? 'Isi lokasi manual' : 'Ketik nama tempat...'}
            value={form.location_name}
            onChange={handleChange}
            onPlaceSelect={handlePlaceSelect}
            disabled={!isLoaded && !loadError}
          />
        </label>
        <div className="form-grid">
          <label>
            Latitude
            <input name="latitude" type="number" step="0.0000001" value={form.latitude} onChange={handleChange} />
          </label>
          <label>
            Longitude
            <input name="longitude" type="number" step="0.0000001" value={form.longitude} onChange={handleChange} />
          </label>
        </div>
        <div className="map-field">
          <div className="map-field-header">
            <span>
              <MapPin size={17} />
              Pin lokasi
            </span>
            <small>Tap map untuk isi koordinat</small>
          </div>
          <Suspense fallback={<Loading label="Memuat map..." />}>
            {isLoaded ? (
              <CoordinatePicker latitude={form.latitude} longitude={form.longitude} onPick={handleCoordinatePick} />
            ) : loadError ? (
              <div className="map-placeholder">Peta belum tersedia. Latitude dan longitude masih bisa diisi manual.</div>
            ) : (
              <div className="map-placeholder">Menunggu peta siap...</div>
            )}
          </Suspense>
        </div>
        <div className="form-grid">
          <label>
            Mulai
            <input name="start_time" type="time" value={form.start_time} onChange={handleChange} />
          </label>
          <label>
            Selesai
            <input name="end_time" type="time" value={form.end_time} onChange={handleChange} />
          </label>
        </div>
        <label>
          Estimasi biaya
          <input name="cost_estimate" min="0" type="number" value={form.cost_estimate} onChange={handleChange} placeholder="Contoh: 150000" />
        </label>
        <label>
          Catatan
          <textarea name="description" value={form.description} onChange={handleChange} placeholder="Contoh: Bawa uang tunai, reservasi atas nama Budi." rows="3" />
        </label>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Menyimpan...' : 'Simpan aktivitas'}
        </button>
      </form>
    </section>
  )
}

export default ItineraryItemForm

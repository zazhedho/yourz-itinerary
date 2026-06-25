import { MapPin } from 'lucide-react'
import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import { getErrorMessage } from '../../services/api'
import itineraryItemService from '../../services/itineraryItemService'

const ItineraryItemForm = () => {
  const { dayId, itemId } = useParams()
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({
    title: '',
    description: '',
    location_name: '',
    latitude: '',
    longitude: '',
    start_time: '',
    end_time: '',
    cost_estimate: 0,
  })

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  const payload = {
    ...form,
    latitude: form.latitude === '' ? null : Number(form.latitude),
    longitude: form.longitude === '' ? null : Number(form.longitude),
    cost_estimate: Number(form.cost_estimate || 0),
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    try {
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
          <input name="title" value={form.title} onChange={handleChange} required />
        </label>
        <label>
          Lokasi
          <input name="location_name" value={form.location_name} onChange={handleChange} />
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
        <Link className="button-secondary" to="/map-picker">
          <MapPin size={17} />
          Pilih dari map
        </Link>
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
          <input name="cost_estimate" min="0" type="number" value={form.cost_estimate} onChange={handleChange} />
        </label>
        <label>
          Catatan
          <textarea name="description" value={form.description} onChange={handleChange} />
        </label>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Menyimpan...' : 'Simpan aktivitas'}
        </button>
      </form>
    </section>
  )
}

export default ItineraryItemForm

import { useMemo, useState } from 'react'
import { Calendar } from 'lucide-react'
import { useNavigate, useParams } from 'react-router-dom'

import AccessDenied from '../../components/common/AccessDenied'
import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import useTripAccess from '../../hooks/useTripAccess'
import tripService from '../../services/tripService'
import { getErrorMessage } from '../../services/api'
import { buildTripPayload } from '../../utils/payloads'

const TripForm = () => {
  const { tripId } = useParams()
  const navigate = useNavigate()
  const { allowed, error: accessError, loading, trip } = useTripAccess(tripId, 'edit')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [draft, setDraft] = useState({})
  const baseForm = useMemo(() => ({
    title: '',
    destination: '',
    start_date: '',
    end_date: '',
    timezone: 'Asia/Jakarta',
    currency_code: 'IDR',
    ...(tripId && trip ? {
      title: trip.title || '',
      destination: trip.destination || '',
      start_date: trip.start_date || '',
      end_date: trip.end_date || '',
      timezone: trip.timezone || 'Asia/Jakarta',
      currency_code: trip.currency_code || 'IDR',
    } : {}),
  }), [tripId, trip])
  const form = { ...baseForm, ...draft }

  const handleChange = (event) => {
    const { name, value } = event.target
    setDraft((current) => {
      const next = { ...form, ...current }
      if (name === 'start_date' && next.end_date && value && next.end_date < value) {
        return { ...current, start_date: value, end_date: value }
      }
      return { ...current, [name]: value }
    })
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    if (form.start_date && form.end_date && form.end_date < form.start_date) {
      setError('Tanggal selesai tidak boleh lebih kecil dari tanggal mulai')
      return
    }
    setSubmitting(true)
    setError('')
    try {
      const payload = buildTripPayload(form)
      const response = tripId ? await tripService.update(tripId, payload) : await tripService.create(payload)
      const id = response?.data?.data?.id || tripId
      navigate(`/trips/${id}`)
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal menyimpan trip'))
    } finally {
      setSubmitting(false)
    }
  }

  if (tripId && loading) return <Loading label="Memeriksa akses trip..." />
  if (tripId && !allowed) return <AccessDenied backTo={`/trips/${tripId}`} message={accessError || 'Hanya owner dan editor yang bisa mengubah trip.'} />

  return (
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">{tripId ? 'Edit trip' : 'Trip baru'}</p>
          <h2>{tripId ? 'Update itinerary' : 'Buat itinerary'}</h2>
        </div>
      </div>
      <form className="form-card" onSubmit={handleSubmit}>
        <ErrorBanner message={error} />
        <label>
          Judul trip
          <input name="title" value={form.title} onChange={handleChange} placeholder="Contoh: Liburan Musim Panas di Bali" required />
        </label>
        <label>
          Destinasi
          <input name="destination" value={form.destination} onChange={handleChange} placeholder="Contoh: Bali, Indonesia" />
        </label>
        <div className="date-range-row">
          <label>
            Mulai
            <div className="date-input-wrapper">
              <Calendar className="date-input-icon" size={16} />
              {!form.start_date && <span className="date-placeholder">dd/mm/yyyy</span>}
              <input 
                name="start_date" 
                type="date"
                className={!form.start_date ? 'empty-date' : ''}
                value={form.start_date} 
                onChange={handleChange} 
              />
            </div>
          </label>
          <span className="date-separator">-</span>
          <label>
            Selesai
            <div className="date-input-wrapper">
              <Calendar className="date-input-icon" size={16} />
              {!form.end_date && <span className="date-placeholder">dd/mm/yyyy</span>}
              <input
                min={form.start_date || undefined}
                name="end_date"
                type="date"
                className={!form.end_date ? 'empty-date' : ''}
                value={form.end_date}
                onChange={handleChange}
              />
            </div>
          </label>
        </div>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Menyimpan...' : 'Simpan trip'}
        </button>
      </form>
    </section>
  )
}

export default TripForm

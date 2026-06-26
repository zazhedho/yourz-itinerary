import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import tripService from '../../services/tripService'
import { getErrorMessage, getResponseData } from '../../services/api'
import { buildTripPayload } from '../../utils/payloads'

const TripForm = () => {
  const { tripId } = useParams()
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({
    title: '',
    destination: '',
    start_date: '',
    end_date: '',
    timezone: 'Asia/Jakarta',
    currency_code: 'IDR',
  })

  useEffect(() => {
    if (!tripId) return
    tripService.getById(tripId).then((response) => {
      const trip = getResponseData(response)
      setForm({
        title: trip.title || '',
        destination: trip.destination || '',
        start_date: trip.start_date || '',
        end_date: trip.end_date || '',
        timezone: trip.timezone || 'Asia/Jakarta',
        currency_code: trip.currency_code || 'IDR',
      })
    })
  }, [tripId])

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
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
          <input name="title" value={form.title} onChange={handleChange} required />
        </label>
        <label>
          Destinasi
          <input name="destination" value={form.destination} onChange={handleChange} />
        </label>
        <div className="form-grid">
          <label>
            Mulai
            <input name="start_date" type="date" value={form.start_date} onChange={handleChange} />
          </label>
          <label>
            Selesai
            <input name="end_date" type="date" value={form.end_date} onChange={handleChange} />
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

import { useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import itineraryDayService from '../../services/itineraryDayService'
import { getErrorMessage } from '../../services/api'
import { buildItineraryDayPayload } from '../../utils/payloads'

const ItineraryDayForm = () => {
  const { tripId, dayId } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const existingDay = location.state?.day
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({
    day_number: existingDay?.day_number || 1,
    title: existingDay?.title || '',
    date: existingDay?.date || '',
  })

  const handleChange = (event) => {
    const value = event.target.name === 'day_number' ? Number(event.target.value) : event.target.value
    setForm((current) => ({ ...current, [event.target.name]: value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      const payload = buildItineraryDayPayload(form)
      if (dayId) {
        await itineraryDayService.update(dayId, payload)
        navigate(-1)
      } else {
        await itineraryDayService.create(tripId, payload)
        navigate(`/trips/${tripId}`)
      }
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal menyimpan hari itinerary'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">Itinerary day</p>
          <h2>{dayId ? 'Edit hari' : 'Tambah hari'}</h2>
        </div>
      </div>
      <form className="form-card" onSubmit={handleSubmit}>
        <ErrorBanner message={error} />
        <label>
          Nomor hari
          <input min="1" name="day_number" type="number" value={form.day_number} onChange={handleChange} required />
        </label>
        <label>
          Judul
          <input name="title" value={form.title} onChange={handleChange} />
        </label>
        <label>
          Tanggal
          <input name="date" type="date" value={form.date} onChange={handleChange} />
        </label>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Menyimpan...' : 'Simpan hari'}
        </button>
      </form>
    </section>
  )
}

export default ItineraryDayForm

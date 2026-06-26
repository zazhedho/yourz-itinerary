import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import { getErrorMessage } from '../../services/api'
import tripMemberService from '../../services/tripMemberService'
import { buildTripMemberPayload } from '../../utils/payloads'

const TripMemberAdd = () => {
  const { tripId } = useParams()
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({ email: '', role: 'viewer' })

  const handleSubmit = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      const payload = buildTripMemberPayload(form)
      await tripMemberService.addMember(tripId, payload)
      navigate(`/trips/${tripId}`)
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal menambahkan member'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">Member</p>
          <h2>Tambah member lewat email</h2>
        </div>
      </div>
      <form className="form-card" onSubmit={handleSubmit}>
        <ErrorBanner message={error} />
        <label>
          Email member
          <input
            type="email"
            value={form.email}
            onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))}
            placeholder="Contoh: teman@email.com"
            required
          />
        </label>
        <label>
          Role
          <select value={form.role} onChange={(event) => setForm((current) => ({ ...current, role: event.target.value }))}>
            <option value="viewer">Viewer</option>
            <option value="editor">Editor</option>
          </select>
        </label>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Menambahkan...' : 'Tambah member'}
        </button>
      </form>
    </section>
  )
}

export default TripMemberAdd

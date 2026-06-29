import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import AccessDenied from '../../components/common/AccessDenied'
import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import useTripAccess from '../../hooks/useTripAccess'
import { getErrorMessage } from '../../services/api'
import tripMemberService from '../../services/tripMemberService'
import { buildTripMemberPayload } from '../../utils/payloads'

const getAddMemberError = (error) => {
  const message = getErrorMessage(error, 'Gagal menambahkan member')
  const normalized = message.toLowerCase()
  if (normalized.includes('user') && normalized.includes('not found')) return 'Email belum terdaftar. Minta member daftar dulu sebelum diundang.'
  if (normalized.includes('duplicate') || normalized.includes('already')) return 'Member sudah ada di itinerary ini.'
  if (normalized.includes('role')) return 'Role tidak valid. Pilih Viewer atau Editor.'
  return message
}

const TripMemberAdd = () => {
  const { tripId } = useParams()
  const navigate = useNavigate()
  const { allowed, error: accessError, loading } = useTripAccess(tripId, 'manageMembers')
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
      setError(getAddMemberError(err))
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) return <Loading label="Memeriksa akses member..." />
  if (!allowed) return <AccessDenied backTo={`/trips/${tripId}`} message={accessError || 'Hanya owner yang bisa menambah member.'} />

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
        <div className="role-help">
          <span>Viewer bisa melihat itinerary tanpa mengubah data.</span>
          <span>Editor bisa mengubah trip, hari, dan aktivitas.</span>
        </div>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Menambahkan...' : 'Tambah member'}
        </button>
      </form>
    </section>
  )
}

export default TripMemberAdd

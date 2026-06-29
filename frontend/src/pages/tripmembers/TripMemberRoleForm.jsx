import { useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'

import AccessDenied from '../../components/common/AccessDenied'
import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import useTripAccess from '../../hooks/useTripAccess'
import { getErrorMessage } from '../../services/api'
import tripMemberService from '../../services/tripMemberService'

const TripMemberRoleForm = () => {
  const { tripId, memberId } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const member = location.state?.member
  const { allowed, error: accessError, loading } = useTripAccess(tripId, 'manageMembers')
  const [role, setRole] = useState(member?.role === 'editor' ? 'editor' : 'viewer')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (event) => {
    event.preventDefault()
    if (member?.role === 'owner') {
      setError('Role owner tidak bisa diubah.')
      return
    }
    setSubmitting(true)
    setError('')
    try {
      await tripMemberService.updateRole(tripId, memberId, { role })
      navigate(`/trips/${tripId}`)
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal mengubah role member'))
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) return <Loading label="Memeriksa akses member..." />
  if (!allowed) return <AccessDenied backTo={`/trips/${tripId}`} message={accessError || 'Hanya owner yang bisa mengubah role member.'} />

  return (
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">Member role</p>
          <h2>Ubah akses member</h2>
        </div>
      </div>
      <form className="form-card" onSubmit={handleSubmit}>
        <ErrorBanner message={error} />
        {member?.role === 'owner' && <div className="field-note">Owner selalu punya akses penuh dan role tidak bisa diubah.</div>}
        <label>
          Role
          <select value={role} onChange={(event) => setRole(event.target.value)} disabled={member?.role === 'owner'}>
            <option value="viewer">Viewer</option>
            <option value="editor">Editor</option>
          </select>
        </label>
        <div className="role-help">
          <span>Viewer bisa melihat itinerary tanpa mengubah data.</span>
          <span>Editor bisa mengubah trip, hari, dan aktivitas.</span>
        </div>
        <button className="button-primary" disabled={submitting || member?.role === 'owner'} type="submit">
          {submitting ? 'Menyimpan...' : 'Simpan role'}
        </button>
      </form>
    </section>
  )
}

export default TripMemberRoleForm

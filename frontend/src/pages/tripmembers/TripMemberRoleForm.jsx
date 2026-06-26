import { useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import { getErrorMessage } from '../../services/api'
import tripMemberService from '../../services/tripMemberService'

const TripMemberRoleForm = () => {
  const { tripId, memberId } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const member = location.state?.member
  const [role, setRole] = useState(member?.role === 'editor' ? 'editor' : 'viewer')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (event) => {
    event.preventDefault()
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
        <label>
          Role
          <select value={role} onChange={(event) => setRole(event.target.value)}>
            <option value="viewer">Viewer</option>
            <option value="editor">Editor</option>
          </select>
        </label>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Menyimpan...' : 'Simpan role'}
        </button>
      </form>
    </section>
  )
}

export default TripMemberRoleForm

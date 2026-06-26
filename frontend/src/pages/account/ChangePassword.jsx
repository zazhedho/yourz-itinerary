import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import accountService from '../../services/accountService'
import { getErrorMessage } from '../../services/api'

const ChangePassword = () => {
  const navigate = useNavigate()
  const [form, setForm] = useState({ old_password: '', new_password: '' })
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      await accountService.changePassword(form)
      navigate('/account')
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal mengubah password'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">Keamanan</p>
          <h2>Ubah password</h2>
        </div>
      </div>
      <form className="form-card" onSubmit={handleSubmit}>
        <ErrorBanner message={error} />
        <label>
          Password lama
          <input name="old_password" type="password" value={form.old_password} onChange={handleChange} required />
        </label>
        <label>
          Password baru
          <input name="new_password" type="password" value={form.new_password} onChange={handleChange} required />
        </label>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Menyimpan...' : 'Simpan password baru'}
        </button>
      </form>
    </section>
  )
}

export default ChangePassword

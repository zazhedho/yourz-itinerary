import { useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import { getErrorMessage } from '../../services/api'
import authService from '../../services/authService'

const ResetPassword = () => {
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const [newPassword, setNewPassword] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const token = params.get('token') || ''

  const handleSubmit = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      await authService.resetPassword({ token, new_password: newPassword })
      navigate('/login', { replace: true })
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal reset password'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="auth-screen">
      <section className="auth-card">
        <h1>Reset password</h1>
        <ErrorBanner message={error} />
        <form className="form-inner" onSubmit={handleSubmit}>
          <label>
            Password baru
            <input type="password" value={newPassword} onChange={(event) => setNewPassword(event.target.value)} required />
          </label>
          <button className="button-primary" disabled={submitting || !token} type="submit">
            {submitting ? 'Menyimpan...' : 'Simpan password'}
          </button>
        </form>
        {!token && <ErrorBanner message="Token reset tidak ditemukan di URL." />}
        <p className="auth-link">
          <Link to="/login">Kembali masuk</Link>
        </p>
      </section>
    </main>
  )
}

export default ResetPassword

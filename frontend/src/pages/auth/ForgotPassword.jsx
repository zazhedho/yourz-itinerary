import { useState } from 'react'
import { Link } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import { getErrorMessage } from '../../services/api'
import authService from '../../services/authService'

const ForgotPassword = () => {
  const [email, setEmail] = useState('')
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    setMessage('')
    try {
      await authService.forgotPassword({ email })
      setMessage('Instruksi reset password dikirim jika email terdaftar.')
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal mengirim reset password'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="auth-screen">
      <section className="auth-card">
        <h1>Lupa password</h1>
        <ErrorBanner message={error} />
        {message && <div className="success-banner">{message}</div>}
        <form className="form-inner" onSubmit={handleSubmit}>
          <label>
            Email
            <input type="email" value={email} onChange={(event) => setEmail(event.target.value)} required />
          </label>
          <button className="button-primary" disabled={submitting} type="submit">
            {submitting ? 'Mengirim...' : 'Kirim instruksi'}
          </button>
        </form>
        <p className="auth-link">
          <Link to="/login">Kembali masuk</Link>
        </p>
      </section>
    </main>
  )
}

export default ForgotPassword

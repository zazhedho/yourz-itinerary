import { useState } from 'react'
import { Link } from 'react-router-dom'
import { MapPin, Plane } from 'lucide-react'

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
    <main className="auth-screen-split auth-screen-recovery auth-screen-forgot">
      <section className="auth-hero">
        <div className="auth-hero-header">
          <Link to="/login" className="auth-brand-badge">
            <MapPin size={20} />
            <span>Yourz Itinerary</span>
          </Link>
        </div>

        <div className="auth-hero-body">
          <p className="auth-kicker">Pemulihan akun</p>
          <h1>Kembali ke rencana perjalanan.</h1>
          <p>Masukkan email akunmu untuk menerima instruksi reset password.</p>
        </div>

        <div className="auth-route" aria-hidden="true">
          <span className="auth-route-track" />
          <span className="auth-route-point start">
            <MapPin size={14} />
          </span>
          <span className="auth-route-point middle">
            <MapPin size={14} />
          </span>
          <span className="auth-route-point end">
            <MapPin size={14} />
          </span>
          <span className="auth-route-traveler">
            <Plane size={18} />
          </span>
        </div>
      </section>

      <div className="auth-form-wrapper">
        <form className="auth-card" onSubmit={handleSubmit}>
          <div className="auth-card-header">
            <div>
              <p className="auth-kicker">Reset password</p>
              <h2>Pulihkan akun</h2>
            </div>
          </div>

          <ErrorBanner message={error} />
          {message && <div className="success-banner">{message}</div>}

          <div className="auth-fields">
            <label>
              Email
              <input
                autoComplete="email"
                maxLength={254}
                type="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                placeholder="nama@email.com"
                required
              />
            </label>
          </div>

          <button className="button-primary button-auth-submit" disabled={submitting} type="submit">
            {submitting ? 'Mengirim...' : 'Kirim instruksi'}
          </button>

          <div className="auth-footer-prompt">
            <Link to="/login" className="auth-footer-link">
              Kembali ke halaman masuk
            </Link>
          </div>
        </form>
      </div>
    </main>
  )
}

export default ForgotPassword

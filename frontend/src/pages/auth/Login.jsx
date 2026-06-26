import { MapPin } from 'lucide-react'
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import GoogleIdentityButton from '../../components/common/GoogleIdentityButton'
import { useAuth } from '../../hooks/useAuth'
import useRegisterStatus from '../../hooks/useRegisterStatus'
import { getGoogleClientId } from '../../utils/runtimeConfig'

const Login = () => {
  const { googleLogin, login, error } = useAuth()
  const { enabled: registerEnabled } = useRegisterStatus()
  const navigate = useNavigate()
  const [form, setForm] = useState({ identifier: '', password: '' })
  const [googleError, setGoogleError] = useState('')
  const [googleSubmitting, setGoogleSubmitting] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const googleClientId = getGoogleClientId()

  const handleSubmit = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    const ok = await login(form)
    setSubmitting(false)
    if (ok) navigate('/trips', { replace: true })
  }

  const handleGoogleCredential = async (idToken) => {
    setGoogleError('')
    setGoogleSubmitting(true)
    const ok = await googleLogin(idToken)
    setGoogleSubmitting(false)
    if (ok) navigate('/trips', { replace: true })
  }

  return (
    <main className="auth-screen">
      <section className="auth-hero">
        <div className="brand-mark">
          <MapPin size={24} />
        </div>
        <h1>Rencanakan trip bersama.</h1>
        <p>Masuk untuk melihat itinerary, anggota, aktivitas, dan pin lokasi dalam satu tempat.</p>
      </section>

      <form className="auth-card" onSubmit={handleSubmit}>
        <ErrorBanner message={error || googleError} />
        {googleClientId && (
          <GoogleIdentityButton
            disabled={submitting || googleSubmitting}
            label="Lanjutkan dengan Google"
            onCredential={handleGoogleCredential}
            onError={setGoogleError}
            text="signin_with"
          />
        )}
        <label>
          Email atau username
          <input
            autoComplete="username"
            value={form.identifier}
            onChange={(event) => setForm((current) => ({ ...current, identifier: event.target.value }))}
            required
          />
        </label>
        <label>
          Password
          <input
            autoComplete="current-password"
            type="password"
            value={form.password}
            onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))}
            required
          />
        </label>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Masuk...' : 'Masuk'}
        </button>
        {registerEnabled !== false && (
          <p className="auth-link">
            Belum punya akun? <Link to="/register">Daftar</Link>
          </p>
        )}
        <p className="auth-link">
          <Link to="/forgot-password">Lupa password?</Link>
        </p>
      </form>
    </main>
  )
}

export default Login

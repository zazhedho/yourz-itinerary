import { Eye, EyeOff, MapPin, Plane } from 'lucide-react'
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import GoogleIdentityButton from '../../components/common/GoogleIdentityButton'
import { useAuth } from '../../hooks/useAuth'
import useRegisterStatus from '../../hooks/useRegisterStatus'
import { sanitizeNoWhitespaceInput } from '../../utils/inputSanitizer'
import { getGoogleClientId } from '../../utils/runtimeConfig'

const Login = () => {
  const { googleLogin, login, error } = useAuth()
  const { enabled: registerEnabled } = useRegisterStatus()
  const navigate = useNavigate()
  const [form, setForm] = useState({ identifier: '', password: '' })
  const [showPassword, setShowPassword] = useState(false)
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

  const handleChange = (event) => {
    const { name, value } = event.target
    setForm((current) => ({
      ...current,
      [name]: name === 'identifier' ? sanitizeNoWhitespaceInput(value) : value,
    }))
  }

  return (
    <main className="auth-screen-split auth-screen-login">
      <section className="auth-hero">
        <div className="auth-hero-header">
          <Link to="/login" className="auth-brand-badge">
            <MapPin size={20} />
            <span>Yourz Itinerary</span>
          </Link>
        </div>

        <div className="auth-hero-body">
          <p className="auth-kicker">Masuk</p>
          <h1>Rencanakan perjalanan bersama.</h1>
          <p>Kelola itinerary, anggota, aktivitas, dan lokasi dalam satu tempat.</p>
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
              <p className="auth-kicker">Masuk</p>
              <h2>Selamat datang kembali</h2>
            </div>
          </div>

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

          <div className="auth-fields">
            <label className="input-field-group">
              <span>Email atau username</span>
              <input
                autoComplete="username"
                maxLength={100}
                name="identifier"
                placeholder="email@domain.com"
                value={form.identifier}
                onChange={handleChange}
                required
              />
            </label>

            <label className="input-field-group">
              <div className="label-with-action">
                <span>Password</span>
                <Link to="/forgot-password" className="forgot-password-link">
                  Lupa password?
                </Link>
              </div>
              <div className="password-field">
                <input
                  autoComplete="current-password"
                  maxLength={64}
                  name="password"
                  placeholder="Masukkan password"
                  type={showPassword ? 'text' : 'password'}
                  value={form.password}
                  onChange={handleChange}
                  required
                />
                <button
                  type="button"
                  className="password-toggle-btn"
                  onClick={() => setShowPassword(!showPassword)}
                  title={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
                  tabIndex={-1}
                >
                  {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </label>
          </div>

          <button className="button-primary button-auth-submit" disabled={submitting} type="submit">
            {submitting ? 'Memproses...' : 'Masuk'}
          </button>

          {registerEnabled !== false && (
            <div className="auth-footer-prompt">
              <span>Belum punya akun?</span>{' '}
              <Link to="/register" className="auth-footer-link">
                Daftar sekarang
              </Link>
            </div>
          )}
        </form>
      </div>
    </main>
  )
}

export default Login

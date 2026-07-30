import { useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { Eye, EyeOff, MapPin, Plane } from 'lucide-react'

import ErrorBanner from '../../components/common/ErrorBanner'
import { getErrorMessage } from '../../services/api'
import authService from '../../services/authService'
import {
  isPasswordValid,
  passwordRequirements,
  passwordStrength,
  passwordStrengthLabel,
  validatePassword,
} from '../../utils/passwordValidation'

const ResetPassword = () => {
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const [form, setForm] = useState({ newPassword: '', confirmPassword: '' })
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const token = params.get('token') || ''
  const validation = validatePassword(form.newPassword)
  const strength = passwordStrength(validation)
  const passwordMatches = form.confirmPassword && form.newPassword === form.confirmPassword

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()

    if (!isPasswordValid(validation)) {
      setError('Password baru belum memenuhi semua syarat.')
      return
    }
    if (!passwordMatches) {
      setError('Konfirmasi password tidak sama.')
      return
    }

    setSubmitting(true)
    setError('')
    try {
      await authService.resetPassword({ token, new_password: form.newPassword })
      navigate('/login', { replace: true })
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal reset password'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="auth-screen-split auth-screen-recovery auth-screen-reset">
      <section className="auth-hero">
        <div className="auth-hero-header">
          <Link to="/login" className="auth-brand-badge">
            <MapPin size={20} />
            <span>Yourz Itinerary</span>
          </Link>
        </div>

        <div className="auth-hero-body">
          <p className="auth-kicker">Password baru</p>
          <h1>Lanjutkan perjalananmu.</h1>
          <p>Buat password baru agar akunmu kembali aman digunakan.</p>
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
              <h2>Buat password baru</h2>
            </div>
          </div>

          <ErrorBanner message={error} />
          {!token && <ErrorBanner message="Token reset tidak ditemukan di URL." />}

          <div className="auth-fields">
            <label>
              Password baru
              <div className="password-field">
                <input
                  autoComplete="new-password"
                  disabled={!token}
                  maxLength={64}
                  name="newPassword"
                  placeholder="Minimal 8 karakter"
                  type={showPassword ? 'text' : 'password'}
                  value={form.newPassword}
                  onChange={handleChange}
                  required
                />
                <button
                  type="button"
                  className="password-toggle-btn"
                  disabled={!token}
                  onClick={() => setShowPassword((current) => !current)}
                  title={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
                >
                  {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </label>

            {form.newPassword && (
              <div className="password-validation-card">
                <div className="password-meter-row">
                  <div className="password-meter">
                    <span style={{ width: `${(strength / 5) * 100}%` }} />
                  </div>
                  <strong>{passwordStrengthLabel(strength)}</strong>
                </div>
                <div className="password-requirements">
                  {passwordRequirements.map(([key, label]) => (
                    <span className={validation[key] ? 'valid' : ''} key={key}>
                      {label}
                    </span>
                  ))}
                </div>
              </div>
            )}

            <label>
              Konfirmasi password
              <div className="password-field">
                <input
                  autoComplete="new-password"
                  disabled={!token}
                  maxLength={64}
                  name="confirmPassword"
                  placeholder="Ulangi password baru"
                  type={showConfirmPassword ? 'text' : 'password'}
                  value={form.confirmPassword}
                  onChange={handleChange}
                  required
                />
                <button
                  type="button"
                  className="password-toggle-btn"
                  disabled={!token}
                  onClick={() => setShowConfirmPassword((current) => !current)}
                  title={showConfirmPassword ? 'Sembunyikan password' : 'Tampilkan password'}
                >
                  {showConfirmPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </label>

            {form.confirmPassword && (
              <div className={`password-match-note ${passwordMatches ? 'valid' : ''}`}>
                {passwordMatches ? 'Password sama' : 'Password belum sama'}
              </div>
            )}
          </div>

          <button className="button-primary button-auth-submit" disabled={submitting || !token} type="submit">
            {submitting ? 'Menyimpan...' : 'Simpan password'}
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

export default ResetPassword

import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import useRegisterStatus from '../../hooks/useRegisterStatus'
import authService from '../../services/authService'
import { getErrorMessage } from '../../services/api'
import {
  isPasswordValid,
  passwordRequirements,
  passwordStrength,
  passwordStrengthLabel,
  validatePassword,
} from '../../utils/passwordValidation'

const Register = () => {
  const navigate = useNavigate()
  const { enabled, otp_enabled: otpEnabled, loading, error: statusError } = useRegisterStatus()
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({ name: '', email: '', phone: '', password: '', confirm_password: '', otp_code: '' })
  const [otpStep, setOtpStep] = useState(false)
  const validation = validatePassword(form.password)
  const strength = passwordStrength(validation)
  const passwordMatches = form.confirm_password && form.password === form.confirm_password

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  const validateRegistrationDetails = () => {
    if (!isPasswordValid(validation)) {
      setError('Password belum memenuhi semua syarat')
      return false
    }
    if (form.password !== form.confirm_password) {
      setError('Konfirmasi password tidak sama')
      return false
    }
    return true
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setError('')
    if (!otpStep && !validateRegistrationDetails()) return

    setSubmitting(true)
    try {
      if (otpEnabled && !otpStep) {
        await authService.sendRegisterOTP({ email: form.email, phone: form.phone })
        setOtpStep(true)
        return
      }

      if (otpEnabled && !form.otp_code.trim()) {
        setError('Kode OTP wajib diisi')
        return
      }

      const payload = { ...form }
      delete payload.confirm_password
      if (!otpEnabled || !payload.otp_code.trim()) delete payload.otp_code
      await authService.register(payload)
      navigate('/login', { replace: true })
    } catch (err) {
      setError(getErrorMessage(err, otpEnabled && !otpStep ? 'Gagal mengirim OTP' : 'Register failed'))
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) {
    return (
      <main className="auth-screen">
        <section className="auth-card">
          <Loading label="Mengecek status registrasi..." />
        </section>
      </main>
    )
  }

  if (!enabled) {
    return (
      <main className="auth-screen">
        <section className="auth-hero compact">
          <h1>Registrasi sedang ditutup.</h1>
          <p>Akun baru belum bisa dibuat saat ini. Silakan masuk jika sudah punya akun.</p>
        </section>
        <section className="auth-card">
          <ErrorBanner message={statusError || 'Public registration is currently disabled.'} />
          <Link className="button-primary" to="/login">
            Masuk
          </Link>
        </section>
      </main>
    )
  }

  return (
    <main className="auth-screen">
      <section className="auth-hero compact">
        <h1>{otpStep ? 'Verifikasi email.' : 'Buat akun Yourz.'}</h1>
        <p>
          {otpStep
            ? `Masukkan kode yang dikirim ke ${form.email}.`
            : 'Akun member bisa membuat trip dan mengundang pasangan atau teman lewat email.'}
        </p>
      </section>
      <form className="auth-card" onSubmit={handleSubmit}>
        <ErrorBanner message={error} />
        {!otpStep ? (
          <>
            <label>
              Nama
              <input name="name" value={form.name} onChange={handleChange} required />
            </label>
            <label>
              Email
              <input name="email" type="email" value={form.email} onChange={handleChange} required />
            </label>
            <label>
              Nomor HP
              <input name="phone" value={form.phone} onChange={handleChange} required />
            </label>
            <label>
              Password
              <input name="password" type="password" value={form.password} onChange={handleChange} required />
            </label>
            {form.password && (
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
              <input
                name="confirm_password"
                type="password"
                value={form.confirm_password}
                onChange={handleChange}
                required
              />
            </label>
            {form.confirm_password && (
              <div className={`password-match-note ${passwordMatches ? 'valid' : ''}`}>
                {passwordMatches ? 'Password sama' : 'Password belum sama'}
              </div>
            )}
          </>
        ) : (
          <>
            <label>
              Kode OTP
              <input
                autoComplete="one-time-code"
                inputMode="numeric"
                maxLength={6}
                name="otp_code"
                value={form.otp_code}
                onChange={handleChange}
                required
              />
            </label>
            <button className="button-secondary" onClick={() => setOtpStep(false)} type="button">
              Ubah email
            </button>
          </>
        )}
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? (otpStep ? 'Memverifikasi...' : 'Memproses...') : otpStep ? 'Verifikasi dan daftar' : 'Daftar'}
        </button>
        <p className="auth-link">
          Sudah punya akun? <Link to="/login">Masuk</Link>
        </p>
      </form>
    </main>
  )
}

export default Register

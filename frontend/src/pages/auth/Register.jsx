import { ArrowLeft, MailCheck, MapPin, RefreshCw, ShieldCheck } from 'lucide-react'
import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import GoogleIdentityButton from '../../components/common/GoogleIdentityButton'
import Loading from '../../components/common/Loading'
import { useAuth } from '../../hooks/useAuth'
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
import { getGoogleClientId } from '../../utils/runtimeConfig'

const isOTPRequiredError = (err) => getErrorMessage(err, '').toLowerCase().includes('otp_code is required')

const getStoredOTPCooldown = () => {
  const expireTime = sessionStorage.getItem('register_otp_cooldown')
  if (!expireTime) return 0

  const remaining = Math.floor((parseInt(expireTime, 10) - Date.now()) / 1000)
  if (remaining > 0) return remaining

  sessionStorage.removeItem('register_otp_cooldown')
  return 0
}

const storeOTPCooldown = (cooldownTime) => {
  sessionStorage.setItem('register_otp_cooldown', Date.now() + cooldownTime * 1000)
}

const Register = () => {
  const navigate = useNavigate()
  const { error: authError, googleLogin } = useAuth()
  const { enabled, otp_enabled: otpEnabled, loading, error: statusError } = useRegisterStatus()
  const [error, setError] = useState('')
  const [googleError, setGoogleError] = useState('')
  const [googleSubmitting, setGoogleSubmitting] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState(() => {
    const saved = sessionStorage.getItem('register_form')
    return saved ? JSON.parse(saved) : { name: '', email: '', phone: '', password: '', confirm_password: '', otp_code: '' }
  })
  const [otpStep, setOtpStep] = useState(() => {
    return sessionStorage.getItem('register_otp_step') === 'true'
  })
  const [cooldown, setCooldown] = useState(getStoredOTPCooldown)

  // Sync state to sessionStorage
  useEffect(() => {
    sessionStorage.setItem('register_form', JSON.stringify(form))
  }, [form])

  useEffect(() => {
    sessionStorage.setItem('register_otp_step', otpStep)
  }, [otpStep])

  useEffect(() => {
    let timer
    if (cooldown > 0) {
      timer = setInterval(() => {
        setCooldown((c) => c - 1)
      }, 1000)
    }
    return () => clearInterval(timer)
  }, [cooldown])

  const handleResendOTP = async () => {
    setError('')
    try {
      await requestRegisterOTP()
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal mengirim ulang OTP'))
    }
  }

  const requestRegisterOTP = async () => {
    const res = await authService.sendRegisterOTP({ email: form.email, phone: form.phone })
    const cooldownTime = res.data?.data?.cooldown || 60
    setCooldown(cooldownTime)
    storeOTPCooldown(cooldownTime)
    setOtpStep(true)
  }

  const handleBackToDetails = () => {
    setError('')
    setOtpStep(false)
    setCooldown(0)
    setForm((current) => ({ ...current, otp_code: '' }))
    sessionStorage.removeItem('register_otp_step')
    sessionStorage.removeItem('register_otp_cooldown')
  }
  const validation = validatePassword(form.password)
  const strength = passwordStrength(validation)
  const passwordMatches = form.confirm_password && form.password === form.confirm_password
  const googleClientId = getGoogleClientId()

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  const handleOTPChange = (event) => {
    const otpCode = event.target.value.replace(/\D/g, '').slice(0, 6)
    setForm((current) => ({ ...current, otp_code: otpCode }))
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
        await requestRegisterOTP()
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
      sessionStorage.removeItem('register_form')
      sessionStorage.removeItem('register_otp_step')
      sessionStorage.removeItem('register_otp_cooldown')
      navigate('/login', { replace: true })
    } catch (err) {
      if (!otpStep && isOTPRequiredError(err)) {
        try {
          await requestRegisterOTP()
          return
        } catch (otpErr) {
          setError(getErrorMessage(otpErr, 'Gagal mengirim OTP'))
          return
        }
      }
      setError(getErrorMessage(err, otpEnabled && !otpStep ? 'Gagal mengirim OTP' : 'Register failed'))
    } finally {
      setSubmitting(false)
    }
  }

  const handleGoogleCredential = async (idToken) => {
    setError('')
    setGoogleError('')
    setGoogleSubmitting(true)
    const ok = await googleLogin(idToken)
    setGoogleSubmitting(false)
    if (ok) navigate('/trips', { replace: true })
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
    <main className="auth-screen auth-screen-register">
      <section className="auth-hero">
        <div className="brand-mark">
          <MapPin size={24} />
        </div>
        <p className="auth-kicker">Yourz Itinerary</p>
        <h1>{otpStep ? 'Verifikasi email.' : 'Buat akun Yourz.'}</h1>
        <p>
          {otpStep
            ? `Masukkan kode yang dikirim ke ${form.email}.`
            : 'Akun member bisa membuat trip dan mengundang pasangan atau teman lewat email.'}
        </p>
        <div className="auth-hero-pills" aria-hidden="true">
          <span>Collaborative</span>
          <span>Mobile</span>
          <span>Private</span>
        </div>
      </section>
      <form className="auth-card" onSubmit={handleSubmit}>
        <div className="auth-card-header">
          <div className="auth-card-icon">
            <ShieldCheck size={20} />
          </div>
          <div>
            <p className="auth-kicker">{otpStep ? 'Verifikasi' : 'Daftar'}</p>
            <h2>{otpStep ? 'Cek email kamu' : 'Mulai dengan akun baru'}</h2>
          </div>
        </div>
        <ErrorBanner message={error || authError || googleError} />
        {!otpStep ? (
          <>
            {googleClientId && (
              <GoogleIdentityButton
                disabled={submitting || googleSubmitting}
                label="Daftar dengan Google"
                onCredential={handleGoogleCredential}
                onError={setGoogleError}
                text="signup_with"
              />
            )}
            <div className="auth-fields">
              <label>
                Nama
                <input name="name" placeholder="Nama lengkap" value={form.name} onChange={handleChange} required />
              </label>
              <label>
                Email
                <input
                  autoComplete="email"
                  name="email"
                  placeholder="email@domain.com"
                  type="email"
                  value={form.email}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                Nomor HP
                <input
                  autoComplete="tel"
                  inputMode="tel"
                  name="phone"
                  placeholder="628123456789"
                  value={form.phone}
                  onChange={handleChange}
                  required
                />
              </label>
              <label>
                Password
                <input
                  autoComplete="new-password"
                  name="password"
                  placeholder="Minimal 8 karakter"
                  type="password"
                  value={form.password}
                  onChange={handleChange}
                  required
                />
              </label>
            </div>
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
                autoComplete="new-password"
                name="confirm_password"
                placeholder="Ulangi password"
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
          <div className="otp-step-container">
            <div className="otp-spotlight">
              <div className="otp-icon-shell">
                <MailCheck size={28} />
              </div>
              <div className="otp-copy">
                <p className="auth-kicker">Kode verifikasi</p>
                <h3>Kami sudah mengirim OTP</h3>
                <p>Masukkan 6 digit kode yang dikirim ke email berikut.</p>
                <span>{form.email}</span>
              </div>
            </div>
            <label className="otp-code-field">
              <span>Kode OTP</span>
              <input
                autoComplete="one-time-code"
                inputMode="numeric"
                maxLength={6}
                name="otp_code"
                value={form.otp_code}
                onChange={handleOTPChange}
                placeholder="000000"
                required
              />
            </label>
            <div className="otp-actions">
              <button
                type="button"
                className={`button-secondary otp-resend-button ${cooldown > 0 ? 'cooldown' : ''}`}
                onClick={handleResendOTP}
                disabled={cooldown > 0}
              >
                <RefreshCw size={16} />
                {cooldown > 0 ? `Kirim ulang ${cooldown}s` : 'Kirim ulang'}
              </button>
              <button className="button-text" onClick={handleBackToDetails} type="button">
                <ArrowLeft size={16} />
                Ubah email
              </button>
            </div>
          </div>
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

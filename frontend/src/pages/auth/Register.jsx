import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import useRegisterStatus from '../../hooks/useRegisterStatus'
import authService from '../../services/authService'
import { getErrorMessage } from '../../services/api'

const Register = () => {
  const navigate = useNavigate()
  const { enabled, loading, error: statusError } = useRegisterStatus()
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({ name: '', email: '', phone: '', password: '', otp_code: '' })
  const [otpSent, setOtpSent] = useState(false)
  const [sendingOTP, setSendingOTP] = useState(false)

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  const handleSendOTP = async () => {
    if (!form.email) return
    setSendingOTP(true)
    setError('')
    try {
      await authService.sendRegisterOTP({ email: form.email })
      setOtpSent(true)
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal mengirim OTP'))
    } finally {
      setSendingOTP(false)
    }
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      await authService.register(form)
      navigate('/login', { replace: true })
    } catch (err) {
      setError(getErrorMessage(err, 'Register failed'))
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
        <h1>Buat akun Yourz.</h1>
        <p>Akun member bisa membuat trip dan mengundang pasangan atau teman lewat email.</p>
      </section>
      <form className="auth-card" onSubmit={handleSubmit}>
        <ErrorBanner message={error} />
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
        <label>
          Kode OTP (opsional)
          <input name="otp_code" value={form.otp_code} onChange={handleChange} />
        </label>
        <button className="button-secondary" disabled={sendingOTP || !form.email} onClick={handleSendOTP} type="button">
          {sendingOTP ? 'Mengirim...' : otpSent ? 'Kirim ulang OTP' : 'Kirim OTP'}
        </button>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Mendaftar...' : 'Daftar'}
        </button>
        <p className="auth-link">
          Sudah punya akun? <Link to="/login">Masuk</Link>
        </p>
      </form>
    </main>
  )
}

export default Register

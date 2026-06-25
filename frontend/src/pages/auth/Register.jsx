import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import authService from '../../services/authService'
import { getErrorMessage } from '../../services/api'

const Register = () => {
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [form, setForm] = useState({ name: '', email: '', phone: '', password: '' })

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
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

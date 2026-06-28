import { useState } from 'react'
import { Link } from 'react-router-dom'
import { KeyRound, ShieldAlert } from 'lucide-react'

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
    <main className="auth-screen auth-screen-login">
      <section className="auth-hero">
        <div className="brand-mark">
          <ShieldAlert size={24} />
        </div>
        <p className="auth-kicker">Pemulihan Akun</p>
        <h1>Lupa Kata Sandi?</h1>
        <p>Jangan panik! Cukup masukkan alamat email Anda, dan kami akan segera mengirimkan tautan ajaib untuk mengatur ulang kata sandi Anda.</p>
      </section>

      <form className="auth-card" onSubmit={handleSubmit}>
        <div className="auth-card-header">
          <div className="auth-card-icon">
            <KeyRound size={20} />
          </div>
          <div>
            <p className="auth-kicker">Lupa Password</p>
            <h2>Atur ulang sandi</h2>
          </div>
        </div>

        <ErrorBanner message={error} />
        {message && <div className="success-banner">{message}</div>}
        
        <div className="auth-fields">
          <label>
            Alamat Email
            <input 
              type="email" 
              value={email} 
              onChange={(event) => setEmail(event.target.value)} 
              placeholder="nama@email.com"
              required 
            />
          </label>
        </div>

        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Mengirim tautan...' : 'Kirim instruksi'}
        </button>
        
        <div className="auth-meta-row">
          <p className="auth-link">
            Ingat password Anda? <Link to="/login">Kembali masuk</Link>
          </p>
        </div>
      </form>
    </main>
  )
}

export default ForgotPassword

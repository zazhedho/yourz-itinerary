import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import accountService from '../../services/accountService'
import { getErrorMessage } from '../../services/api'
import {
  isPasswordValid,
  passwordRequirements,
  passwordStrength,
  passwordStrengthLabel,
  validatePassword,
} from '../../utils/passwordValidation'

const ChangePassword = () => {
  const navigate = useNavigate()
  const [form, setForm] = useState({ old_password: '', new_password: '', confirm_password: '' })
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const validation = validatePassword(form.new_password)
  const strength = passwordStrength(validation)
  const passwordMatches = form.confirm_password && form.new_password === form.confirm_password

  const handleChange = (event) => {
    setForm((current) => ({ ...current, [event.target.name]: event.target.value }))
  }

  const handleSubmit = async (event) => {
    event.preventDefault()
    
    if (!isPasswordValid(validation)) {
      setError('Password baru belum memenuhi semua persyaratan keamanan.')
      return
    }
    
    if (!passwordMatches) {
      setError('Konfirmasi password tidak cocok dengan password baru.')
      return
    }

    setSubmitting(true)
    setError('')
    try {
      await accountService.changePassword({
        old_password: form.old_password,
        new_password: form.new_password
      })
      navigate('/account')
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal mengubah password'))
    } finally {
      setSubmitting(false)
    }
  }

  const isFormValid = form.old_password && isPasswordValid(validation) && passwordMatches

  return (
    <section className="screen-stack" style={{ paddingBottom: '40px' }}>
      <div className="section-header">
        <div>
          <p className="eyebrow">Keamanan</p>
          <h2>Ubah Password</h2>
        </div>
      </div>
      <ErrorBanner message={error} />
      
      <form className="settings-group" onSubmit={handleSubmit}>
        <h3>Perbarui Kata Sandi</h3>
        <div className="settings-list" style={{ padding: '20px 16px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
          
          <label style={{ display: 'block' }}>
            <span style={{ display: 'block', marginBottom: '6px', fontSize: '13px', fontWeight: '600', color: 'var(--color-ink)' }}>Password Lama</span>
            <input 
              name="old_password" 
              type="password" 
              placeholder="Masukkan password saat ini"
              value={form.old_password} 
              onChange={handleChange} 
              required 
            />
          </label>
          
          <div style={{ height: '1px', background: 'var(--color-hairline)', margin: '4px 0' }} />

          <div>
            <label style={{ display: 'block', marginBottom: '8px' }}>
              <span style={{ display: 'block', marginBottom: '6px', fontSize: '13px', fontWeight: '600', color: 'var(--color-ink)' }}>Password Baru</span>
              <input 
                name="new_password" 
                type="password" 
                placeholder="Buat password yang kuat"
                value={form.new_password} 
                onChange={handleChange} 
                required 
              />
            </label>
            
            {form.new_password && (
              <div className="password-validation-card" style={{ marginTop: '12px' }}>
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
          </div>
          
          <label style={{ display: 'block' }}>
            <span style={{ display: 'block', marginBottom: '6px', fontSize: '13px', fontWeight: '600', color: 'var(--color-ink)' }}>Konfirmasi Password Baru</span>
            <input 
              name="confirm_password" 
              type="password" 
              placeholder="Ulangi password baru"
              value={form.confirm_password} 
              onChange={handleChange} 
              required 
            />
          </label>
          
          {form.confirm_password && (
            <div className={`password-match-note ${passwordMatches ? 'valid' : ''}`} style={{ marginTop: '-4px', fontSize: '13px' }}>
              {passwordMatches ? 'Password sama' : 'Password belum sama'}
            </div>
          )}

          <button 
            className="button-primary" 
            disabled={submitting || !isFormValid} 
            type="submit" 
            style={{ marginTop: '12px' }}
          >
            {submitting ? 'Menyimpan...' : 'Simpan Password Baru'}
          </button>
        </div>
      </form>
    </section>
  )
}

export default ChangePassword

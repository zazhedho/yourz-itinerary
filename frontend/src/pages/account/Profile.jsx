import { KeyRound, LogOut, MonitorSmartphone, Trash2, UserRound } from 'lucide-react'
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import ConfirmDialog from '../../components/common/ConfirmDialog'
import ErrorBanner from '../../components/common/ErrorBanner'
import { useAuth } from '../../hooks/useAuth'
import { useConfirm } from '../../hooks/useConfirm'
import accountService from '../../services/accountService'
import { getErrorMessage } from '../../services/api'

const Profile = () => {
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const { confirm, dialogProps } = useConfirm()
  const [name, setName] = useState(user?.name || '')
  const [phone, setPhone] = useState(user?.phone || '')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleUpdate = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      await accountService.updateProfile({ name, phone })
      await logout()
      navigate('/login', { replace: true })
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal update profil'))
    } finally {
      setSubmitting(false)
    }
  }

  const handleDeleteAccount = async () => {
    const ok = await confirm({
      title: 'Hapus akun',
      message: 'Akun dan semua data trip akan dihapus permanen. Lanjutkan?',
      confirmLabel: 'Hapus akun',
      danger: true,
    })
    if (!ok) return
    try {
      await accountService.deleteAccount()
      await logout()
      navigate('/login', { replace: true })
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal menghapus akun'))
    }
  }

  const handleLogout = async () => {
    await logout()
    navigate('/login', { replace: true })
  }

  return (
    <section className="screen-stack" style={{ paddingBottom: '40px' }}>
      <div className="profile-hero">
        {user?.avatar_url ? (
          <img 
            src={user.avatar_url} 
            alt="Profile Avatar" 
            className="profile-avatar-large" 
            style={{ objectFit: 'cover' }}
          />
        ) : (
          <div className="profile-avatar-large">
            {user?.name?.charAt(0)?.toUpperCase() || 'U'}
          </div>
        )}
        <div className="profile-hero-text">
          <h2>{user?.name || 'Pengguna'}</h2>
          <p>{user?.email}</p>
        </div>
      </div>

      <ErrorBanner message={error} />

      <form className="settings-group" onSubmit={handleUpdate}>
        <h3>Informasi Pribadi</h3>
        <div className="settings-list" style={{ padding: '20px 16px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <label style={{ display: 'block' }}>
            <span style={{ display: 'block', marginBottom: '6px', fontSize: '13px', fontWeight: '600', color: 'var(--color-ink)' }}>Nama lengkap</span>
            <input value={name} onChange={(event) => setName(event.target.value)} required />
          </label>
          <label style={{ display: 'block' }}>
            <span style={{ display: 'block', marginBottom: '6px', fontSize: '13px', fontWeight: '600', color: 'var(--color-ink)' }}>Nomor HP</span>
            <input value={phone} onChange={(event) => setPhone(event.target.value)} placeholder="08..." />
          </label>
          <label style={{ display: 'block' }}>
            <span style={{ display: 'block', marginBottom: '6px', fontSize: '13px', fontWeight: '600', color: 'var(--color-ink)' }}>Alamat email</span>
            <input value={user?.email || ''} disabled />
          </label>
          <button className="button-primary" disabled={submitting || (name === user?.name && phone === (user?.phone || ''))} type="submit" style={{ marginTop: '8px' }}>
            {submitting ? 'Menyimpan...' : 'Simpan Perubahan'}
          </button>
        </div>
      </form>

      <div className="settings-group">
        <h3>Keamanan & Akses</h3>
        <div className="settings-list">
          <Link className="settings-list-item" to="/account/password">
            <div className="settings-list-icon">
              <KeyRound size={18} />
            </div>
            <div className="settings-list-text">
              <strong>Ubah Password</strong>
              <span>Perbarui kata sandi akun Anda</span>
            </div>
          </Link>
          <Link className="settings-list-item" to="/account/sessions">
            <div className="settings-list-icon">
              <MonitorSmartphone size={18} />
            </div>
            <div className="settings-list-text">
              <strong>Sesi Aktif</strong>
              <span>Kelola perangkat yang sedang login</span>
            </div>
          </Link>
          <button className="settings-list-item" onClick={handleLogout} type="button">
            <div className="settings-list-icon">
              <LogOut size={18} />
            </div>
            <div className="settings-list-text">
              <strong>Keluar</strong>
              <span>Akhiri sesi pada perangkat ini</span>
            </div>
          </button>
        </div>
      </div>

      <div className="settings-group danger-zone">
        <h3>Zona Berbahaya</h3>
        <div className="settings-list">
          <button className="settings-list-item danger-item" onClick={handleDeleteAccount} type="button">
            <div className="settings-list-icon">
              <Trash2 size={18} />
            </div>
            <div className="settings-list-text">
              <strong>Hapus Akun</strong>
              <span>Tindakan ini permanen dan menghapus semua data Anda</span>
            </div>
          </button>
        </div>
      </div>

      <ConfirmDialog {...dialogProps} />
    </section>
  )
}

export default Profile

import { LogOut, Trash2 } from 'lucide-react'
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
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleUpdate = async (event) => {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    try {
      await accountService.updateProfile({ name })
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
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">Akun</p>
          <h2>Profil kamu</h2>
        </div>
      </div>
      <ErrorBanner message={error} />
      <form className="form-card" onSubmit={handleUpdate}>
        <label>
          Nama
          <input value={name} onChange={(event) => setName(event.target.value)} required />
        </label>
        <label>
          Email
          <input value={user?.email || ''} disabled />
        </label>
        <button className="button-primary" disabled={submitting} type="submit">
          {submitting ? 'Menyimpan...' : 'Simpan profil'}
        </button>
      </form>
      <div className="form-card">
        <Link className="button-secondary" to="/account/password">
          Ubah password
        </Link>
        <Link className="button-secondary" to="/account/sessions">
          Kelola sesi aktif
        </Link>
        <button className="button-secondary" onClick={handleLogout} type="button">
          <LogOut size={17} />
          Keluar
        </button>
      </div>
      <button className="button-secondary danger-action" onClick={handleDeleteAccount} type="button">
        <Trash2 size={17} />
        Hapus akun
      </button>
      <ConfirmDialog {...dialogProps} />
    </section>
  )
}

export default Profile

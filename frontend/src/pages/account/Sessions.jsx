import { MonitorSmartphone, Trash2, LogOut } from 'lucide-react'
import { useEffect, useState } from 'react'

import ConfirmDialog from '../../components/common/ConfirmDialog'
import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import { useConfirm } from '../../hooks/useConfirm'
import { getErrorMessage, getResponseData } from '../../services/api'
import sessionService from '../../services/sessionService'

const Sessions = () => {
  const { confirm, dialogProps } = useConfirm()
  const [sessions, setSessions] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const loadSessions = () => {
    setLoading(true)
    sessionService
      .getActiveSessions()
      .then((response) => setSessions(getResponseData(response)?.sessions || []))
      .catch((err) => setError(getErrorMessage(err, 'Gagal memuat sesi')))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    let active = true
    sessionService
      .getActiveSessions()
      .then((response) => {
        if (active) setSessions(getResponseData(response)?.sessions || [])
      })
      .catch((err) => {
        if (active) setError(getErrorMessage(err, 'Gagal memuat sesi'))
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => { active = false }
  }, [])

  const handleRevoke = async (session) => {
    const ok = await confirm({
      title: 'Cabut sesi',
      message: 'Hapus sesi ini dari daftar? Perangkat akan logout otomatis.',
      confirmLabel: 'Cabut',
      danger: true,
    })
    if (!ok) return
    try {
      await sessionService.revoke(session.session_id)
      await loadSessions()
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal mencabut sesi'))
    }
  }

  const handleRevokeOthers = async () => {
    const ok = await confirm({
      title: 'Cabut sesi lain',
      message: 'Semua sesi selain yang aktif saat ini akan dicabut. Lanjutkan?',
      confirmLabel: 'Cabut semua',
      danger: true,
    })
    if (!ok) return
    try {
      await sessionService.revokeOthers()
      await loadSessions()
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal mencabut sesi'))
    }
  }

  if (loading) return <Loading label="Memuat sesi..." />

  return (
    <section className="screen-stack" style={{ paddingBottom: '40px' }}>
      <div className="section-header">
        <div>
          <p className="eyebrow">Keamanan</p>
          <h2>Sesi Aktif</h2>
        </div>
      </div>
      <ErrorBanner message={error} />
      <div className="settings-group">
        <h3 style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', paddingRight: '16px' }}>
          <span>Perangkat yang Login</span>
          {sessions.length > 0 && (
            <span style={{ background: 'var(--color-surface)', color: 'var(--color-brand)', padding: '2px 8px', borderRadius: '12px', fontSize: '11px', fontWeight: '700', letterSpacing: 'normal' }}>
              {sessions.length} perangkat
            </span>
          )}
        </h3>
        {sessions.length ? (
          <div className="settings-list">
            {sessions.map((session) => (
              <div className="settings-list-item" key={session.session_id} style={{ cursor: 'default', background: session.is_current_session ? 'var(--color-surface)' : 'transparent' }}>
                <div className="settings-list-icon" style={{ background: session.is_current_session ? 'var(--color-brand)' : 'var(--color-soft)', color: session.is_current_session ? '#ffffff' : 'var(--color-ink)' }}>
                  <MonitorSmartphone size={18} />
                </div>
                <div className="settings-list-text">
                  <strong style={{ display: '-webkit-box', WebkitLineClamp: 1, WebkitBoxOrient: 'vertical', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                    {session.device_info || 'Unknown Device'}
                    {session.is_current_session && <span style={{ marginLeft: '8px', fontSize: '12px', fontWeight: 'normal', color: 'var(--color-brand)', background: 'rgba(52, 211, 153, 0.1)', padding: '2px 6px', borderRadius: '4px' }}>Saat Ini</span>}
                  </strong>
                  <span>IP: {session.ip || '-'}</span>
                </div>
                {!session.is_current_session && (
                  <button 
                    className="icon-link danger" 
                    onClick={() => handleRevoke(session)} 
                    type="button" 
                    title="Cabut sesi"
                    style={{ marginLeft: 'auto', flexShrink: 0 }}
                  >
                    <Trash2 size={16} />
                  </button>
                )}
              </div>
            ))}
            <button 
              className="settings-list-item" 
              onClick={handleRevokeOthers} 
              type="button" 
              disabled={sessions.length <= 1}
              style={{ opacity: sessions.length <= 1 ? 0.5 : 1, cursor: sessions.length <= 1 ? 'not-allowed' : 'pointer', borderTop: '1px solid var(--color-hairline)' }}
            >
              <div className="settings-list-icon" style={{ background: 'rgba(193, 53, 21, 0.1)', color: 'var(--color-primary)' }}>
                <LogOut size={18} />
              </div>
              <div className="settings-list-text">
                <strong style={{ color: 'var(--color-primary)' }}>Cabut Semua Sesi Lain</strong>
                <span>Keluarkan akun dari semua perangkat lain</span>
              </div>
            </button>
          </div>
        ) : (
          <div className="empty-card" style={{ marginTop: 0 }}>Tidak ada sesi aktif yang bisa ditampilkan.</div>
        )}
      </div>

      <ConfirmDialog {...dialogProps} />
    </section>
  )
}

export default Sessions

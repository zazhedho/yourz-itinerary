import { Trash2 } from 'lucide-react'
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
      .then((response) => setSessions(getResponseData(response) || []))
      .catch((err) => setError(getErrorMessage(err, 'Gagal memuat sesi')))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    let active = true
    sessionService
      .getActiveSessions()
      .then((response) => {
        if (active) setSessions(getResponseData(response) || [])
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
      await sessionService.revoke(session.id)
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
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">Sesi</p>
          <h2>Sesi aktif</h2>
        </div>
      </div>
      <ErrorBanner message={error} />
      {sessions.length ? (
        <div className="reorder-list">
          {sessions.map((session) => (
            <article className="reorder-row" key={session.id}>
              <div>
                <strong>{session.user_agent || 'Unknown device'}</strong>
                <span>{session.ip_address || '-'}</span>
              </div>
              <button className="icon-link danger" onClick={() => handleRevoke(session)} type="button" title="Cabut sesi">
                <Trash2 size={16} />
              </button>
            </article>
          ))}
        </div>
      ) : (
        <div className="empty-card">Tidak ada sesi aktif yang bisa ditampilkan.</div>
      )}
      <button className="button-secondary danger-action" onClick={handleRevokeOthers} type="button">
        Cabut semua sesi lain
      </button>
      <ConfirmDialog {...dialogProps} />
    </section>
  )
}

export default Sessions

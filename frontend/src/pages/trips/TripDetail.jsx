import { Edit3, LogOut, Plus, Trash2, UserMinus, UserPlus } from 'lucide-react'
import { useCallback, useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'

import ConfirmDialog from '../../components/common/ConfirmDialog'
import DayTimeline from '../../components/itinerary/DayTimeline'
import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import { useConfirm } from '../../hooks/useConfirm'
import itineraryDayService from '../../services/itineraryDayService'
import itineraryItemService from '../../services/itineraryItemService'
import tripMemberService from '../../services/tripMemberService'
import tripService from '../../services/tripService'
import { getErrorMessage, getResponseData } from '../../services/api'
import { formatDateRange, roleLabel } from '../../utils/formatters'

const shortId = (value = '') => value.slice(0, 8)
const memberDisplayName = (member) => member.user_name || `User ${shortId(member.user_id)}`
const memberDisplayMeta = (member) => member.user_email || member.user_id

const TripDetail = () => {
  const { tripId } = useParams()
  const navigate = useNavigate()
  const { confirm, dialogProps } = useConfirm()
  const [trip, setTrip] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const loadTrip = useCallback(() => {
    setLoading(true)
    tripService
      .getById(tripId)
      .then((response) => setTrip(getResponseData(response)))
      .catch((err) => setError(getErrorMessage(err, 'Gagal memuat detail trip')))
      .finally(() => setLoading(false))
  }, [tripId])

  useEffect(() => {
    let active = true
    tripService
      .getById(tripId)
      .then((response) => {
        if (active) setTrip(getResponseData(response))
      })
      .catch((err) => {
        if (active) setError(getErrorMessage(err, 'Gagal memuat detail trip'))
      })
      .finally(() => {
        if (active) setLoading(false)
      })

    return () => {
      active = false
    }
  }, [tripId])

  const runAction = async (confirmMessage, action, fallback) => {
    if (confirmMessage) {
      const ok = await confirm({
        title: 'Konfirmasi aksi',
        message: confirmMessage,
        confirmLabel: 'Lanjutkan',
        danger: true,
      })
      if (!ok) return
    }
    setError('')
    try {
      await action()
    } catch (err) {
      setError(getErrorMessage(err, fallback))
    }
  }

  const handleLeaveTrip = () =>
    runAction('Keluar dari itinerary ini?', async () => {
      await tripMemberService.leaveTrip(trip.id)
      navigate('/trips', { replace: true })
    }, 'Gagal keluar dari trip')

  const handleDeleteTrip = () =>
    runAction('Hapus trip ini?', async () => {
      await tripService.delete(trip.id)
      navigate('/trips', { replace: true })
    }, 'Gagal menghapus trip')

  const handleRemoveMember = (member) =>
    runAction('Hapus member dari trip?', async () => {
      await tripMemberService.remove(trip.id, member.id)
      await loadTrip()
    }, 'Gagal menghapus member')

  const handleDeleteDay = (day) =>
    runAction('Hapus hari itinerary ini?', async () => {
      await itineraryDayService.delete(day.id)
      await loadTrip()
    }, 'Gagal menghapus hari')

  const handleDeleteItem = (item) =>
    runAction('Hapus aktivitas ini?', async () => {
      await itineraryItemService.delete(item.id)
      await loadTrip()
    }, 'Gagal menghapus aktivitas')

  if (loading) return <Loading label="Memuat detail trip..." />
  if (!trip) return <ErrorBanner message={error || 'Trip tidak ditemukan'} />

  const members = trip.members || []
  const days = trip.days || []
  const itemCount = days.reduce((total, day) => total + (day.items?.length || 0), 0)

  return (
    <section className="screen-stack trip-detail-screen">
      <div className="detail-cover">
        <div>
          <p className="detail-kicker">{trip.destination || 'Shared itinerary'}</p>
          <h2>{trip.title}</h2>
          <div className="detail-meta-row">
            <span>{formatDateRange(trip.start_date, trip.end_date)}</span>
            <span>{trip.currency_code}</span>
          </div>
        </div>
      </div>

      <div className="trip-summary-grid">
        <div>
          <strong>{days.length}</strong>
          <span>Hari</span>
        </div>
        <div>
          <strong>{itemCount}</strong>
          <span>Aktivitas</span>
        </div>
        <div>
          <strong>{members.length}</strong>
          <span>Member</span>
        </div>
      </div>

      <ErrorBanner message={error} />

      <section className="action-panel" aria-label="Trip actions">
        <Link className="action-tile" to={`/trips/${trip.id}/edit`}>
          <Edit3 size={18} />
          <span>Edit trip</span>
        </Link>
        <Link className="action-tile" to={`/trips/${trip.id}/members/add`}>
          <UserPlus size={18} />
          <span>Tambah member</span>
        </Link>
        <Link className="action-tile primary" to={`/trips/${trip.id}/days/new`}>
          <Plus size={19} />
          <span>Tambah hari</span>
        </Link>
      </section>

      <section className="content-section">
        <div className="section-heading">
          <div>
            <p className="eyebrow">Members</p>
            <h2>Akses itinerary</h2>
          </div>
          <Link className="icon-link" to={`/trips/${trip.id}/members/add`} title="Tambah member">
            <UserPlus size={16} />
          </Link>
        </div>
        <div className="member-list">
          {members.map((member) => (
            <article className="member-row" key={member.id}>
              <div className="member-avatar">
                {member.avatar_url ? <img alt="" src={member.avatar_url} /> : memberDisplayName(member).charAt(0)}
              </div>
              <div>
                <strong>{memberDisplayName(member)}</strong>
                <span>{memberDisplayMeta(member)}</span>
                <small>{roleLabel(member.role)}</small>
              </div>
              {member.role !== 'owner' && (
                <div className="member-actions">
                  <Link state={{ member }} to={`/trips/${trip.id}/members/${member.id}/role`}>
                    Edit role
                  </Link>
                  <button onClick={() => handleRemoveMember(member)} type="button" title="Hapus member">
                    <UserMinus size={14} />
                  </button>
                </div>
              )}
            </article>
          ))}
        </div>
      </section>

      <section className="content-section">
        <div className="section-heading">
          <div>
            <p className="eyebrow">Timeline</p>
            <h2>Rencana perjalanan</h2>
          </div>
          <Link className="icon-link" to={`/trips/${trip.id}/days/new`} title="Tambah hari">
            <Plus size={18} />
          </Link>
        </div>
        <DayTimeline
          currency={trip.currency_code}
          days={days}
          onDeleteDay={handleDeleteDay}
          onDeleteItem={handleDeleteItem}
        />
      </section>

      <section className="danger-zone">
        <div>
          <p className="eyebrow">Danger zone</p>
          <h2>Aksi trip</h2>
        </div>
        <button aria-label="Keluar dari trip" onClick={handleLeaveTrip} type="button">
          <LogOut size={17} />
          Keluar
        </button>
        <button aria-label="Hapus trip" onClick={handleDeleteTrip} type="button">
          <Trash2 size={17} />
          Hapus
        </button>
      </section>
      <ConfirmDialog {...dialogProps} />
    </section>
  )
}

export default TripDetail

import { CalendarDays, Edit3, ListChecks, LogOut, Plus, Trash2, UserMinus, UserPlus, UsersRound, Calendar, Wallet, Settings, X } from 'lucide-react'
import { useCallback, useEffect, useState } from 'react'
import { Link, useNavigate, useParams, useLocation } from 'react-router-dom'

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
import { getDestinationPhoto } from '../../services/unsplashService'
import { formatDateRange, formatMoney, roleLabel } from '../../utils/formatters'

const shortId = (value = '') => value.slice(0, 8)
const memberDisplayName = (member) => member.user_name || `User ${shortId(member.user_id)}`
const memberDisplayMeta = (member) => member.user_email || member.user_id

const TripDetail = () => {
  const { tripId } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const index = location.state?.index || 0
  const { confirm, dialogProps } = useConfirm()
  const [trip, setTrip] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showMembers, setShowMembers] = useState(false)
  const [showSettings, setShowSettings] = useState(false)
  const [coverPhoto, setCoverPhoto] = useState(null)

  useEffect(() => {
    if (!trip?.destination) return
    let active = true
    getDestinationPhoto(trip.destination, index).then((url) => {
      if (active) setCoverPhoto(url)
    })
    return () => {
      active = false
    }
  }, [trip?.destination, index])

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
  const calculatedCostEstimate = days.reduce(
    (total, day) => total + (day.items || []).reduce((dayTotal, item) => dayTotal + Number(item.cost_estimate || 0), 0),
    0,
  )
  const totalCostEstimate = trip.total_cost_estimate ?? calculatedCostEstimate

  const getNextDate = () => {
    if (days.length === 0) return trip.start_date ? String(trip.start_date).split('T')[0] : ''
    const validDates = days
      .filter(d => d.date)
      .map(d => new Date(d.date).getTime())
      .filter(t => !isNaN(t))
      
    if (validDates.length === 0) return trip.start_date ? String(trip.start_date).split('T')[0] : ''
    
    const nextDate = new Date(Math.max(...validDates))
    nextDate.setDate(nextDate.getDate() + 1)
    return nextDate.toISOString().split('T')[0]
  }

  const getNextDayNumber = () => {
    if (days.length === 0) return 1
    const nums = days.map(d => Number(d.day_number)).filter(n => !isNaN(n))
    if (nums.length === 0) return 1
    return Math.max(...nums) + 1
  }

  return (
    <section className="screen-stack trip-detail-screen">
      <div 
        className="detail-cover"
        style={coverPhoto ? { backgroundImage: `linear-gradient(180deg, rgba(0, 0, 0, 0.05) 0%, rgba(0, 0, 0, 0.75) 100%), url("${coverPhoto}")` } : undefined}
      >
        <div className="cover-top-actions">
          <Link className="detail-cover-action" to={`/trips/${trip.id}/edit`} title="Edit trip" aria-label="Edit trip">
            <Edit3 size={17} />
          </Link>
          <button className="detail-cover-action" onClick={() => setShowSettings(true)} title="Pengaturan lanjutan" aria-label="Pengaturan lanjutan" type="button">
            <Settings size={17} />
          </button>
        </div>
        <div>
          <p className="detail-kicker">{trip.destination || 'Shared itinerary'}</p>
          <h2>{trip.title}</h2>
          <div className="detail-meta-row">
            <span>
              <Calendar size={14} /> {formatDateRange(trip.start_date, trip.end_date)}
            </span>
            <span>
              <Wallet size={14} /> {trip.currency_code}
            </span>
          </div>
        </div>
      </div>

      <div className="trip-summary-grid">
        <div>
          <span className="summary-icon">
            <CalendarDays size={16} />
          </span>
          <strong>{days.length}</strong>
          <span>Hari</span>
        </div>
        <div>
          <span className="summary-icon">
            <ListChecks size={16} />
          </span>
          <strong>{itemCount}</strong>
          <span>Aktivitas</span>
        </div>
        <div>
          <span className="summary-icon">
            <Wallet size={16} />
          </span>
          <strong className="summary-money">{formatMoney(totalCostEstimate, trip.currency_code)}</strong>
          <span>Budget</span>
        </div>
        <button
          aria-expanded={showMembers}
          aria-controls="trip-member-list"
          className="summary-button"
          onClick={() => setShowMembers((value) => !value)}
          type="button"
        >
          <span className="summary-icon">
            <UsersRound size={16} />
          </span>
          <strong>{members.length}</strong>
          <span>Member</span>
        </button>
      </div>

      <ErrorBanner message={error} />

      {showMembers && (
        <section className="content-section" id="trip-member-list">
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
      )}

      <section className="content-section">
        <div className="section-heading">
          <div>
            <p className="eyebrow">Timeline</p>
            <h2>Rencana perjalanan</h2>
          </div>
          <Link className="icon-link" state={{ nextDayNumber: getNextDayNumber(), nextDate: getNextDate() }} to={`/trips/${trip.id}/days/new`} title="Tambah hari">
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

      {showSettings && (
        <div className="dialog-backdrop" onClick={() => setShowSettings(false)}>
          <div className="dialog-card" onClick={(e) => e.stopPropagation()}>
            <div className="section-heading" style={{ marginBottom: 20 }}>
              <h2>Pengaturan Lanjutan</h2>
              <button className="modal-close" onClick={() => setShowSettings(false)} type="button" aria-label="Tutup">
                <X size={18} />
              </button>
            </div>
            <div className="danger-list">
              <div className="danger-row">
                <div className="danger-text">
                  <strong>Keluar dari Trip</strong>
                  <span>Anda akan kehilangan akses ke itinerary ini.</span>
                </div>
                <button aria-label="Keluar dari trip" onClick={() => { setShowSettings(false); handleLeaveTrip(); }} type="button">
                  <LogOut size={16} />
                  Keluar
                </button>
              </div>
              <div className="danger-row">
                <div className="danger-text">
                  <strong>Hapus Trip</strong>
                  <span>Tindakan ini permanen dan menghapus data.</span>
                </div>
                <button aria-label="Hapus trip" onClick={() => { setShowSettings(false); handleDeleteTrip(); }} type="button">
                  <Trash2 size={16} />
                  Hapus
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      <ConfirmDialog {...dialogProps} />
    </section>
  )
}

export default TripDetail

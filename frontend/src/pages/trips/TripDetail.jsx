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

  return (
    <section className="screen-stack">
      <div className="detail-cover">
        <div>
          <p>{trip.destination || 'Shared trip'}</p>
          <h2>{trip.title}</h2>
          <span>{formatDateRange(trip.start_date, trip.end_date)}</span>
        </div>
      </div>

      <div className="quick-actions">
        <Link to={`/trips/${trip.id}/edit`}>
          <Edit3 size={17} />
          Edit
        </Link>
        <Link to={`/trips/${trip.id}/members/add`}>
          <UserPlus size={17} />
          Member
        </Link>
        <Link to={`/trips/${trip.id}/days/new`}>
          <Plus size={17} />
          Day
        </Link>
        <button aria-label="Keluar dari trip" onClick={handleLeaveTrip} type="button">
          <LogOut size={17} />
          Leave
        </button>
        <button aria-label="Hapus trip" className="danger-action" onClick={handleDeleteTrip} type="button">
          <Trash2 size={17} />
          Delete
        </button>
      </div>

      <ErrorBanner message={error} />

      <section className="member-strip">
        {(trip.members || []).map((member) => (
          <div className="member-chip" key={member.id}>
            {member.role !== 'owner' ? (
              <Link state={{ member }} to={`/trips/${trip.id}/members/${member.id}/role`}>
                <span>{roleLabel(member.role)}</span>
              </Link>
            ) : (
              <span>{roleLabel(member.role)}</span>
            )}
            <small>{member.user_id.slice(0, 8)}</small>
            {member.role !== 'owner' && (
              <button onClick={() => handleRemoveMember(member)} type="button" title="Hapus member">
                <UserMinus size={13} />
              </button>
            )}
          </div>
        ))}
      </section>

      <DayTimeline
        currency={trip.currency_code}
        days={trip.days || []}
        onDeleteDay={handleDeleteDay}
        onDeleteItem={handleDeleteItem}
      />
      <ConfirmDialog {...dialogProps} />
    </section>
  )
}

export default TripDetail

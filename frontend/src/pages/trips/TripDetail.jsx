import { Edit3, Plus, UserPlus } from 'lucide-react'
import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'

import DayTimeline from '../../components/itinerary/DayTimeline'
import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import tripService from '../../services/tripService'
import { getErrorMessage, getResponseData } from '../../services/api'
import { formatDateRange, roleLabel } from '../../utils/formatters'

const TripDetail = () => {
  const { tripId } = useParams()
  const [trip, setTrip] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    tripService
      .getById(tripId)
      .then((response) => setTrip(getResponseData(response)))
      .catch((err) => setError(getErrorMessage(err, 'Gagal memuat detail trip')))
      .finally(() => setLoading(false))
  }, [tripId])

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
      </div>

      <ErrorBanner message={error} />

      <section className="member-strip">
        {(trip.members || []).map((member) => (
          <div className="member-chip" key={member.id}>
            <span>{roleLabel(member.role)}</span>
            <small>{member.user_id.slice(0, 8)}</small>
          </div>
        ))}
      </section>

      <DayTimeline currency={trip.currency_code} days={trip.days || []} />
    </section>
  )
}

export default TripDetail

import { CalendarDays, MapPin, UsersRound } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useEffect, useState } from 'react'

import { formatDateRange } from '../../utils/formatters'
import { getDestinationPhoto } from '../../services/unsplashService'

const TripCard = ({ trip, index = 0 }) => {
  const [coverPhoto, setCoverPhoto] = useState(null)

  useEffect(() => {
    let active = true
    getDestinationPhoto(trip.destination, index).then((url) => {
      if (active) setCoverPhoto(url)
    })
    return () => {
      active = false
    }
  }, [trip.destination, index])

  return (
    <Link className="trip-card" to={`/trips/${trip.id}`}>
      {coverPhoto ? (
        <img alt={trip.destination || trip.title} src={coverPhoto} />
      ) : (
        <div className="trip-card-placeholder" style={{ width: '100%', height: '140px', background: '#e2e8f0' }} />
      )}
      <div className="trip-card-body">
        <div>
          <h2>{trip.title}</h2>
          <p>
            <MapPin size={14} />
            {trip.destination || 'Destinasi belum diatur'}
          </p>
        </div>
        <div className="trip-meta-grid">
          <span>
            <CalendarDays size={14} />
            {formatDateRange(trip.start_date, trip.end_date)}
          </span>
          <span>
            <UsersRound size={14} />
            {trip.member_count || 1} member
          </span>
        </div>
      </div>
    </Link>
  )
}

export default TripCard

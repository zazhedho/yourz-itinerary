import { CalendarDays, MapPin, UsersRound } from 'lucide-react'
import { Link } from 'react-router-dom'

import { formatDateRange } from '../../utils/formatters'

const photos = [
  'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=900&q=80',
  'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=900&q=80',
  'https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?auto=format&fit=crop&w=900&q=80',
]

const TripCard = ({ trip, index = 0 }) => (
  <Link className="trip-card" to={`/trips/${trip.id}`}>
    <img alt={trip.title} src={photos[index % photos.length]} />
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

export default TripCard

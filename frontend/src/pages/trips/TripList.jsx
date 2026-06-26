import { Plus } from 'lucide-react'
import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import Loading from '../../components/common/Loading'
import TripCard from '../../components/trips/TripCard'
import tripService from '../../services/tripService'
import { getErrorMessage, getResponseData } from '../../services/api'

const TripList = () => {
  const [trips, setTrips] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    tripService
      .getAll()
      .then((response) => setTrips(getResponseData(response) || []))
      .catch((err) => setError(getErrorMessage(err, 'Gagal memuat trips')))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <Loading label="Memuat trips..." />

  return (
    <section className="screen-stack">
      <div className="section-header" style={{ marginBottom: '8px', marginTop: '12px' }}>
        <div>
          <p className="eyebrow" style={{ color: 'var(--color-brand)' }}>Rencana Mendatang</p>
          <h1 style={{ margin: 0, fontSize: '28px', fontWeight: '800', letterSpacing: '-0.5px' }}>Destinasi Pilihan</h1>
        </div>
        <Link className="button-circle" to="/trips/new" title="Buat trip">
          <Plus size={20} />
        </Link>
      </div>
      <ErrorBanner message={error} />
      {trips.length ? (
        <div className="trip-list">
          {trips.map((trip, index) => (
            <TripCard index={index} key={trip.id} trip={trip} />
          ))}
        </div>
      ) : (
        <div className="empty-card">Belum ada trip. Buat itinerary pertama dan undang member lewat email.</div>
      )}
    </section>
  )
}

export default TripList

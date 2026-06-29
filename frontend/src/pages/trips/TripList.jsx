import { Plus } from 'lucide-react'
import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import ErrorBanner from '../../components/common/ErrorBanner'
import PageSkeleton from '../../components/common/PageSkeleton'
import RetryState from '../../components/common/RetryState'
import TripCard from '../../components/trips/TripCard'
import tripService from '../../services/tripService'
import { getErrorMessage, getResponseData } from '../../services/api'

const TripList = () => {
  const [trips, setTrips] = useState([])
  const [error, setError] = useState('')
  const [requestKey, setRequestKey] = useState(0)
  const [loadedKey, setLoadedKey] = useState(null)

  const retryLoadTrips = () => {
    setRequestKey((current) => current + 1)
  }

  useEffect(() => {
    let active = true
    tripService
      .getAll()
      .then((response) => {
        if (!active) return
        setTrips(getResponseData(response) || [])
        setError('')
      })
      .catch((err) => {
        if (active) setError(getErrorMessage(err, 'Gagal memuat trips'))
      })
      .finally(() => {
        if (active) setLoadedKey(requestKey)
      })

    return () => {
      active = false
    }
  }, [requestKey])

  const loading = loadedKey !== requestKey
  if (loading) return <PageSkeleton label="Memuat trips" rows={3} />

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
      {error ? <RetryState message={error} onRetry={retryLoadTrips} /> : <ErrorBanner message={error} />}
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

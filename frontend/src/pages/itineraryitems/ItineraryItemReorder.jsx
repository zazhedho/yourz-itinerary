import { useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'

import AccessDenied from '../../components/common/AccessDenied'
import ErrorBanner from '../../components/common/ErrorBanner'
import ReorderList from '../../components/itinerary/ReorderList'
import Loading from '../../components/common/Loading'
import useTripAccess from '../../hooks/useTripAccess'
import { getErrorMessage } from '../../services/api'
import itineraryItemService from '../../services/itineraryItemService'

const ItineraryItemReorder = () => {
  const { dayId } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const accessTripId = location.state?.tripId || location.state?.day?.trip_id
  const { allowed, error: accessError, loading } = useTripAccess(accessTripId, 'edit')
  const [items, setItems] = useState(location.state?.day?.items || [])
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const move = (from, to) => {
    setItems((current) => {
      const next = [...current]
      const [item] = next.splice(from, 1)
      next.splice(to, 0, item)
      return next
    })
  }

  const save = async () => {
    setSubmitting(true)
    setError('')
    try {
      await itineraryItemService.reorder(dayId, items.map((item) => item.id))
      navigate(-1)
    } catch (err) {
      setError(getErrorMessage(err, 'Gagal menyimpan urutan aktivitas'))
    } finally {
      setSubmitting(false)
    }
  }

  if (!accessTripId) return <AccessDenied message="Buka form dari detail trip agar akses bisa diverifikasi." />
  if (loading) return <Loading label="Memeriksa akses trip..." />
  if (!allowed) return <AccessDenied backTo={`/trips/${accessTripId}`} message={accessError || 'Hanya owner dan editor yang bisa menyusun aktivitas.'} />

  return (
    <section className="screen-stack">
      <div className="section-header">
        <div>
          <p className="eyebrow">Reorder</p>
          <h2>Susun ulang aktivitas</h2>
        </div>
      </div>
      <ErrorBanner message={error} />
      <ReorderList items={items} onMove={move} />
      <button className="button-primary" disabled={submitting || !items.length} onClick={save} type="button">
        {submitting ? 'Menyimpan...' : 'Simpan urutan'}
      </button>
    </section>
  )
}

export default ItineraryItemReorder

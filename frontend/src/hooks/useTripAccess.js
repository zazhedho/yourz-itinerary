import { useEffect, useMemo, useState } from 'react'

import tripService from '../services/tripService'
import { getErrorMessage, getResponseData } from '../services/api'
import { canAccessTripAction, getTripAccess } from '../utils/tripAccess'
import { useAuth } from './useAuth'

const useTripAccess = (tripId, action) => {
  const { user } = useAuth()
  const [result, setResult] = useState({ error: '', trip: null, tripId: null })

  useEffect(() => {
    if (!tripId) return

    let active = true
    tripService
      .getById(tripId)
      .then((response) => {
        if (active) setResult({ error: '', trip: getResponseData(response), tripId })
      })
      .catch((err) => {
        if (active) {
          setResult({
            error: getErrorMessage(err, 'Gagal memeriksa akses trip'),
            trip: null,
            tripId,
          })
        }
      })

    return () => {
      active = false
    }
  }, [tripId])

  const loading = Boolean(tripId) && result.tripId !== tripId
  const trip = !loading && result.tripId === tripId ? result.trip : null
  const error = !loading && result.tripId === tripId ? result.error : ''
  const access = useMemo(() => getTripAccess(trip, user), [trip, user])
  const allowed = Boolean(tripId) && !error && canAccessTripAction(access, action)

  return { access, allowed, error, loading, trip }
}

export default useTripAccess

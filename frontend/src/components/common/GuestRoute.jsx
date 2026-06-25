import { Navigate, Outlet } from 'react-router-dom'

import { useAuth } from '../../hooks/useAuth'
import Loading from './Loading'

const GuestRoute = () => {
  const { booting, isAuthenticated } = useAuth()

  if (booting) return <Loading label="Menyiapkan sesi..." />
  if (isAuthenticated) return <Navigate to="/trips" replace />

  return <Outlet />
}

export default GuestRoute

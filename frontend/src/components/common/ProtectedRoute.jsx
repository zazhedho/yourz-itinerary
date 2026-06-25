import { Navigate, Outlet } from 'react-router-dom'

import { useAuth } from '../../hooks/useAuth'
import Loading from './Loading'

const ProtectedRoute = () => {
  const { booting, isAuthenticated } = useAuth()

  if (booting) return <Loading label="Menyiapkan sesi..." />
  if (!isAuthenticated) return <Navigate to="/login" replace />

  return <Outlet />
}

export default ProtectedRoute

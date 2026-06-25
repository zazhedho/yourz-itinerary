import { CalendarDays, MapPin, Plus, UserRound } from 'lucide-react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'

import { useAuth } from '../../hooks/useAuth'

const AppShell = () => {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login', { replace: true })
  }

  return (
    <div className="app-shell">
      <header className="mobile-header">
        <div>
          <p className="eyebrow">Yourz Itinerary</p>
          <h1>Trip bersama</h1>
        </div>
        <button className="avatar-button" onClick={handleLogout} type="button" title="Logout">
          <UserRound size={20} />
        </button>
      </header>

      <main className="app-main">
        <Outlet context={{ user }} />
      </main>

      <nav className="bottom-nav" aria-label="Primary">
        <NavLink to="/trips">
          <CalendarDays size={20} />
          <span>Trips</span>
        </NavLink>
        <NavLink to="/trips/new">
          <Plus size={22} />
          <span>Buat</span>
        </NavLink>
        <NavLink to="/trips">
          <MapPin size={20} />
          <span>Map</span>
        </NavLink>
      </nav>
    </div>
  )
}

export default AppShell

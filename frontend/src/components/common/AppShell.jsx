import { CalendarDays, MapPin, Plus, UserRound } from 'lucide-react'
import { Link, NavLink, Outlet } from 'react-router-dom'

import { useAuth } from '../../hooks/useAuth'

const AppShell = () => {
  const { user } = useAuth()

  return (
    <div className="app-shell">
      <header className="mobile-header">
        <div>
          <p className="eyebrow">Yourz Itinerary</p>
          <h1>Trip bersama</h1>
        </div>
        <Link className="avatar-button" to="/account" title="Akun">
          <UserRound size={20} />
        </Link>
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

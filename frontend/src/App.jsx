import { lazy, Suspense } from 'react'
import { Navigate, Route, Routes } from 'react-router-dom'

import AppShell from './components/common/AppShell'
import GuestRoute from './components/common/GuestRoute'
import Loading from './components/common/Loading'
import ProtectedRoute from './components/common/ProtectedRoute'

const ChangePassword = lazy(() => import('./pages/account/ChangePassword'))
const ForgotPassword = lazy(() => import('./pages/auth/ForgotPassword'))
const Login = lazy(() => import('./pages/auth/Login'))
const Profile = lazy(() => import('./pages/account/Profile'))
const Register = lazy(() => import('./pages/auth/Register'))
const ResetPassword = lazy(() => import('./pages/auth/ResetPassword'))
const Sessions = lazy(() => import('./pages/account/Sessions'))
const TripList = lazy(() => import('./pages/trips/TripList'))
const TripDetail = lazy(() => import('./pages/trips/TripDetail'))
const TripForm = lazy(() => import('./pages/trips/TripForm'))
const TripMemberAdd = lazy(() => import('./pages/tripmembers/TripMemberAdd'))
const ItineraryDayForm = lazy(() => import('./pages/itinerarydays/ItineraryDayForm'))
const ItineraryItemForm = lazy(() => import('./pages/itineraryitems/ItineraryItemForm'))
const ItineraryItemReorder = lazy(() => import('./pages/itineraryitems/ItineraryItemReorder'))
const TripMemberRoleForm = lazy(() => import('./pages/tripmembers/TripMemberRoleForm'))

const App = () => (
  <Suspense fallback={<Loading />}>
    <Routes>
      <Route element={<GuestRoute />}>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/reset-password" element={<ResetPassword />} />
      </Route>

      <Route element={<ProtectedRoute />}>
        <Route element={<AppShell />}>
          <Route index element={<Navigate to="/trips" replace />} />
          <Route path="/account" element={<Profile />} />
          <Route path="/account/password" element={<ChangePassword />} />
          <Route path="/account/sessions" element={<Sessions />} />
          <Route path="/trips" element={<TripList />} />
          <Route path="/trips/new" element={<TripForm />} />
          <Route path="/trips/:tripId" element={<TripDetail />} />
          <Route path="/trips/:tripId/edit" element={<TripForm />} />
          <Route path="/trips/:tripId/members/add" element={<TripMemberAdd />} />
          <Route path="/trips/:tripId/members/:memberId/role" element={<TripMemberRoleForm />} />
          <Route path="/trips/:tripId/days/new" element={<ItineraryDayForm />} />
          <Route path="/itinerary-days/:dayId/edit" element={<ItineraryDayForm />} />
          <Route path="/itinerary-days/:dayId/items/new" element={<ItineraryItemForm />} />
          <Route path="/itinerary-days/:dayId/items/reorder" element={<ItineraryItemReorder />} />
          <Route path="/itinerary-items/:itemId/edit" element={<ItineraryItemForm />} />
        </Route>
      </Route>

      <Route path="*" element={<Navigate to="/trips" replace />} />
    </Routes>
  </Suspense>
)

export default App

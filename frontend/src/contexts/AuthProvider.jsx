import { useEffect, useMemo, useState } from 'react'

import authService from '../services/authService'
import { getErrorMessage, getResponseData } from '../services/api'
import { AuthContext } from './auth-context'

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null)
  const [booting, setBooting] = useState(() => Boolean(localStorage.getItem('token')))
  const [error, setError] = useState('')

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) return

    authService
      .me()
      .then((response) => setUser(getResponseData(response)))
      .catch(() => localStorage.removeItem('token'))
      .finally(() => setBooting(false))
  }, [])

  const login = async (payload) => {
    setError('')
    try {
      const response = await authService.login(payload)
      const data = getResponseData(response) || {}
      const token = data.access_token || data.token
      if (token) localStorage.setItem('token', token)
      const profile = await authService.me()
      setUser(getResponseData(profile))
      return true
    } catch (err) {
      setError(getErrorMessage(err, 'Login failed'))
      return false
    }
  }

  const logout = async () => {
    try {
      await authService.logout()
    } finally {
      localStorage.removeItem('token')
      setUser(null)
    }
  }

  const value = useMemo(
    () => ({
      user,
      booting,
      error,
      isAuthenticated: Boolean(user),
      login,
      logout,
    }),
    [user, booting, error],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

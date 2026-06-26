import { useCallback, useEffect, useMemo, useState } from 'react'

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

  const applyAuthResponse = useCallback(async (response) => {
    const data = getResponseData(response) || {}
    const token = data.access_token || data.token
    const refreshToken = data.refresh_token
    if (token) localStorage.setItem('token', token)
    if (refreshToken) localStorage.setItem('refresh_token', refreshToken)
    const profile = await authService.me()
    setUser(getResponseData(profile))
  }, [])

  const login = useCallback(async (payload) => {
    setError('')
    try {
      const response = await authService.login(payload)
      await applyAuthResponse(response)
      return true
    } catch (err) {
      setError(getErrorMessage(err, 'Login failed'))
      return false
    }
  }, [applyAuthResponse])

  const googleLogin = useCallback(async (idToken) => {
    setError('')
    try {
      const response = await authService.googleLogin({ id_token: idToken })
      await applyAuthResponse(response)
      return true
    } catch (err) {
      setError(getErrorMessage(err, 'Google login failed'))
      return false
    }
  }, [applyAuthResponse])

  const logout = useCallback(async () => {
    try {
      await authService.logout()
    } finally {
      localStorage.removeItem('token')
      localStorage.removeItem('refresh_token')
      setUser(null)
    }
  }, [])

  const value = useMemo(
    () => ({
      user,
      booting,
      error,
      isAuthenticated: Boolean(user),
      googleLogin,
      login,
      logout,
    }),
    [user, booting, error, googleLogin, login, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

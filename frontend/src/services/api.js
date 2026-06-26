import axios from 'axios'

const API_BASE_URL = window.ENV_CONFIG?.API_URL || import.meta.env.VITE_API_URL || '/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

let refreshPromise = null

const refreshAccessToken = async () => {
  const refreshToken = localStorage.getItem('refresh_token')
  if (!refreshToken) throw new Error('missing refresh token')
  const response = await axios.post(`${API_BASE_URL}/user/refresh-token`, { refresh_token: refreshToken })
  const data = response?.data?.data || {}
  const token = data.access_token || data.token
  if (!token) throw new Error('missing access token')
  localStorage.setItem('token', token)
  if (data.refresh_token) localStorage.setItem('refresh_token', data.refresh_token)
  return token
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config
    if (error.response?.status === 401 && !original?._retry) {
      original._retry = true
      try {
        refreshPromise = refreshPromise || refreshAccessToken()
        const token = await refreshPromise
        refreshPromise = null
        original.headers.Authorization = `Bearer ${token}`
        return api(original)
      } catch (refreshError) {
        refreshPromise = null
        localStorage.removeItem('token')
        localStorage.removeItem('refresh_token')
        return Promise.reject(refreshError)
      }
    }
    return Promise.reject(error)
  },
)

export const getResponseData = (response) => response?.data?.data ?? null

export const getErrorMessage = (error, fallback = 'Request failed') => {
  const payload = error?.response?.data
  if (Array.isArray(payload?.error)) return payload.error.map((item) => item.message).join(', ')
  if (payload?.error?.message) return payload.error.message
  if (typeof payload?.error === 'string') return payload.error
  if (payload?.message) return payload.message
  return fallback
}

export default api

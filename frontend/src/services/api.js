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

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
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

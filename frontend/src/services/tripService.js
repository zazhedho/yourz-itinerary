import api from './api'

const tripService = {
  getAll: () => api.get('/trips'),
  create: (payload) => api.post('/trips', payload),
  getById: (id) => api.get(`/trips/${id}`),
  update: (id, payload) => api.put(`/trips/${id}`, payload),
  delete: (id) => api.delete(`/trips/${id}`),
}

export default tripService

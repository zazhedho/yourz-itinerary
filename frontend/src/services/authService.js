import api from './api'

const authService = {
  login: (payload) => api.post('/user/login', payload),
  register: (payload) => api.post('/user/register', payload),
  me: () => api.get('/user'),
  logout: () => api.post('/user/logout'),
}

export default authService

import api from './api'

const sessionService = {
  getActiveSessions: () => api.get('/user/sessions'),
  revoke: (sessionId) => api.delete(`/user/session/${sessionId}`),
  revokeOthers: () => api.post('/user/sessions/revoke-others'),
}

export default sessionService

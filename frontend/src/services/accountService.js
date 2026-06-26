import api from './api'

const accountService = {
  updateProfile: (payload) => api.put('/user', payload),
  changePassword: (payload) => api.put('/user/change/password', payload),
  deleteAccount: () => api.delete('/user'),
}

export default accountService

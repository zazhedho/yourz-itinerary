import api from './api'

const tripMemberService = {
  addMember: (tripId, payload) => api.post(`/trips/${tripId}/members`, payload),
  updateRole: (tripId, memberId, payload) => api.put(`/trips/${tripId}/members/${memberId}`, payload),
  remove: (tripId, memberId) => api.delete(`/trips/${tripId}/members/${memberId}`),
  leaveTrip: (tripId) => api.delete(`/trips/${tripId}/leave`),
}

export default tripMemberService

import api from './api'

const itineraryDayService = {
  create: (tripId, payload) => api.post(`/trips/${tripId}/days`, payload),
  update: (dayId, payload) => api.put(`/itinerary-days/${dayId}`, payload),
  delete: (dayId) => api.delete(`/itinerary-days/${dayId}`),
}

export default itineraryDayService

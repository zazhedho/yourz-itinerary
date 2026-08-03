import api from './api'

const weatherService = {
  getByDay: (dayId) => api.get(`/itinerary-days/${dayId}/weather`),
}

export default weatherService

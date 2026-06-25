import api from './api'

const itineraryItemService = {
  create: (dayId, payload) => api.post(`/itinerary-days/${dayId}/items`, payload),
  update: (itemId, payload) => api.put(`/itinerary-items/${itemId}`, payload),
  delete: (itemId) => api.delete(`/itinerary-items/${itemId}`),
  reorder: (dayId, itemIds) => api.put(`/itinerary-days/${dayId}/items/reorder`, { item_ids: itemIds }),
}

export default itineraryItemService

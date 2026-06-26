const UNSPLASH_API_URL = 'https://api.unsplash.com'
const ACCESS_KEY = import.meta.env.VITE_UNSPLASH_ACCESS_KEY

const defaultPhotos = [
  'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=900&q=80',
  'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=900&q=80',
  'https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?auto=format&fit=crop&w=900&q=80',
]

/**
 * Fetch a cover photo from Unsplash based on the destination name.
 * Uses localStorage caching to prevent hitting the 50 req/hr rate limit during dev.
 */
export const getDestinationPhoto = async (destination, index = 0) => {
  if (!destination) {
    return defaultPhotos[index % defaultPhotos.length]
  }

  const query = destination.trim().toLowerCase()
  const cacheKey = `unsplash_covers_${query}`

  // Check cache first
  const cached = localStorage.getItem(cacheKey)
  if (cached) {
    try {
      const urls = JSON.parse(cached)
      if (Array.isArray(urls) && urls.length > 0) {
        return urls[index % urls.length]
      }
    } catch (e) {
      // old cache or invalid JSON, will refetch
    }
  }

  // If no access key or rate limit exceeded, return default
  if (!ACCESS_KEY) {
    return defaultPhotos[index % defaultPhotos.length]
  }

  try {
    const response = await fetch(
      `${UNSPLASH_API_URL}/search/photos?query=${encodeURIComponent(query + ' landmark')}&orientation=landscape&per_page=5`,
      {
        headers: {
          Authorization: `Client-ID ${ACCESS_KEY}`,
        },
      }
    )

    if (!response.ok) {
      throw new Error('Unsplash API error')
    }

    const data = await response.json()
    if (data.results && data.results.length > 0) {
      const urls = data.results.map(res => res.urls.regular)
      localStorage.setItem(cacheKey, JSON.stringify(urls))
      return urls[index % urls.length]
    }
    
    // No results found
    return defaultPhotos[index % defaultPhotos.length]
  } catch (error) {
    console.error('Failed to fetch from Unsplash:', error)
    return defaultPhotos[index % defaultPhotos.length]
  }
}

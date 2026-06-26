# Yourz Itinerary Frontend

## Environment

Local development can use `frontend/.env`:

```env
VITE_API_URL=http://localhost:8080/api
VITE_GOOGLE_CLIENT_ID=your-google-oauth-client-id.apps.googleusercontent.com
VITE_GOOGLE_MAPS_API_KEY=your-google-maps-api-key
VITE_GOOGLE_MAPS_MAP_ID=your-google-maps-map-id
```

Production can override runtime config without rebuilding by replacing `public/env-config.js`:

```js
window.ENV_CONFIG = {
  API_URL: 'https://your-domain.com/api',
  GOOGLE_CLIENT_ID: 'your-google-oauth-client-id.apps.googleusercontent.com',
  GOOGLE_MAPS_API_KEY: 'your-google-maps-api-key',
  GOOGLE_MAPS_MAP_ID: 'your-google-maps-map-id',
}
```

Notes:

- Google login/register buttons render only when `GOOGLE_CLIENT_ID` or `VITE_GOOGLE_CLIENT_ID` is configured.
- Backend must use the matching `GOOGLE_CLIENT_ID` or include it in `GOOGLE_CLIENT_IDS`.
- Google Maps search requires `GOOGLE_MAPS_API_KEY`.
- Advanced markers require `GOOGLE_MAPS_MAP_ID`; local fallback uses `DEMO_MAP_ID`.

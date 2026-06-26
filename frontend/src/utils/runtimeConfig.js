export const getRuntimeConfigValue = (runtimeKey, viteKey, fallback = '') =>
  String(window.ENV_CONFIG?.[runtimeKey] || import.meta.env[viteKey] || fallback).trim()

export const getGoogleClientId = () => getRuntimeConfigValue('GOOGLE_CLIENT_ID', 'VITE_GOOGLE_CLIENT_ID')

export const getGoogleMapsApiKey = () => getRuntimeConfigValue('GOOGLE_MAPS_API_KEY', 'VITE_GOOGLE_MAPS_API_KEY')

export const getGoogleMapsMapId = () => getRuntimeConfigValue('GOOGLE_MAPS_MAP_ID', 'VITE_GOOGLE_MAPS_MAP_ID', 'DEMO_MAP_ID')

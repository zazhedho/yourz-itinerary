import { Cloud, CloudRain, CloudSun, Sun } from 'lucide-react'

const icons = {
  CLEAR: Sun,
  MOSTLY_CLEAR: Sun,
  PARTLY_CLOUDY: CloudSun,
  MOSTLY_CLOUDY: Cloud,
  CLOUDY: Cloud,
  LIGHT_RAIN_SHOWERS: CloudRain,
  CHANCE_OF_SHOWERS: CloudRain,
  SCATTERED_SHOWERS: CloudRain,
  RAIN_SHOWERS: CloudRain,
  HEAVY_RAIN_SHOWERS: CloudRain,
  LIGHT_TO_MODERATE_RAIN: CloudRain,
  MODERATE_TO_HEAVY_RAIN: CloudRain,
  RAIN: CloudRain,
  LIGHT_RAIN: CloudRain,
  HEAVY_RAIN: CloudRain,
}

const ItemWeather = ({ weather, onOpen }) => {
  if (!weather || weather.status !== 'available') return null
  const Icon = icons[weather.condition_code] || CloudSun
  return (
    <button
      aria-label={`Lihat detail cuaca: ${weather.condition_description || 'prakiraan cuaca'}`}
      className="item-weather"
      onClick={onOpen}
      type="button"
    >
      <span className="item-weather-condition">
        <Icon aria-hidden="true" size={16} />
        <span>{weather.condition_description || 'Prakiraan cuaca'}</span>
      </span>
      <strong>{Math.round(weather.min_temperature_c)}–{Math.round(weather.max_temperature_c)}°C</strong>
      <small><CloudRain aria-hidden="true" size={12} /> Peluang hujan {weather.precipitation_probability}%</small>
    </button>
  )
}

export default ItemWeather

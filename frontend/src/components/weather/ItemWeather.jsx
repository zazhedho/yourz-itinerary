import { Cloud, CloudRain, CloudSun, Sun } from 'lucide-react'

import { formatWeatherTime } from '../../utils/formatters'

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
  const isHourly = weather.forecast_type === 'hourly'
  const temperature = isHourly
    ? `${Math.round(weather.temperature_c)}°C`
    : `${Math.round(weather.min_temperature_c)}–${Math.round(weather.max_temperature_c)}°C`
  const precipitation = isHourly
    ? `Hujan ${weather.precipitation_probability}% pukul ${formatWeatherTime(weather.forecast_time, weather.time_zone)}`
    : `Hujan ${weather.precipitation_probability}%`
  return (
    <button
      aria-label={`Lihat detail cuaca: ${weather.condition_description || 'prakiraan cuaca'}`}
      className={`item-weather ${isHourly ? 'item-weather-hourly' : ''}`}
      onClick={onOpen}
      type="button"
    >
      <span className="item-weather-condition">
        <Icon aria-hidden="true" size={16} />
        <span>{weather.condition_description || 'Prakiraan cuaca'}</span>
      </span>
      <span className="item-weather-values">
        <strong>{temperature}</strong>
        <small><CloudRain aria-hidden="true" size={12} /> {precipitation}</small>
      </span>
    </button>
  )
}

export default ItemWeather

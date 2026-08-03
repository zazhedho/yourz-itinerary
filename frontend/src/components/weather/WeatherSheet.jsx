import { Wind, X } from 'lucide-react'
import { useEffect, useRef } from 'react'

const WeatherSheet = ({ weather, onClose, returnFocus }) => {
  const dialogRef = useRef(null)

  useEffect(() => {
    dialogRef.current?.focus()
    const handleKeyDown = (event) => {
      if (event.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', handleKeyDown)
    return () => {
      document.removeEventListener('keydown', handleKeyDown)
      returnFocus?.focus()
    }
  }, [onClose, returnFocus])

  if (!weather) return null
  const date = weather.forecast_date
    ? new Date(`${weather.forecast_date}T00:00:00`).toLocaleDateString('id-ID', { dateStyle: 'long' })
    : 'Tanggal tidak tersedia'
  const feelsLike = [weather.feels_like_min_c, weather.feels_like_max_c]
    .filter((value) => value != null)
    .map((value) => `${Math.round(value)}°C`)
    .join(' – ')

  return (
    <div className="weather-sheet-backdrop" onClick={onClose} role="presentation">
      <section
        aria-labelledby="weather-sheet-title"
        aria-modal="true"
        className="weather-sheet"
        onClick={(event) => event.stopPropagation()}
        ref={dialogRef}
        role="dialog"
        tabIndex="-1"
      >
        <div className="weather-sheet-heading">
          <div>
            <p className="eyebrow">Prakiraan cuaca</p>
            <h2 id="weather-sheet-title">{weather.condition_description || weather.condition_code}</h2>
            <span>{date}</span>
          </div>
          <button aria-label="Tutup detail cuaca" className="modal-close" onClick={onClose} type="button">
            <X size={18} />
          </button>
        </div>
        <div className="weather-sheet-temperature">
          <strong>{Math.round(weather.min_temperature_c)}–{Math.round(weather.max_temperature_c)}°C</strong>
          <span>Peluang hujan {weather.precipitation_probability}%</span>
        </div>
        <dl className="weather-detail-grid">
          {feelsLike && <div><dt>Suhu terasa</dt><dd>{feelsLike}</dd></div>}
          <div><dt>Kelembapan udara</dt><dd>{weather.humidity_percent}%</dd></div>
          <div><dt><Wind aria-hidden="true" size={14} /> Kecepatan angin</dt><dd>{weather.wind_speed_kph} km/j</dd></div>
        </dl>
        <p className="weather-sheet-attribution">
          <span>Source: Includes weather data from Google</span>
          <span className="google-maps-attribution" translate="no">Google Maps</span>
        </p>
      </section>
    </div>
  )
}

export default WeatherSheet

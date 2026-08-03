import { GripVertical, MapPin, Pencil, Plus, Trash2, ChevronDown, MoreHorizontal } from 'lucide-react'
import { useRef, useState } from 'react'
import { Link } from 'react-router-dom'

import ItemWeather from '../weather/ItemWeather'
import WeatherSheet from '../weather/WeatherSheet'
import { formatDate, formatMoney, formatTime } from '../../utils/formatters'

const DayTimeline = ({ days = [], currency = 'IDR', canEdit = true, onDeleteDay, onDeleteItem, onDayExpanded, onRetryWeather, weatherByDay = {} }) => {
  const [expandedDays, setExpandedDays] = useState({})
  const [expandedActions, setExpandedActions] = useState({})
  const [selectedWeather, setSelectedWeather] = useState(null)
  const [weatherTrigger, setWeatherTrigger] = useState(null)
  const weatherRequested = useRef({})

  const toggleDay = (dayId) => {
    const day = days.find((item) => item.id === dayId)
    const willExpand = !expandedDays[dayId]
    const hasCoordinates = (day?.items || []).some((item) => item.latitude != null && item.longitude != null)
    setExpandedDays(prev => ({
      ...prev,
      [dayId]: !prev[dayId]
    }))
    if (willExpand && day && hasCoordinates && !weatherRequested.current[day.id]) {
      weatherRequested.current[day.id] = true
      onDayExpanded?.(day)
    }
  }

  const toggleActions = (e, dayId) => {
    e.preventDefault()
    e.stopPropagation()
    setExpandedActions(prev => ({
      ...prev,
      [dayId]: !prev[dayId]
    }))
  }

  if (!days.length) {
    return <div className="empty-card">{canEdit ? 'Belum ada hari itinerary. Tambahkan hari pertama untuk mulai menyusun rencana.' : 'Belum ada hari itinerary.'}</div>
  }

  return (
    <div className="day-timeline">
      {days.map((day) => (
        <section className="day-card" key={day.id}>
          <div className="day-card-header" onClick={() => toggleDay(day.id)} style={{ cursor: 'pointer' }}>
            <div className="day-title-group">
              <button 
                type="button" 
                className={`collapse-toggle ${!expandedDays[day.id] ? 'collapsed' : ''}`}
                aria-label="Toggle collapse"
                onClick={(event) => { event.stopPropagation(); toggleDay(day.id) }}
              >
                <ChevronDown size={20} />
              </button>
              <div className="day-title-content">
                <div className="day-badge-row">
                  <span className="day-badge">Day {day.day_number}</span>
                  {day.title && day.date && <span className="day-date">{formatDate(day.date)}</span>}
                </div>
                <h3>{day.title || (day.date ? formatDate(day.date) : 'Rencana hari ini')}</h3>
              </div>
            </div>
            {canEdit && (
              <div className="day-actions-wrapper" onClick={e => e.stopPropagation()}>
                <div className={`inline-actions expander ${expandedActions[day.id] ? 'expanded' : ''}`}>
                  <Link className="icon-link" state={{ day }} to={`/itinerary-days/${day.id}/edit`} title="Edit hari">
                    <Pencil size={16} />
                  </Link>
                  <Link className="icon-link" state={{ day }} to={`/itinerary-days/${day.id}/items/reorder`} title="Susun item">
                    <GripVertical size={16} />
                  </Link>
                  <Link className="icon-link" state={{ tripId: day.trip_id }} to={`/itinerary-days/${day.id}/items/new`} title="Tambah item">
                    <Plus size={18} />
                  </Link>
                  <button aria-label="Hapus hari" className="icon-link danger" onClick={() => onDeleteDay?.(day)} type="button" title="Hapus hari">
                    <Trash2 size={16} />
                  </button>
                </div>
                <button
                  type="button"
                  className={`icon-link action-trigger ${expandedActions[day.id] ? 'active' : ''}`}
                  onClick={(e) => toggleActions(e, day.id)}
                  aria-label="Tampilkan opsi"
                >
                  <MoreHorizontal size={18} />
                </button>
              </div>
            )}
          </div>

          <div className={`item-list-wrapper ${!expandedDays[day.id] ? 'collapsed' : ''}`}>
            <div className="item-list-inner">
              {expandedDays[day.id] && weatherByDay[day.id]?.loading && <div aria-label="Memuat cuaca" className="weather-loading" role="status" />}
              {expandedDays[day.id] && weatherByDay[day.id]?.error && (
                <div className="weather-day-message" role="alert">
                  <span>{weatherByDay[day.id].error}</span>
                  <button onClick={() => onRetryWeather?.(day)} type="button">Coba lagi</button>
                </div>
              )}
              {expandedDays[day.id] && weatherByDay[day.id]?.data?.status === 'out_of_range' && <p className="weather-day-message">Prakiraan cuaca hanya tersedia untuk 10 hari ke depan.</p>}
              {expandedDays[day.id] && weatherByDay[day.id]?.data?.status === 'past_date' && <p className="weather-day-message">Prakiraan cuaca untuk tanggal ini sudah tidak tersedia.</p>}
              {expandedDays[day.id] && weatherByDay[day.id]?.data?.status === 'limit_reached' && <p className="weather-day-message">Prakiraan cuaca sedang tidak tersedia</p>}
              {expandedDays[day.id] && weatherByDay[day.id]?.data?.status === 'provider_unavailable' && (
                <div className="weather-day-message" role="alert">
                  <span>Prakiraan cuaca sedang tidak tersedia.</span>
                  <button onClick={() => onRetryWeather?.(day)} type="button">Coba lagi</button>
                </div>
              )}
              <div className="item-list">
                {(day.items || []).length ? (
                  (day.items || []).map((item) => (
                    <article className="item-row" key={item.id}>
                      <div className="item-time-column">
                        <div className="item-time">
                          {formatTime(item.start_time) || '--:--'}
                          {item.end_time && <span className="time-separator"> - {formatTime(item.end_time)}</span>}
                        </div>
                        <ItemWeather
                          onOpen={(event) => {
                            setSelectedWeather(weatherByDay[day.id]?.data?.items?.find((weather) => weather.item_id === item.id))
                            setWeatherTrigger(event.currentTarget)
                          }}
                          weather={weatherByDay[day.id]?.data?.items?.find((weather) => weather.item_id === item.id)}
                        />
                      </div>
                      <div className="item-content">
                        <div className="item-title-row">
                          <h4>{item.title}</h4>
                          {canEdit && (
                            <div className="inline-actions compact">
                              <Link className="icon-link" state={{ item, tripId: day.trip_id }} to={`/itinerary-items/${item.id}/edit`} title="Edit aktivitas">
                                <Pencil size={14} />
                              </Link>
                              <button aria-label="Hapus aktivitas" className="icon-link danger" onClick={() => onDeleteItem?.(item)} type="button" title="Hapus aktivitas">
                                <Trash2 size={14} />
                              </button>
                            </div>
                          )}
                        </div>
                        {item.description && <p className="item-note">{item.description}</p>}
                          <div className="item-meta">
                          {item.location_name && (
                            <a
                              href={`https://www.google.com/maps/search/?api=1&query=${
                                item.latitude && item.longitude 
                                  ? `${item.latitude},${item.longitude}` 
                                  : encodeURIComponent(item.location_name)
                              }`}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="item-location"
                              title="Buka di Google Maps"
                            >
                              <MapPin size={12} />
                              {item.location_name}
                            </a>
                          )}
                          <span className="item-cost">
                            {formatMoney(item.cost_estimate, currency)}
                          </span>
                        </div>
                      </div>
                    </article>
                  ))
                ) : (
                  <div className="day-empty-state">
                    <span>Belum ada aktivitas di hari ini.</span>
                    {canEdit && <Link state={{ tripId: day.trip_id }} to={`/itinerary-days/${day.id}/items/new`}>Tambah aktivitas</Link>}
                  </div>
                )}
              </div>
              {expandedDays[day.id] && weatherByDay[day.id]?.data?.items?.some((weather) => weather.status === 'available') && (
                <p className="weather-inline-attribution" translate="no">Google Maps</p>
              )}
            </div>
          </div>
        </section>
      ))}
      <WeatherSheet onClose={() => setSelectedWeather(null)} returnFocus={weatherTrigger} weather={selectedWeather} />
    </div>
  )
}

export default DayTimeline

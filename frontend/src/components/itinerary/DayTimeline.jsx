import { GripVertical, MapPin, Pencil, Plus, Trash2, ChevronDown, MoreHorizontal } from 'lucide-react'
import { useState } from 'react'
import { Link } from 'react-router-dom'

import { formatDate, formatMoney, formatShortDateTime, formatTime } from '../../utils/formatters'

const DayTimeline = ({ days = [], currency = 'IDR', canEdit = true, memberNameByUserId = {}, onDeleteDay, onDeleteItem }) => {
  const [expandedDays, setExpandedDays] = useState({})
  const [expandedActions, setExpandedActions] = useState({})

  const toggleDay = (dayId) => {
    setExpandedDays(prev => ({
      ...prev,
      [dayId]: !prev[dayId]
    }))
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
              <div className="item-list">
                {(day.items || []).length ? (
                  (day.items || []).map((item) => (
                    <article className="item-row" key={item.id}>
                      <div className="item-time">
                        {formatTime(item.start_time) || '--:--'}
                        {item.end_time && <span className="time-separator"> - {formatTime(item.end_time)}</span>}
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
                          {item.created_at && (
                            <span className="item-audit">
                              Dibuat {memberNameByUserId[item.created_by] || 'member'} • {formatShortDateTime(item.created_at)}
                            </span>
                          )}
                          {item.updated_at && (
                            <span className="item-audit">
                              Diubah {memberNameByUserId[item.updated_by] || 'member'} • {formatShortDateTime(item.updated_at)}
                            </span>
                          )}
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
            </div>
          </div>
        </section>
      ))}
    </div>
  )
}

export default DayTimeline

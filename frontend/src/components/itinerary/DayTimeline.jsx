import { GripVertical, MapPin, Pencil, Plus, Trash2 } from 'lucide-react'
import { Link } from 'react-router-dom'

import { formatMoney } from '../../utils/formatters'

const DayTimeline = ({ days = [], currency = 'IDR', onDeleteDay, onDeleteItem }) => {
  if (!days.length) {
    return <div className="empty-card">Belum ada hari itinerary. Tambahkan hari pertama untuk mulai menyusun rencana.</div>
  }

  return (
    <div className="day-timeline">
      {days.map((day) => (
        <section className="day-card" key={day.id}>
          <div className="day-card-header">
            <div>
              <p className="eyebrow">Day {day.day_number}</p>
              <h3>{day.title || day.date || 'Rencana hari ini'}</h3>
            </div>
            <div className="inline-actions">
              <Link className="icon-link" state={{ day }} to={`/itinerary-days/${day.id}/edit`} title="Edit hari">
                <Pencil size={16} />
              </Link>
              <Link className="icon-link" state={{ day }} to={`/itinerary-days/${day.id}/items/reorder`} title="Susun item">
                <GripVertical size={16} />
              </Link>
              <Link className="icon-link" to={`/itinerary-days/${day.id}/items/new`} title="Tambah item">
                <Plus size={18} />
              </Link>
              <button aria-label="Hapus hari" className="icon-link danger" onClick={() => onDeleteDay?.(day)} type="button" title="Hapus hari">
                <Trash2 size={16} />
              </button>
            </div>
          </div>

          <div className="item-list">
            {(day.items || []).length ? (
              (day.items || []).map((item) => (
                <article className="item-row" key={item.id}>
                  <div className="item-time">{item.start_time || '--:--'}</div>
                  <div className="item-content">
                    <div className="item-title-row">
                      <h4>{item.title}</h4>
                      <div className="inline-actions compact">
                        <Link className="icon-link" state={{ item }} to={`/itinerary-items/${item.id}/edit`} title="Edit aktivitas">
                          <Pencil size={14} />
                        </Link>
                        <button aria-label="Hapus aktivitas" className="icon-link danger" onClick={() => onDeleteItem?.(item)} type="button" title="Hapus aktivitas">
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </div>
                    {item.location_name && (
                      <p>
                        <MapPin size={13} />
                        {item.location_name}
                      </p>
                    )}
                    <span>{formatMoney(item.cost_estimate, currency)}</span>
                  </div>
                </article>
              ))
            ) : (
              <div className="day-empty-state">
                <span>Belum ada aktivitas di hari ini.</span>
                <Link to={`/itinerary-days/${day.id}/items/new`}>Tambah aktivitas</Link>
              </div>
            )}
          </div>
        </section>
      ))}
    </div>
  )
}

export default DayTimeline

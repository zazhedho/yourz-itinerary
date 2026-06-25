import { MapPin, Plus } from 'lucide-react'
import { Link } from 'react-router-dom'

import { formatMoney } from '../../utils/formatters'

const DayTimeline = ({ days = [], currency = 'IDR' }) => {
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
            <Link className="icon-link" to={`/itinerary-days/${day.id}/items/new`} title="Tambah item">
              <Plus size={18} />
            </Link>
          </div>

          <div className="item-list">
            {(day.items || []).map((item) => (
              <article className="item-row" key={item.id}>
                <div className="item-time">{item.start_time || '--:--'}</div>
                <div className="item-content">
                  <h4>{item.title}</h4>
                  {item.location_name && (
                    <p>
                      <MapPin size={13} />
                      {item.location_name}
                    </p>
                  )}
                  <span>{formatMoney(item.cost_estimate, currency)}</span>
                </div>
              </article>
            ))}
          </div>
        </section>
      ))}
    </div>
  )
}

export default DayTimeline

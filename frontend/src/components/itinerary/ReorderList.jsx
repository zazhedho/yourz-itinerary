import { ArrowDown, ArrowUp } from 'lucide-react'

const ReorderList = ({ items, onMove }) => (
  <div className="reorder-list">
    {items.map((item, index) => (
      <article className="reorder-row" key={item.id}>
        <div className="reorder-index">
          {index + 1}
        </div>
        <div className="reorder-content">
          <strong>{item.title}</strong>
          <span>{item.location_name || 'Tanpa lokasi'}</span>
        </div>
        <div className="reorder-controls">
          <button className="reorder-btn" disabled={index === 0} onClick={() => onMove(index, index - 1)} type="button" aria-label="Naikkan">
            <ArrowUp size={16} />
          </button>
          <button className="reorder-btn" disabled={index === items.length - 1} onClick={() => onMove(index, index + 1)} type="button" aria-label="Turunkan">
            <ArrowDown size={16} />
          </button>
        </div>
      </article>
    ))}
  </div>
)

export default ReorderList

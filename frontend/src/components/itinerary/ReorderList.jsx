import { ArrowDown, ArrowUp } from 'lucide-react'

const ReorderList = ({ items, onMove }) => (
  <div className="reorder-list">
    {items.map((item, index) => (
      <article className="reorder-row" key={item.id}>
        <div>
          <strong>{item.title}</strong>
          <span>{item.location_name || 'Tanpa lokasi'}</span>
        </div>
        <div className="inline-actions">
          <button disabled={index === 0} onClick={() => onMove(index, index - 1)} type="button">
            <ArrowUp size={16} />
          </button>
          <button disabled={index === items.length - 1} onClick={() => onMove(index, index + 1)} type="button">
            <ArrowDown size={16} />
          </button>
        </div>
      </article>
    ))}
  </div>
)

export default ReorderList

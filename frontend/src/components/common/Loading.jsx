const Loading = ({ label = 'Memuat data...' }) => (
  <div className="loading-state" role="status">
    <span className="loading-dot" />
    <span>{label}</span>
  </div>
)

export default Loading

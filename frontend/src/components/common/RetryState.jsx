const RetryState = ({ message = 'Koneksi bermasalah. Coba lagi.', onRetry }) => (
  <div className="retry-state" role="alert">
    <strong>Data belum termuat</strong>
    <span>{message}</span>
    {onRetry && (
      <button className="button-secondary" onClick={onRetry} type="button">
        Coba lagi
      </button>
    )}
  </div>
)

export default RetryState

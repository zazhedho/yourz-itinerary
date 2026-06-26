const ConfirmDialog = ({ open, title, message, confirmLabel = 'Konfirmasi', cancelLabel = 'Batal', danger, onCancel, onConfirm }) => {
  if (!open) return null

  return (
    <div className="dialog-backdrop" role="presentation">
      <section className="dialog-card" role="dialog" aria-modal="true" aria-labelledby="confirm-title">
        <h2 id="confirm-title">{title}</h2>
        <p>{message}</p>
        <div className="dialog-actions">
          <button className="button-secondary" onClick={onCancel} type="button">
            {cancelLabel}
          </button>
          <button className={danger ? 'button-primary danger-button' : 'button-primary'} onClick={onConfirm} type="button">
            {confirmLabel}
          </button>
        </div>
      </section>
    </div>
  )
}

export default ConfirmDialog

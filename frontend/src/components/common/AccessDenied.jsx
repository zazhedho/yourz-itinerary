import { Link } from 'react-router-dom'

const AccessDenied = ({ backTo, message = 'Kamu tidak punya akses untuk aksi ini.' }) => (
  <section className="screen-stack">
    <div className="empty-card access-denied-card">
      <strong>Akses ditolak</strong>
      <span>{message}</span>
      {backTo && (
        <Link className="button-secondary" to={backTo}>
          Kembali
        </Link>
      )}
    </div>
  </section>
)

export default AccessDenied

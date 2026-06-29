const PageSkeleton = ({ label = 'Memuat data', rows = 3 }) => (
  <section aria-label={label} className="screen-stack skeleton-screen" role="status">
    <div className="skeleton-header">
      <span />
      <strong />
    </div>
    {Array.from({ length: rows }, (_, index) => (
      <div className="skeleton-row" data-testid="skeleton-row" key={index}>
        <span />
        <strong />
        <em />
      </div>
    ))}
  </section>
)

export default PageSkeleton

export const formatDateRange = (startDate, endDate) => {
  if (startDate && endDate) return `${startDate} - ${endDate}`
  if (startDate) return startDate
  if (endDate) return endDate
  return 'Tanggal belum diatur'
}

export const formatMoney = (value, currency = 'IDR') =>
  new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency,
    maximumFractionDigits: 0,
  }).format(Number(value || 0))

export const roleLabel = (role) => {
  if (role === 'owner') return 'Owner'
  if (role === 'editor') return 'Editor'
  if (role === 'viewer') return 'Viewer'
  return role || '-'
}

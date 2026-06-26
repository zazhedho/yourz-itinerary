export const formatDate = (dateString) => {
  if (!dateString) return ''
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return dateString
    return new Intl.DateTimeFormat('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    }).format(date)
  } catch {
    return dateString
  }
}

export const formatDateRange = (startDate, endDate) => {
  if (startDate && endDate) return `${formatDate(startDate)} - ${formatDate(endDate)}`
  if (startDate) return formatDate(startDate)
  if (endDate) return formatDate(endDate)
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

export const formatTime = (timeString) => {
  if (!timeString) return ''
  const parts = timeString.split(':')
  if (parts.length >= 2) {
    return `${parts[0]}:${parts[1]}`
  }
  return timeString
}

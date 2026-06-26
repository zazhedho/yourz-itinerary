export const emptyToUndefined = (value) => {
  if (value === '') return undefined
  return value
}

export const normalizeClockTime = (value) => {
  if (!value) return value
  const trimmed = String(value).trim()
  const match = trimmed.match(/^(\d{2}):(\d{2})(?::\d{2}(?:\.\d+)?)?$/)
  if (!match) return trimmed
  return `${match[1]}:${match[2]}`
}

const stripUndefined = (payload) =>
  Object.fromEntries(Object.entries(payload).filter(([, value]) => value !== undefined && value !== null))

export const buildTripPayload = (form) =>
  stripUndefined({
    title: form.title?.trim(),
    destination: emptyToUndefined(form.destination?.trim()),
    start_date: emptyToUndefined(form.start_date),
    end_date: emptyToUndefined(form.end_date),
    timezone: emptyToUndefined(form.timezone),
    currency_code: emptyToUndefined(form.currency_code),
    status: emptyToUndefined(form.status),
  })

export const buildItineraryDayPayload = (form) =>
  stripUndefined({
    day_number: Number(form.day_number),
    title: emptyToUndefined(form.title?.trim()),
    date: emptyToUndefined(form.date),
  })

export const buildItineraryItemPayload = (form) =>
  stripUndefined({
    title: form.title?.trim(),
    description: emptyToUndefined(form.description?.trim()),
    location_name: emptyToUndefined(form.location_name?.trim()),
    latitude: form.latitude === '' ? undefined : Number(form.latitude),
    longitude: form.longitude === '' ? undefined : Number(form.longitude),
    start_time: emptyToUndefined(normalizeClockTime(form.start_time)),
    end_time: emptyToUndefined(normalizeClockTime(form.end_time)),
    cost_estimate: Number(form.cost_estimate || 0),
    sort_order: form.sort_order ? Number(form.sort_order) : undefined,
  })

export const buildTripMemberPayload = (form) =>
  stripUndefined({
    email: form.email?.trim().toLowerCase(),
    role: form.role,
  })

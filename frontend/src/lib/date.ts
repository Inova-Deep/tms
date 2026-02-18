const pad = (n: number) => (n < 10 ? `0${n}` : `${n}`)

const toDate = (value: string | Date): Date | null => {
  if (value instanceof Date) return isNaN(value.getTime()) ? null : value
  const d = new Date(value)
  return isNaN(d.getTime()) ? null : d
}

export function formatDate(value: string | Date): string {
  const d = toDate(value)
  if (!d) return ''
  const year = d.getFullYear()
  const month = pad(d.getMonth() + 1)
  const day = pad(d.getDate())
  return `${year}-${month}-${day}`
}

export function formatDateTime(value: string | Date): string {
  const d = toDate(value)
  if (!d) return ''
  const date = formatDate(d)
  const hours = pad(d.getHours())
  const minutes = pad(d.getMinutes())
  return `${date} ${hours}:${minutes}`
}

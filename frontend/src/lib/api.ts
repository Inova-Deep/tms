import { useErrorStore } from '@/stores/errors'
import { useAuthStore } from '@/stores/auth'

const baseUrl = import.meta.env.VITE_API_BASE ?? '/api'

type ApiOptions = RequestInit & { parseJson?: boolean }

type ApiError = {
  status: number
  message: string
  detail?: string
}

function buildUrl(path: string) {
  if (path.startsWith('http')) return path
  const trimmed = path.startsWith('/') ? path : `/${path}`
  return `${baseUrl}${trimmed}`
}

function notifyError(err: ApiError) {
  const store = useErrorStore()
  store.push(err.message, err.detail)
}

export async function api<T>(path: string, options: ApiOptions = {}): Promise<T> {
  const { parseJson = true, headers, ...rest } = options
  const auth = useAuthStore()

  const response = await fetch(buildUrl(path), {
    headers: {
      'Content-Type': 'application/json',
      ...(auth.token ? { Authorization: `Bearer ${auth.token}` } : {}),
      ...headers,
    },
    ...rest,
  })

  if (!response.ok) {
    let detail = ''
    try {
      const body = await response.json()
      detail = body?.error || body?.message || ''
    } catch (_) {
      // fallback to text if not JSON
      try {
        detail = await response.text()
      } catch {
        detail = ''
      }
    }
    const err: ApiError = {
      status: response.status,
      message: `Request failed (${response.status})`,
      detail: detail || response.statusText,
    }
    notifyError(err)
    throw err
  }

  if (!parseJson) return undefined as T

  const contentType = response.headers.get('content-type') || ''
  if (contentType.includes('application/json')) {
    return (await response.json()) as T
  }

  // If backend accidentally returns HTML (e.g., dev server), surface a friendly error
  const text = await response.text()
  const err: ApiError = {
    status: response.status,
    message: 'Unexpected response format',
    detail: text.slice(0, 200),
  }
  notifyError(err)
  throw err
}

import type { ApiErrorBody } from '@/types/auth'

const API_BASE = '/api'

export class ApiError extends Error {
  readonly status: number
  readonly code: string

  constructor(status: number, code: string, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

async function parseError(response: Response): Promise<ApiError> {
  let code = 'REQUEST_FAILED'
  let message = response.statusText || 'Request failed'

  try {
    const body = (await response.json()) as ApiErrorBody
    if (body.error?.code) {
      code = body.error.code
    }
    if (body.error?.message) {
      message = body.error.message
    }
  } catch {
    // ignore JSON parse errors
  }

  return new ApiError(response.status, code, message)
}

export async function apiRequest<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const headers = new Headers(options.headers)
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    credentials: 'include',
    headers,
  })

  if (!response.ok) {
    throw await parseError(response)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

export const api = {
  get<T>(path: string) {
    return apiRequest<T>(path)
  },
  post<T>(path: string, body?: unknown) {
    return apiRequest<T>(path, {
      method: 'POST',
      body: body === undefined ? undefined : JSON.stringify(body),
    })
  },
  put<T>(path: string, body?: unknown) {
    return apiRequest<T>(path, {
      method: 'PUT',
      body: body === undefined ? undefined : JSON.stringify(body),
    })
  },
  delete(path: string) {
    return apiRequest<void>(path, { method: 'DELETE' })
  },
}

function buildQuery(params?: Record<string, string>): string {
  if (!params) {
    return ''
  }
  const search = new URLSearchParams(params)
  const qs = search.toString()
  return qs ? `?${qs}` : ''
}

export function apiGet<T>(path: string, params?: Record<string, string>): Promise<T> {
  return api.get<T>(`${path}${buildQuery(params)}`)
}

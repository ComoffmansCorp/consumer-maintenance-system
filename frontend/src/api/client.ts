import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios'
import type { AuthResponse, Problem } from '@/types'

const ACCESS_TOKEN_KEY = 'cms.accessToken'
const REFRESH_TOKEN_KEY = 'cms.refreshToken'

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY)
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
}

export function storeSession(auth: AuthResponse) {
  localStorage.setItem(ACCESS_TOKEN_KEY, auth.accessToken)
  localStorage.setItem(REFRESH_TOKEN_KEY, auth.refreshToken)
}

export function clearSession() {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
}

export const api = axios.create({ baseURL: '/api' })

api.interceptors.request.use((config) => {
  const token = getAccessToken()
  if (token) {
    config.headers.set('Authorization', `Bearer ${token}`)
  }
  return config
})

let refreshPromise: Promise<string | null> | null = null

async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = getRefreshToken()
  if (!refreshToken) return null
  try {
    const { data } = await axios.post<AuthResponse>('/api/auth/refresh', { refreshToken })
    storeSession(data)
    return data.accessToken
  } catch {
    clearSession()
    return null
  }
}

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as (InternalAxiosRequestConfig & { _retried?: boolean }) | undefined
    const isAuthEndpoint = original?.url?.includes('/auth/')

    if (error.response?.status === 401 && original && !original._retried && !isAuthEndpoint) {
      original._retried = true
      if (!refreshPromise) {
        refreshPromise = refreshAccessToken().finally(() => {
          refreshPromise = null
        })
      }
      const newToken = await refreshPromise
      if (newToken) {
        original.headers.set('Authorization', `Bearer ${newToken}`)
        return api(original)
      }
      clearSession()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  },
)

// fetchBlobUrl loads an authenticated file (PDF, photo) through the axios
// instance — a plain <a href> or <img src> would hit the API without the
// bearer token and get a 401, since browsers don't attach custom headers to
// those requests. Callers must revoke the returned URL when done with it.
export async function fetchBlobUrl(path: string): Promise<string> {
  const response = await api.get(path, { responseType: 'blob' })
  return URL.createObjectURL(response.data as Blob)
}

export function extractErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const problem = error.response?.data as Problem | undefined
    if (problem?.detail) return problem.detail
    if (problem?.title) return problem.title
    if (error.message) return error.message
  }
  return 'Что-то пошло не так, попробуйте ещё раз'
}

import { api } from './client'
import type { AuthResponse } from '@/types'

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  fullName: string
}

export const authApi = {
  login: (req: LoginRequest) => api.post<AuthResponse>('/auth/login', req).then((r) => r.data),
  registerClient: (req: RegisterRequest) =>
    api.post<AuthResponse>('/auth/register-client', req).then((r) => r.data),
  registerMaster: (req: RegisterRequest) =>
    api.post<AuthResponse>('/auth/register-master', req).then((r) => r.data),
  logout: (refreshToken: string) => api.post('/auth/logout', { refreshToken }),
}

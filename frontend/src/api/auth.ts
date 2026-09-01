import { api } from './client'
import type { AuthResponse, Role, UserDTO } from '@/types'

export interface LoginRequest {
  tenantCode: string
  username: string
  password: string
}

export interface CompanyRegistrationRequest {
  tenantName: string
  tenantCode: string
  plan: string
  username: string
  password: string
  fullName: string
}

export interface BootstrapSuperAdminRequest {
  username: string
  password: string
  fullName: string
}

export interface RegisterMarketplaceRequest {
  username: string
  password: string
  fullName: string
}

export const authApi = {
  login: (req: LoginRequest) => api.post<AuthResponse>('/auth/login', req).then((r) => r.data),
  registerCompany: (req: CompanyRegistrationRequest) =>
    api.post<AuthResponse>('/auth/register-company', req).then((r) => r.data),
  registerClient: (req: RegisterMarketplaceRequest) =>
    api.post<AuthResponse>('/auth/register-client', req).then((r) => r.data),
  registerMaster: (req: RegisterMarketplaceRequest) =>
    api.post<AuthResponse>('/auth/register-master', req).then((r) => r.data),
  bootstrapSuperAdmin: (req: BootstrapSuperAdminRequest) =>
    api.post<AuthResponse>('/auth/bootstrap-super-admin', req).then((r) => r.data),
  logout: (refreshToken: string) => api.post('/auth/logout', { refreshToken }),
}

export interface CreateTenantUserRequest {
  username: string
  password: string
  fullName: string
  role: Role
}

export const usersApi = {
  list: (role?: Role) => api.get<UserDTO[]>('/users', { params: role ? { role } : {} }).then((r) => r.data),
  create: (req: CreateTenantUserRequest) => api.post<UserDTO>('/users', req).then((r) => r.data),
}

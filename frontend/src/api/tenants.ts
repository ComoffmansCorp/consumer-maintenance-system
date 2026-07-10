import { api } from './client'
import type { TenantDTO, TenantPlan } from '@/types'

export const tenantsApi = {
  list: () => api.get<TenantDTO[]>('/platform/tenants').then((r) => r.data),
  create: (req: { name: string; code: string; plan: TenantPlan }) =>
    api.post<TenantDTO>('/platform/tenants', req).then((r) => r.data),
  updatePlan: (id: number, plan: TenantPlan) =>
    api.patch<TenantDTO>(`/platform/tenants/${id}/plan`, { plan }).then((r) => r.data),
  deactivate: (id: number) => api.delete<TenantDTO>(`/platform/tenants/${id}`).then((r) => r.data),
}

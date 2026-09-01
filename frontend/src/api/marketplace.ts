import { api } from './client'
import type { CategoryDTO, MasterProfileDTO, Page, RequestDTO, RequestStatus, ServiceDTO } from '@/types'

export interface CreateRequestRequest {
  serviceId: number
  description: string
  addressText: string
  latitude?: number
  longitude?: number
}

export interface RequestListParams {
  page?: number
  pageSize?: number
}

export interface UpdateMasterProfileRequest {
  city: string
  bio: string
  specializationIds: number[]
}

export const marketplaceApi = {
  listCategories: () => api.get<CategoryDTO[]>('/marketplace/categories').then((r) => r.data),
  listServices: (categoryId?: number) =>
    api
      .get<ServiceDTO[]>('/marketplace/services', { params: categoryId ? { categoryId } : {} })
      .then((r) => r.data),

  createRequest: (req: CreateRequestRequest) =>
    api.post<RequestDTO>('/marketplace/requests', req).then((r) => r.data),
  listMyRequests: (params: RequestListParams = {}) =>
    api.get<Page<RequestDTO>>('/marketplace/requests', { params }).then((r) => r.data),
  listOpenRequests: (params: RequestListParams = {}) =>
    api.get<Page<RequestDTO>>('/marketplace/requests/open', { params }).then((r) => r.data),
  getRequest: (id: number) => api.get<RequestDTO>(`/marketplace/requests/${id}`).then((r) => r.data),
  claimRequest: (id: number) => api.post<RequestDTO>(`/marketplace/requests/${id}/claim`).then((r) => r.data),
  completeRequest: (id: number) =>
    api.post<RequestDTO>(`/marketplace/requests/${id}/complete`).then((r) => r.data),
  cancelRequest: (id: number, reason: string) =>
    api.post<RequestDTO>(`/marketplace/requests/${id}/cancel`, { reason }).then((r) => r.data),

  getMasterProfile: () => api.get<MasterProfileDTO>('/marketplace/master/profile').then((r) => r.data),
  updateMasterProfile: (req: UpdateMasterProfileRequest) =>
    api.put<MasterProfileDTO>('/marketplace/master/profile', req).then((r) => r.data),
}

export const requestStatusLabels: Record<RequestStatus, string> = {
  OPEN: 'Открыта',
  IN_PROGRESS: 'В работе',
  COMPLETED: 'Выполнена',
  CANCELED: 'Отменена',
}

// "Мастерская" badge tones for the marketplace's own screens (indigo accent,
// not the admin's brass) — see MarketplaceLandingView.vue for the palette.
export const requestStatusStyles: Record<RequestStatus, { bg: string; fg: string }> = {
  OPEN: { bg: '#EFEBE1', fg: '#55524A' },
  IN_PROGRESS: { bg: '#E7E3FC', fg: '#5B4BE0' },
  COMPLETED: { bg: 'oklch(0.94 0.04 145)', fg: 'oklch(0.42 0.1 145)' },
  CANCELED: { bg: 'oklch(0.95 0.03 25)', fg: '#B3261E' },
}

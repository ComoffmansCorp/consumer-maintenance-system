import { api } from './client'
import type { FavoriteDTO, OfferDTO, Page, RequestDTO, RequestStatus } from '@/types'

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

export interface SubmitOfferRequest {
  price: number
  comment?: string
}

// GET /api/requests, /api/requests/open and /api/admin/requests are all
// paginated (Page<RequestDTO>) per internal/request/handler.go -- the old
// flat marketplace API returned bare arrays here, this is the real change.
export const requestsApi = {
  createRequest: (req: CreateRequestRequest) => api.post<RequestDTO>('/requests', req).then((r) => r.data),
  listMine: (params: RequestListParams = {}) =>
    api.get<Page<RequestDTO>>('/requests', { params }).then((r) => r.data),
  listOpen: (params: RequestListParams = {}) =>
    api.get<Page<RequestDTO>>('/requests/open', { params }).then((r) => r.data),
  getRequest: (id: number) => api.get<RequestDTO>(`/requests/${id}`).then((r) => r.data),
  completeRequest: (id: number) => api.post<RequestDTO>(`/requests/${id}/complete`).then((r) => r.data),
  cancelRequest: (id: number, reason: string) =>
    api.post<RequestDTO>(`/requests/${id}/cancel`, { reason }).then((r) => r.data),

  submitOffer: (requestId: number, req: SubmitOfferRequest) =>
    api.post<OfferDTO>(`/requests/${requestId}/offers`, req).then((r) => r.data),
  listOffers: (requestId: number) => api.get<OfferDTO[]>(`/requests/${requestId}/offers`).then((r) => r.data),
  acceptOffer: (requestId: number, offerId: number) =>
    api.post<RequestDTO>(`/requests/${requestId}/offers/${offerId}/accept`).then((r) => r.data),

  // Favorites: add/remove return 204 No Content, not a body.
  listFavorites: () => api.get<FavoriteDTO[]>('/requests/favorites').then((r) => r.data),
  addFavorite: (masterId: number) => api.post<void>('/requests/favorites', { masterId }).then(() => undefined),
  removeFavorite: (masterId: number) => api.delete<void>(`/requests/favorites/${masterId}`).then(() => undefined),

  listAdmin: (params: RequestListParams & { status?: RequestStatus } = {}) =>
    api.get<Page<RequestDTO>>('/admin/requests', { params }).then((r) => r.data),
}

export const requestStatusLabels: Record<RequestStatus, string> = {
  OPEN: 'Открыта',
  ASSIGNED: 'В работе',
  COMPLETED: 'Выполнена',
  CANCELED: 'Отменена',
}

// "Мастерская" badge tones (indigo accent, matching MarketplaceShell.vue's
// palette) -- reused across public request screens and the admin panel.
export const requestStatusStyles: Record<RequestStatus, { bg: string; fg: string }> = {
  OPEN: { bg: '#EFEBE1', fg: '#55524A' },
  ASSIGNED: { bg: '#E7E3FC', fg: '#5B4BE0' },
  COMPLETED: { bg: 'oklch(0.94 0.04 145)', fg: 'oklch(0.42 0.1 145)' },
  CANCELED: { bg: 'oklch(0.95 0.03 25)', fg: '#B3261E' },
}

export const offerStatusLabels: Record<string, string> = {
  PENDING: 'Ожидает',
  ACCEPTED: 'Принят',
  REJECTED: 'Отклонён',
  WITHDRAWN: 'Отозван',
}

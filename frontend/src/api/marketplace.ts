import { api } from './client'
import type {
  CategoryDTO,
  ServiceDTO,
  ProfileDTO,
  RequestDTO,
  OfferDTO,
  FavoriteDTO,
  ReviewDTO,
  PaymentDTO,
  MessageDTO,
  Page,
} from '@/types'

export const catalogApi = {
  listCategories: () => api.get<CategoryDTO[]>('/catalog/categories').then((r) => r.data),
  listServices: (categoryId?: number) =>
    api
      .get<ServiceDTO[]>('/catalog/services', { params: categoryId ? { categoryId } : undefined })
      .then((r) => r.data),
  createCategory: (body: { name: string; parentCategoryId?: number }) =>
    api.post<CategoryDTO>('/admin/categories', body).then((r) => r.data),
  updateCategory: (id: number, body: { name: string; active: boolean }) =>
    api.put<CategoryDTO>(`/admin/categories/${id}`, body).then((r) => r.data),
  createService: (body: {
    categoryId: number
    name: string
    description?: string
    priceFrom?: number
    priceTo?: number
    unit?: string
  }) => api.post<ServiceDTO>('/admin/services', body).then((r) => r.data),
  updateService: (
    id: number,
    body: {
      name: string
      description?: string
      priceFrom?: number
      priceTo?: number
      unit?: string
      active: boolean
    },
  ) => api.put<ServiceDTO>(`/admin/services/${id}`, body).then((r) => r.data),
}

export const masterApi = {
  getProfile: () => api.get<ProfileDTO>('/master/profile').then((r) => r.data),
  updateProfile: (body: { city: string; bio: string; specializationIds: number[] }) =>
    api.put<ProfileDTO>('/master/profile', body).then((r) => r.data),
  listReviews: (masterId: number, page = 1, pageSize = 20) =>
    api
      .get<Page<ReviewDTO>>(`/masters/${masterId}/reviews`, { params: { page, pageSize } })
      .then((r) => r.data),
}

export const requestsApi = {
  create: (body: {
    serviceId: number
    description: string
    addressText: string
    latitude?: number
    longitude?: number
  }) => api.post<RequestDTO>('/requests', body).then((r) => r.data),
  listMine: (page = 1, pageSize = 20) =>
    api.get<Page<RequestDTO>>('/requests', { params: { page, pageSize } }).then((r) => r.data),
  listOpen: (page = 1, pageSize = 20) =>
    api
      .get<Page<RequestDTO>>('/requests/open', { params: { page, pageSize } })
      .then((r) => r.data),
  get: (id: number) => api.get<RequestDTO>(`/requests/${id}`).then((r) => r.data),
  cancel: (id: number, reason: string) =>
    api.post<RequestDTO>(`/requests/${id}/cancel`, { reason }).then((r) => r.data),
  complete: (id: number) => api.post<RequestDTO>(`/requests/${id}/complete`).then((r) => r.data),
  listOffers: (id: number) => api.get<OfferDTO[]>(`/requests/${id}/offers`).then((r) => r.data),
  submitOffer: (id: number, body: { price: number; comment?: string }) =>
    api.post<OfferDTO>(`/requests/${id}/offers`, body).then((r) => r.data),
  acceptOffer: (id: number, offerId: number) =>
    api.post<RequestDTO>(`/requests/${id}/offers/${offerId}/accept`).then((r) => r.data),
  listFavorites: () => api.get<FavoriteDTO[]>('/requests/favorites').then((r) => r.data),
  addFavorite: (masterId: number) =>
    api.post<FavoriteDTO>('/requests/favorites', { masterId }).then((r) => r.data),
  removeFavorite: (masterId: number) => api.delete(`/requests/favorites/${masterId}`),
}

export const reviewsApi = {
  create: (body: { requestId: number; rating: number; comment?: string }) =>
    api.post<ReviewDTO>('/reviews', body).then((r) => r.data),
}

export const paymentsApi = {
  getForRequest: (requestId: number) =>
    api.get<PaymentDTO>(`/requests/${requestId}/payment`).then((r) => r.data),
}

export const chatApi = {
  listMessages: (requestId: number, sinceId?: number) =>
    api
      .get<MessageDTO[]>(`/requests/${requestId}/messages`, {
        params: sinceId ? { sinceId } : undefined,
      })
      .then((r) => r.data),
  send: (requestId: number, text: string) =>
    api.post<MessageDTO>(`/requests/${requestId}/messages`, { text }).then((r) => r.data),
}

export const adminApi = {
  listMasters: (page = 1, pageSize = 20) =>
    api.get<Page<ProfileDTO>>('/admin/masters', { params: { page, pageSize } }).then((r) => r.data),
  listRequests: (status?: string, page = 1, pageSize = 20) =>
    api
      .get<Page<RequestDTO>>('/admin/requests', { params: { status, page, pageSize } })
      .then((r) => r.data),
  hideReview: (id: number) => api.put(`/admin/reviews/${id}/hide`),
  listPayments: (page = 1, pageSize = 20) =>
    api
      .get<Page<PaymentDTO>>('/admin/payments', { params: { page, pageSize } })
      .then((r) => r.data),
}

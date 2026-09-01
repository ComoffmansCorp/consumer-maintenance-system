import { api } from './client'
import type { Page, ReviewDTO } from '@/types'

export interface CreateReviewRequest {
  requestId: number
  rating: number
  comment?: string
}

export const reviewsApi = {
  createReview: (req: CreateReviewRequest) => api.post<ReviewDTO>('/reviews', req).then((r) => r.data),
  // Public, and paginated -- GET /api/masters/{id}/reviews (mounted at root,
  // outside JWTAuth, see internal/server/router.go).
  listByMaster: (masterId: number, page = 1, pageSize = 20) =>
    api.get<Page<ReviewDTO>>(`/masters/${masterId}/reviews`, { params: { page, pageSize } }).then((r) => r.data),
  hideReview: (id: number) => api.put<ReviewDTO>(`/admin/reviews/${id}/hide`).then((r) => r.data),
}
